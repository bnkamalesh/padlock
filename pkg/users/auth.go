package users

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/bnkamalesh/padlock/pkg/platform/cache"
)

type ctxKey string

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sessionIDLength = 24
	letterIdxBits   = 6                    // 6 bits to represent a letter index
	letterIdxMask   = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax    = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	app        = "padlock.dev"
	ctxUserKey = ctxKey("user")
)

var (
	src                 = rand.NewSource(time.Now().UnixNano())
	signKey             = []byte(rdmStr(64))
	ErrSessionID        = errors.New("Invalid session ID received")
	ErrSessionIDExpired = errors.New("Session ID expired")
)

func tokenKeyFunc(token *jwt.Token) (interface{}, error) {
	return signKey, nil
}

// rdmStr returns a random string of length n
// ref: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func rdmStr(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

type TokenClaims struct {
	Source string
	jwt.StandardClaims
}

func sessionID(source string, u *User) (string, *TokenClaims) {
	id := rdmStr(48)
	now := time.Now()
	expire := now.Add(time.Hour * 12)
	tc := TokenClaims{
		source,
		jwt.StandardClaims{
			Audience:  app,
			Issuer:    app,
			Id:        id,
			ExpiresAt: expire.Unix(),
			IssuedAt:  now.Unix(),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		tc,
	)
	str, _ := token.SignedString(signKey)

	return str, &tc
}

func sessionDetails(tokenStr string) (*TokenClaims, error) {
	c := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, c, tokenKeyFunc)
	if err != nil {
		switch err.Error() {
		case jwt.ErrInvalidKey.Error(), jwt.ErrSignatureInvalid.Error():
			{
				return nil, ErrSessionID
			}
		}
		return nil, err
	}
	if !token.Valid {
		return nil, ErrSessionID
	}

	err = token.Claims.Valid()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Login signs in a user based on the for the given email & password
// And returns the user instance, as well as the authenticated session ID
func (us *Users) Login(ctx context.Context, source, email, password string) (*User, string, error) {
	if !emailRegex.Match([]byte(email)) {
		return nil, "", ErrInvalidEmail
	}

	u, err := us.store.ReadByEmail(ctx, email)
	if err != nil {
		if us.appCtx.Logging {
			us.appCtx.Logger.Error(err)
		}
		return nil, "", ErrUnexpected
	}

	pwd := hash(u.Salt, password)
	if pwd != u.Password {
		return nil, "", ErrInvalidLogin
	}

	token, claims := sessionID(source, u)
	expiry := time.Until(time.Unix(claims.ExpiresAt, 0))
	err = us.cache.Set(claims.Id, u, expiry)
	if err != nil {
		if us.appCtx.Logging {
			us.appCtx.Logger.Error(err)
		}
		return nil, "", ErrUnexpected
	}

	return u, token, nil
}

func (us *Users) AuthUser(ctx context.Context, source, token string) (*User, error) {
	claims, err := sessionDetails(token)
	if err != nil {
		return nil, err
	}

	if source != claims.Source {
		return nil, ErrSessionID
	}

	u := &User{}
	err = us.cache.Get(claims.Id, u)
	if err != nil {
		if err == cache.ErrNotFound {
			return nil, ErrSessionIDExpired
		}
		if us.appCtx.Logging {
			us.appCtx.Logger.Error(err)
		}
		return nil, err
	}

	return u, nil
}

func SetContext(ctx context.Context, u *User) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(
		ctx,
		ctxUserKey,
		u,
	)
	return ctx
}

func FromContext(ctx context.Context) *User {
	if ctx == nil {
		return nil
	}
	u, ok := ctx.Value(ctxUserKey).(*User)
	if !ok {
		return nil
	}
	return u
}
