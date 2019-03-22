package http

import (
	"net/http"

	"github.com/bnkamalesh/webgo"
)

func helloworld(w http.ResponseWriter, req *http.Request) {
	webgo.R200(w, "hello world!")
}

func Routes() []*webgo.Route {
	return []*webgo.Route{
		&webgo.Route{
			Name:     "home",
			Pattern:  "/",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{helloworld},
		},
	}
}
