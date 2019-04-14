package appcontext

import (
	"context"
	"time"

	"github.com/bnkamalesh/padlock/pkg/platform/logger"
)

type ctxKey string

const (
	reqCtxKey = ctxKey("request")
)

// AppContext holds all the app level context
type AppContext struct {
	Logger  logger.Logger
	Logging bool
	Debug   bool
}

func New(l logger.Logger) *AppContext {
	a := &AppContext{
		Logger: l,
	}
	return a
}

type RequestContext struct {
	StartAt *time.Time `json:"startAt,omitempty"`
	EndAt   *time.Time `json:"endAt,omitempty"`
	Source  string     `json:"source,omitempty"`
	Debug   bool       `json:"debug,omitempty"`
}

func (ac *AppContext) NewReqContext(ctx context.Context, source string) (context.Context, *RequestContext) {
	rctx := &RequestContext{
		Source: source,
	}

	now := time.Now()
	if rctx.StartAt == nil {
		rctx.StartAt = &now
	}
	c := context.WithValue(
		ctx,
		reqCtxKey,
		rctx,
	)

	return c, rctx
}

func (ac *AppContext) ReqContext(ctx context.Context) *RequestContext {
	if ctx == nil {
		ctx = context.Background()
	}

	rc, ok := ctx.Value(reqCtxKey).(*RequestContext)
	if !ok {
		return nil
	}

	return rc
}
