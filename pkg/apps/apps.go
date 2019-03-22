// Package apps are the registerd applications
package apps

import (
	"time"

	"github.com/bnkamalesh/padlock/pkg/totp"
)

// App holds all the info related to an application registered on this platform
type App struct {
	// ID is the unique ID of an application
	ID string `json:"id,omitempty"`
	// Name is human friendly name for the application
	Name string `json:"name,omitempty"`
	// Description is a short description for the application
	Description string `json:"description,omitempty"`
	// Secret of the app which is used for generating the TOTP
	Secret string `json:"secret,omitempty"`
	// CreatedAt is the timestamp at which the application was registered on this platform
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// UpdatedAt is the timestamp at which the applicated was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type store interface {
	// Insert inserts a new record in the data store, and returns the ID of the inserted record
	Insert(data interface{}) (string, error)
	// Update updates the record of the given ID, with the data. It should overwrite the whole
	// record which exists on the data store.
	Update(id string, data interface{}) error
	// Delete deletes a record for the given ID
	Delete(id string) error
}

// Apps handles all the service methods made available by this package
type Apps struct {
	store store
	totp  *totp.TOTP
}

// Create accepts an App instance and inserts it in the data store. On success it'll return
// the pointer of the app instance which was insterted
func (a *Apps) Create(app App) (*App, error) {
	now := time.Now().UTC()
	app.CreatedAt = &now
	id, err := a.store.Insert(app)
	if err != nil {
		return nil, err
	}
	app.ID = id
	return &app, nil
}

// Update updates the details of the app in the store
func (a *Apps) Update(app App) (*App, error) {
	now := time.Now().UTC()
	app.UpdatedAt = &now
	err := a.store.Update(app.ID, app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Delete deletes the app from the store
func (a *Apps) Delete(app App) (*App, error) {
	err := a.store.Delete(app.ID)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func New() *Apps {
	return &Apps{}
}
