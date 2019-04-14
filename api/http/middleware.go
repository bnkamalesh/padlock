package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bnkamalesh/padlock/pkg/users"

	"github.com/OneOfOne/xxhash"
	"github.com/bnkamalesh/webgo"
)

func sourceID(r *http.Request) string {
	h := xxhash.New64()
	rdr := strings.NewReader(r.RemoteAddr + r.Header.Get("User-Agent"))
	io.Copy(h, rdr)
	return fmt.Sprintf("%d", h.Sum64())
}

type responseWriter struct {
	StartAt time.Time
	http.ResponseWriter
	code int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
	// server timing duration in milliseconds
	rw.Header().Set("Server-Timing", fmt.Sprintf("total;dur=%v", time.Since(rw.StartAt).Seconds()/1000))
	rw.ResponseWriter.WriteHeader(code)
}

// MiddlewareReqCtx injects appcontext.RequestContext into HTTP request context
func (s *Server) MiddlewareReqCtx(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	ctx, rctx := s.appCtx.NewReqContext(r.Context(), sourceID(r))
	r = r.WithContext(
		ctx,
	)

	w = &responseWriter{
		*rctx.StartAt,
		w,
		0,
	}

	next(w, r)
}

func (s *Server) Authentication(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			if s.appCtx.Debug && s.appCtx.Logging {
				s.appCtx.Logger.Error(err)
			}
		} else {
			if cookie.Expires.Before(time.Now()) {
				webgo.R403(w, "Sorry, you should be authorized to access this page")
				return
			}
			token = cookie.String()
		}
	}

	if token == "" {
		webgo.R403(w, "Sorry, you should be authorized to access this page")
		return
	}

	ctx := r.Context()
	u, err := s.api.AuthenticatedUser(ctx, sourceID(r), token)
	if err != nil {
		switch err {
		case users.ErrSessionID, users.ErrSessionIDExpired:
			{
				webgo.R403(w, err.Error())
			}
		default:
			{
				webgo.R500(w, err.Error())
			}
		}
		return
	}

	r = r.WithContext(
		users.SetContext(ctx, u),
	)
}
