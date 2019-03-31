// Package api has all the methods which are exposed as the APIs of this application
// It can be exposed using HTTP or GRPC or any protocol
package api

import (
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
