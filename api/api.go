// Package api has all the methods which are exposed as the APIs of this application
// It can be exposed using HTTP or GRPC or any protocol
package api

import (
	"github.com/bnkamalesh/padlock/pkg/appcontext"
)

type API struct {
	appCtx *appcontext.AppContext
}
