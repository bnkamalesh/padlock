package appcontext

import "github.com/bnkamalesh/padlock/pkg/platform/logger"

// AppContext holds all the app level context
type AppContext struct {
	Logger logger.Logger
}

func New(l logger.Logger) *AppContext {
	a := &AppContext{
		Logger: l,
	}
	return a
}
