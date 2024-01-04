package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yap-pwkeeper/internal/app/server"
	"yap-pwkeeper/internal/app/server/config"
	"yap-pwkeeper/internal/app/server/grpcapi"
	"yap-pwkeeper/internal/pkg/aaa"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/mongodb"
	"yap-pwkeeper/pkg/jwtToken"
)

var (
	// go build -ldflags " \
	// -X 'main.buildVersion=$(git describe --tag --always 2>/dev/null)' \
	// -X 'main.buildDate=$(date)' \
	// "
	buildVersion, buildDate = "N/A", "N/A"
)

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()
	// print version
	version()

	// get config
	conf := config.New()

	// version flag
	if conf.Version {
		return
	}

	// setup logging
	if conf.Debug {
		logger.SetMode(logger.ModeDev)
		logger.SetLevel(-1)
		conf.Print()
	} else {
		logger.SetMode(logger.ModeProd)
		logger.SetLevel(conf.LogLevel)
	}

	// set jwt key
	if conf.TokenKey != "" {
		jwtToken.SetKey(conf.TokenKey)
	}

	logger.Log().Info("starting server")
	defer func() { logger.Log().Info("server stopped") }()
	// notify context
	nCtx, nStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer nStop()

	// connect database
	logger.Log().Info("connecting database")
	db, err := mongodb.New(nCtx, conf.DbUri)
	if err != nil {
		logger.Log().WithErr(err).Error("database setup failed")
		exitCode = 1
		return
	}
	defer func() {
		logger.Log().Info("closing database connection")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := db.Close(ctx)
		if err != nil {
			logger.Log().WithErr(err).Warn("database connection terminated")
			exitCode = 2
		} else {
			logger.Log().Info("database connection closed gracefully")
		}
		cancel()
	}()

	// auth controller
	auth := aaa.New(db)

	// setup grpc
	gs := grpcapi.New(
		grpcapi.WithAddress(conf.Address),
		grpcapi.WithAuthHandlers(grpcapi.NewAuthHandlers(auth)),
	)

	// init and run server
	serverApp := server.New(
		server.WithGRPCServer(gs),
	)
	err = serverApp.Run(nCtx)
	if err != nil {
		logger.Log().WithErr(err).Error("unclean exit")
		exitCode = 2
	}
}

func version() {
	_, _ = fmt.Fprintf(
		os.Stdout,
		`Build version: %s
Build date: %s
`, buildVersion, buildDate)
}
