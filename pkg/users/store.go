package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
)

const (
	usersTable = "users"
)

type store interface {
	Create(ctx context.Context, u User) (*User, error)
	ReadByEmail(ctx context.Context, email string) (*User, error)
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
		"INSERT INTO %s %s RETURNING id",
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

	result := dbs.db.QueryRow(
		stmt,
		u.Name,
		u.Email,
		u.Phone,
		u.Password,
		u.Salt,
		u.CreatedAt,
	)
	err := result.Scan(&u.ID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (dbs *dbStore) ReadByEmail(ctx context.Context, email string) (*User, error) {
	stmt := fmt.Sprintf(
		"SELECT id,name,email,phone,password,salt,createdat,updatedat FROM %s WHERE email=$1",
		usersTable,
	)

	result := dbs.db.QueryRow(stmt, email)
	u := User{}

	createdAt := pq.NullTime{}
	updatedAt := pq.NullTime{}
	err := result.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Phone,
		&u.Password,
		&u.Salt,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	u.CreatedAt = &createdAt.Time
	u.UpdatedAt = &updatedAt.Time

	return &u, nil
}

func (dbs *dbStore) Update(ctx context.Context, u User) (*User, error) {
	stmt := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=%d",
		usersTable,
		"SET name=$1, email=$2, phone=$3, password=$4, salt=$5, createdat=$6, updatedat=$7",
		u.ID,
	)

	_, err := dbs.db.Exec(
		stmt,
		u.Name,
		u.Email,
		u.Phone,
		u.Password,
		u.Salt,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
