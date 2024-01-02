package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/server"
	"yap-pwkeeper/internal/server/config"
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

	// notify context
	nCtx, nStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer nStop()

	logger.Log().Info("starting server")
	serverApp := server.New()
	err := serverApp.Run(nCtx)
	if err != nil {
		logger.Log().WithErr(err).Error("unclean exit")
		exitCode = 2
	}
	logger.Log().Info("server stopped")

}

func version() {
	_, _ = fmt.Fprintf(
		os.Stdout,
		`Build version: %s
Build date: %s
`, buildVersion, buildDate)
}
