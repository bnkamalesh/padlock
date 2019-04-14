package http

import (
	"net/http"

	"github.com/bnkamalesh/webgo"
)

func helloworld(w http.ResponseWriter, req *http.Request) {
	webgo.R200(w, "hello world!")
}

func routes(s *Server) []*webgo.Route {
	return []*webgo.Route{
		&webgo.Route{
			Name:     "home",
			Pattern:  "/",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{helloworld},
		},
		&webgo.Route{
			Name:     "login",
			Pattern:  "/login",
			Method:   http.MethodPost,
			Handlers: []http.HandlerFunc{s.Login},
		},
		&webgo.Route{
			Name:     "auth.required",
			Pattern:  "/restricted",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{s.Authentication, helloworld},
		},
	}
}
