package apps

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
	"github.com/bnkamalesh/padlock/pkg/users"
)

const (
	appTable      = "applications"
	appOwnerTable = "applicationOwners"
)

type store interface {
	Create(ctx context.Context, app App) (*App, error)
	Read(ctx context.Context, id int64) (*App, error)
	ReadAll(ctx context.Context, ids ...int64) ([]App, error)
	Update(ctx context.Context, app App) (*App, error)
	Delete(ctx context.Context, app App) error
	SetOwner(ctx context.Context, app App, u users.User) (*App, error)
	CreateAndSetOwner(ctx context.Context, app App, u users.User) (*App, error)
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

func (dbs *dbStore) Create(ctx context.Context, app App) (*App, error) {
	stmt := fmt.Sprintf(
		"INSERT INTO %s %s RETURNING id",
		appTable,
		dbs.prepColVals("name", "description", "totp", "createdAt", "updatedAt"),
	)

	result := dbs.db.QueryRow(
		stmt,
		app.Name,
		app.Description,
		app.TOTP,
		app.CreatedAt,
		app.UpdatedAt,
	)
	err := result.Scan(&app.ID)
	if err != nil {
		return nil, err
	}
	fmt.Println("app.ID=", app.ID)
	return &app, nil
}

func (dbs *dbStore) Read(ctx context.Context, id int64) (*App, error) {
	return nil, nil
}
func (dbs *dbStore) ReadAll(ctx context.Context, ids ...int64) ([]App, error) {
	return nil, nil
}
func (dbs *dbStore) Update(ctx context.Context, app App) (*App, error) {
	return nil, nil
}
func (dbs *dbStore) Delete(ctx context.Context, app App) error {
	return nil
}

func (dbs *dbStore) SetOwner(ctx context.Context, app App, u users.User) (*App, error) {
	return &app, nil
}

func (dbs *dbStore) CreateAndSetOwner(ctx context.Context, app App, u users.User) (*App, error) {
	tx, err := dbs.db.Begin()
	if err != nil {
		return nil, err
	}
	defer dbs.rollback(tx)

	stmt := fmt.Sprintf(
		"INSERT INTO %s %s RETURNING id",
		appTable,
		dbs.prepColVals("name", "description", "totp", "createdAt", "updatedAt"),
	)

	b, err := json.Marshal(app.TOTP)
	if err != nil {
		return nil, err
	}

	result := tx.QueryRow(
		stmt,
		app.Name,
		app.Description,
		string(b),
		app.CreatedAt,
		app.UpdatedAt,
	)
	err = result.Scan(&app.ID)
	if err != nil {
		return nil, err
	}

	stmt = fmt.Sprintf(
		"INSERT INTO %s %s",
		appOwnerTable,
		dbs.prepColVals("appid", "userid", "createdat"),
	)

	now := time.Now()
	_, err = tx.Exec(
		stmt,
		app.ID,
		u.ID,
		now,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &app, nil
}
