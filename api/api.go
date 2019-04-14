// Package api has all the methods which are exposed as the APIs of this application
// It can be exposed using HTTP or GRPC or any protocol
package api

import (
	"context"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
	"github.com/bnkamalesh/padlock/pkg/apps"
	"github.com/bnkamalesh/padlock/pkg/users"
)

type API struct {
	appCtx *appcontext.AppContext
	apps   *apps.Apps
	users  *users.Users
}

func New(appCtx *appcontext.AppContext, a *apps.Apps, u *users.Users) *API {
	api := &API{
		appCtx: appCtx,
		apps:   a,
		users:  u,
	}

	return api
}

// AuthenticatedUser gets the details of the user based on the ource & auth token provided
func (a *API) Login(ctx context.Context, source, email, password string) (*users.User, string, error) {
	return a.users.Login(ctx, source, email, password)
}

// AuthenticatedUser gets the details of the user based on the ource & auth token provided
func (a *API) AuthenticatedUser(ctx context.Context, source, token string) (*users.User, error) {
	return a.users.AuthUser(ctx, source, token)
}
