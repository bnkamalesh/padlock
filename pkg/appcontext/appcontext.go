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
	Logger logger.Logger
	Debug  bool
}

func New(l logger.Logger) *AppContext {
	a := &AppContext{
		Logger: l,
	}
	return a
}

type RequestContext struct {
	StartAt time.Time
	EndAt   time.Time
	Source  string
	Debug   bool
}

func (ac *AppContext) NewReqContext(ctx context.Context, src string) (context.Context, *RequestContext) {
	rctx := RequestContext{
		StartAt: time.Now(),
		Source:  src,
		Debug:   ac.Debug,
	}
	c := context.WithValue(
		ctx,
		reqCtxKey,
		rctx,
	)

	return c, &rctx
}

func ReqContext(ctx context.Context) *RequestContext {
	if ctx == nil {
		ctx = context.Background()
	}

	rc, ok := ctx.Value(reqCtxKey).(RequestContext)
	if !ok {
		return nil
	}

	return &rc
}
