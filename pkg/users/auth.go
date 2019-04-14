package users

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sessionIDLength = 24
	letterIdxBits   = 6                    // 6 bits to represent a letter index
	letterIdxMask   = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax    = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	app = "Padlock.dev"
)

var (
	src          = rand.NewSource(time.Now().UnixNano())
	signKey      = []byte(rdmStr(64))
	ErrSessionID = errors.New("Invalid session ID received")
)

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

type claims struct {
	IP string
	jwt.StandardClaims
}

func sessionID(referrer string, u *User) string {
	id := rdmStr(48)
	now := time.Now()
	expire := now.Add(time.Hour * 12)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims{
			referrer,
			jwt.StandardClaims{
				Audience:  app,
				Issuer:    app,
				Id:        id,
				ExpiresAt: expire.Unix(),
				IssuedAt:  now.Unix(),
			},
		},
	)
	str, _ := token.SignedString(signKey)

	return str
}

func sessionDetails(tokenStr string) (*User, *claims, error) {
	c := &claims{}
	token, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !token.Valid {
		return nil, nil, ErrSessionID
	}

	err = token.Claims.Valid()
	if err != nil {
		return nil, nil, err
	}
	return nil, c, nil
}

// Login signs in a user based on the for the given email & password
// And returns the user instance, as well as the authenticated session ID
func (us *Users) Login(ctx context.Context, referrer, email, password string) (*User, string, error) {
	if !emailRegex.Match([]byte(email)) {
		return nil, "", ErrInvalidEmail
	}

	u, err := us.store.ReadByEmail(ctx, email)
	if err != nil {
		us.appCtx.Logger.Error(err)
		return nil, "", ErrUnexpected
	}

	pwd := hash(u.Salt, password)
	if pwd != u.Password {
		return nil, "", ErrInvalidLogin
	}

	return u, sessionID(referrer, u), nil
}
