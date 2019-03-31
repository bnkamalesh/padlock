// Package apps are the registerd applications
package apps

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
	"github.com/bnkamalesh/padlock/pkg/totp"
	"github.com/bnkamalesh/padlock/pkg/users"
)

var (
	ErrUnexpected   = errors.New("Sorry, an unexpected error occurred")
	ErrInvalidName  = errors.New("Sorry, invalid/no application name provided")
	ErrInvalidID    = errors.New("Sorry, invalid/no application ID provided")
	ErrInvalidOwner = errors.New("Sorry, invalid/no owner provided")
)

// App holds all the info related to an application registered on this platform
type App struct {
	// ID is the unique ID of an application registered on this platform
	ID int64 `json:"id,omitempty"`
	// Name is a human friendly name for the application
	Name string `json:"name,omitempty"`
	// Description is a short description for the application
	Description string `json:"description,omitempty"`
	// TOTP stores all the base settings required for generating TOTP
	TOTP *totp.TOTP `json:"totp,omitempty"`
	// CreatedAt is the timestamp at which the application was registered on this platform
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// UpdatedAt is the timestamp at which the applicated was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Apps handles all the service methods made available by this package
type Apps struct {
	appCtx *appcontext.AppContext
	store  store
}

// Create accepts an App instance and inserts it in the data store. On success it'll return
// the pointer of the app instance which was insterted
func (a *Apps) Create(ctx context.Context, app App) (*App, error) {
	now := time.Now().UTC()
	app.CreatedAt = &now
	// resetting ID, since it should ideally be ignored while creating/registering a new application
	app.ID = 0
	app.Name = strings.TrimSpace(app.Name)
	if app.Name == "" {
		return nil, ErrInvalidName
	}

	ap, err := a.store.Create(ctx, app)
	if err != nil {
		a.appCtx.Logger.Error(err)
		return nil, ErrUnexpected
	}

	return ap, nil
}

// Read reads a single app from the store based on the given ID
func (a *Apps) Read(ctx context.Context, id int64) (*App, error) {
	if id < 1 {
		return nil, ErrInvalidID
	}

	app, err := a.store.Read(ctx, id)
	if err != nil {
		a.appCtx.Logger.Error(err)
		return nil, ErrUnexpected
	}
	return app, nil
}

// ReadAll reads all the records from the store for the given IDs
func (a *Apps) ReadAll(ctx context.Context, ids ...int64) ([]App, error) {
	validIDs := make([]int64, 0, len(ids))

	for _, id := range ids {
		if id < 1 {
			continue
		}
		validIDs = append(validIDs, id)
	}

	apps, err := a.store.ReadAll(ctx, validIDs...)
	if err != nil {
		a.appCtx.Logger.Error(err)
		return nil, ErrUnexpected
	}
	return apps, nil
}

// Update updates the details of the app in the store
func (a *Apps) Update(ctx context.Context, app App) (*App, error) {
	now := time.Now().UTC()
	app.UpdatedAt = &now

	oldApp, err := a.Read(ctx, app.ID)
	if err != nil {
		return nil, err
	}

	app.Name = strings.TrimSpace(app.Name)
	if app.Name == "" {
		app.Name = oldApp.Name
	}

	app.Description = strings.TrimSpace(app.Description)
	if app.Description == "" {
		app.Description = oldApp.Description
	}

	updatedApp, err := a.store.Update(ctx, app)
	if err != nil {
		a.appCtx.Logger.Error(err)
		return nil, ErrUnexpected
	}
	return updatedApp, nil
}

// Delete deletes the app from the store
func (a *Apps) Delete(ctx context.Context, app App) (*App, error) {
	existingApp, err := a.Read(ctx, app.ID)
	if err != nil {
		return nil, err
	}

	err = a.store.Delete(ctx, *existingApp)
	if err != nil {
		a.appCtx.Logger.Error(err)
		return nil, ErrUnexpected
	}
	return existingApp, nil
}

func (a *Apps) SetOwner(ctx context.Context, app App, u users.User) (*App, error) {
	if u.ID == 0 {
		return nil, ErrInvalidOwner
	}
	ap, err := a.store.SetOwner(ctx, app, u)
	if err != nil {
		a.appCtx.Logger.Error(err)
		return nil, ErrUnexpected
	}
	return ap, err
}

func (a *Apps) CreateAndSetOwner(ctx context.Context, app App, u users.User) (*App, error) {
	now := time.Now().UTC()
	app.CreatedAt = &now
	// resetting ID, since it should ideally be ignored while creating/registering a new application
	app.ID = 0
	app.Name = strings.TrimSpace(app.Name)
	if app.Name == "" {
		return nil, ErrInvalidName
	}

	if u.ID == 0 {
		return nil, ErrInvalidOwner
	}

	ap, err := a.store.CreateAndSetOwner(ctx, app, u)
	if err != nil {
		a.appCtx.Logger.Error(err)
		return nil, ErrUnexpected
	}

	return ap, nil
}

func New(appCtx *appcontext.AppContext, sdb *sql.DB) *Apps {
	dbs := &dbStore{
		db:     sdb,
		appCtx: appCtx,
	}
	return &Apps{
		appCtx: appCtx,
		store:  dbs,
	}
}
