package http

import (
	"encoding/json"
	"net/http"

	"github.com/bnkamalesh/webgo"
)

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	rctx := s.appCtx.ReqContext(ctx)
	source := ""
	if rctx != nil {
		source = rctx.Source
	}

	payload := make(map[string]string)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		webgo.R400(w, err)
		return
	}

	u, token, err := s.api.Login(ctx, source, payload["email"], payload["password"])
	if err != nil {
		webgo.R400(w, err.Error())
		return
	}

	w.Header().Set("Authorization", token)
	http.SetCookie(
		w,
		&http.Cookie{
			Name:  "Authorization",
			Value: token,
			// seconds for 1 week
			MaxAge:   60 * 60 * 24 * 7,
			HttpOnly: true,
			Path:     "/",
		},
	)
	webgo.R200(w, u)
}
