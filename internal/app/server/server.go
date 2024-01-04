package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"yap-pwkeeper/internal/app/server/grpcapi"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

type UserController interface {
	Register(credentials models.UserCredentials) (models.User, error)
}

type App struct {
	wg sync.WaitGroup
	gs *grpcapi.GRCPServer
}

// New is a new server instance constructor
func New(options ...func(app *App)) *App {
	app := new(App)
	for _, opt := range options {
		opt(app)
	}
	return app
}

func WithGRPCServer(gs *grpcapi.GRCPServer) func(app *App) {
	return func(app *App) {
		app.gs = gs
	}
}

// Run starts server instance
func (a *App) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("server start terminated: %w", ctx.Err())
	default:
	}

	grpcError := make(chan error)
	a.wg.Add(1)
	go func(stop chan error) {
		defer a.wg.Done()
		defer func() { close(stop) }()
		logger.Log().Infof("starting grpc server at %s", a.gs.GetAddress())
		listener, err := net.Listen("tcp", a.gs.GetAddress())
		if err != nil {
			stop <- err
			return
		}
		if err := a.gs.Server.Serve(listener); err != nil {
			stop <- err
			return
		}
		logger.Log().Info("grpc server stopped")
		return
	}(grpcError)

	// waiting for main context to be cancelled
	select {
	case err := <-grpcError:
		return fmt.Errorf("grpc server error: %w", err)
	case <-ctx.Done():
		logger.Log().Info("stop request received")
	}

	// gracefully stop grpc server
	logger.Log().Info("stopping grpc server")
	grpcStop := make(chan struct{})
	go func() {
		a.gs.Server.GracefulStop()
		close(grpcStop)
	}()
	select {
	case <-grpcStop:
		logger.Log().Info("grpc server graceful shutdown")
	case <-time.After(10 * time.Second):
		logger.Log().Info("grpc server forced shutdown")
	}
	a.gs.Server.Stop()

	// wait until tasks stopped
	a.wg.Wait()

	return nil
}
