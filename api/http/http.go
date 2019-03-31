// Package http has all the HTTP handlers for the APIs
package http

import (
	"github.com/bnkamalesh/webgo"

	"github.com/bnkamalesh/padlock/api"
)

type Server struct {
	router *webgo.Router
	api    *api.API
}

func (s *Server) Start() error {
	s.router.Start()
	return nil
}

func NewServer(host, port string, api *api.API) (*Server, error) {
	s := &Server{}
	router := webgo.NewRouter(&webgo.Config{
		Host: host,
		Port: port,
	}, routes())
	s.router = router
	return s, nil
}
