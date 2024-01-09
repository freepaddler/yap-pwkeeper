package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yap-pwkeeper/internal/app/server"
	"yap-pwkeeper/internal/app/server/aaa"
	"yap-pwkeeper/internal/app/server/config"
	"yap-pwkeeper/internal/app/server/documents"
	"yap-pwkeeper/internal/app/server/grpcapi"
	"yap-pwkeeper/internal/app/server/serial"
	"yap-pwkeeper/internal/pkg/grpc/interceptors"
	"yap-pwkeeper/internal/pkg/jwtToken"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/mongodb"
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
		jwtToken.SetTTL(2*time.Hour + 10*time.Second)
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

	// serials setup
	serial.SetSource(db)
	serial.SetBatchSize(10)

	// auth controller
	auth := aaa.New(db)

	// documents controller
	docs := documents.New(db)

	// setup grpc
	gs := grpcapi.New(
		grpcapi.WithAddress(conf.Address),
		grpcapi.WithUnaryInterceptors(interceptors.AuthUnaryServer(auth.Validate, "Docs/")),
		grpcapi.WithStreamInterceptors(interceptors.AuthStreamServer(auth.Validate, "Docs/")),
		grpcapi.WithAuthHandlers(grpcapi.NewAuthHandlers(auth)),
		grpcapi.WithDocsHandlers(grpcapi.NewDocsHandlers(docs)),
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
