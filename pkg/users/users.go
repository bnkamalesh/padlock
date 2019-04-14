package users

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
	"github.com/bnkamalesh/padlock/pkg/platform/cache"
)

var (
	hasher = sha512.New()

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	ErrInvalidEmail = errors.New("Invalid/no email provided")
	ErrInvalidUser  = errors.New("Invalid user")
	ErrInvalidLogin = errors.New("Email/password did not match")
	ErrUnexpected   = errors.New("Sorry, an unexpected error occurred")
)

type Users struct {
	appCtx *appcontext.AppContext
	store  store
	cache  cache.Cache
}

func (us *Users) Create(ctx context.Context, u User, password string) (*User, error) {
	u.Email = strings.TrimSpace(u.Email)

	if u.Email == "" || !emailRegex.Match([]byte(u.Email)) {
		return nil, ErrInvalidEmail
	}

	u.setPassword(password)
	now := time.Now()
	u.CreatedAt = &now

	usr, err := us.store.Create(ctx, u)
	if err != nil {
		if us.appCtx.Logging {
			us.appCtx.Logger.Error(err)
		}
		return nil, ErrUnexpected
	}

	return usr, nil
}

func (us *Users) Read(ctx context.Context, id int) (*User, error) {
	return nil, nil
}

func (us *Users) Update(ctx context.Context, u User) (*User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	u.UpdatedAt = &now

	usr, err := us.store.Update(ctx, u)
	if err != nil {
		if us.appCtx.Logging {
			us.appCtx.Logger.Error(err)
		}
		return nil, ErrUnexpected
	}

	return usr, nil
}

func (us *Users) ReadByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func New(appCtx *appcontext.AppContext, sdb *sql.DB, c cache.Cache) *Users {
	dbs := &dbStore{
		db: sdb,
	}
	u := &Users{
		appCtx: appCtx,
		cache:  c,
		store:  dbs,
	}
	return u
}

type User struct {
	ID        int64      `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Password  string     `json:"-"`
	Salt      string     `json:"-"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidUser
	}

	return nil
}

func (u *User) setPassword(password string) {
	u.Salt = uuid.New().String()
	u.Password = hash(u.Salt, password)
}

func hash(ss ...string) string {
	chksum := hasher.Sum([]byte(strings.Join(ss, "")))
	dst := make([]byte, len(chksum)*2)
	hex.Encode(dst, chksum)
	return string(dst)
}
