// Package http has all the HTTP handlers for the APIs
package http

import (
	"github.com/bnkamalesh/webgo"
	"github.com/bnkamalesh/webgo/middleware"

	"github.com/bnkamalesh/padlock/api"
	"github.com/bnkamalesh/padlock/pkg/appcontext"
)

type Server struct {
	appDomain string
	appCtx    *appcontext.AppContext
	api       *api.API
	router    *webgo.Router
}

func (s *Server) Start() error {
	s.router.Start()
	return nil
}

func NewServer(host, port, appdomain string, api *api.API, appCtx *appcontext.AppContext) (*Server, error) {
	s := &Server{
		appCtx:    appCtx,
		api:       api,
		appDomain: appdomain,
	}

	router := webgo.NewRouter(&webgo.Config{
		Host: host,
		Port: port,
	}, routes(s))
	webgo.LOGHANDLER = appCtx.Logger
	s.router = router
	// This should be final middleware added, so that execution starts with this
	defer router.Use(s.MiddlewareReqCtx)

	if s.appCtx.Debug {
		router.Use(middleware.AccessLog)
	}

	return s, nil
}
