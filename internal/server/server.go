package server

import (
	"context"
	"fmt"

	"yap-pwkeeper/internal/pkg/logger"
)

type App struct {
}

// New is a new server instance constructor
func New(options ...func(app *App)) *App {
	app := new(App)
	for _, opt := range options {
		opt(app)
	}
	return app
}

// Run starts server instance
func (a *App) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("server start terminated: %w", ctx.Err())
	default:
	}
	logger.Log().Info("run called")
	return nil
}
