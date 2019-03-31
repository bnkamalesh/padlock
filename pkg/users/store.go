package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
)

const (
	usersTable = "users"
)

type store interface {
	Create(ctx context.Context, u User) (*User, error)
	Read(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, u User) (*User, error)
}

type dbStore struct {
	appCtx *appcontext.AppContext
	db     *sql.DB
}

func (dbs *dbStore) rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		dbs.appCtx.Logger.Error(err)
	}
}

func (dbs *dbStore) prepColVals(cols ...string) string {
	count := len(cols)
	if count == 0 {
		return ""
	}

	vals := make([]string, 0, count)
	for i := 0; i < count; i++ {
		vals = append(vals, fmt.Sprintf("$%d", i+1))
	}

	return fmt.Sprintf("(%s) VALUES(%s)", strings.Join(cols, ", "), strings.Join(vals, ", "))
}

func (dbs *dbStore) Create(ctx context.Context, u User) (*User, error) {
	stmt := fmt.Sprintf(
		"INSERT INTO %s %s",
		usersTable,
		dbs.prepColVals(
			"name",
			"email",
			"phone",
			"password",
			"salt",
			"createdat",
		),
	)

	result, err := dbs.db.Query(
		stmt,
		u.Name,
		u.Email,
		u.Phone,
		u.Password,
		u.Salt,
		u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		err = result.Scan(&u.ID)
		if err != nil {
			return nil, err
		}
	}

	return &u, nil
}

func (dbs *dbStore) Read(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}

func (dbs *dbStore) Update(ctx context.Context, u User) (*User, error) {
	return nil, nil
}
