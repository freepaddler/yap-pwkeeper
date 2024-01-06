package main

import (
	"fmt"
	"os"
	"time"

	"yap-pwkeeper/internal/app/client/auth"
	"yap-pwkeeper/internal/app/client/config"
	"yap-pwkeeper/internal/app/client/grpccli"
	"yap-pwkeeper/internal/pkg/logger"
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
	logger.SetMode(logger.ModeDev)
	logger.SetLevel(conf.LogLevel)
	conf.Print()

	// setup grpc client
	client, err := grpccli.New(conf.Address)
	if err != nil {
		logger.Log().WithErr(err).Error("setup server connection failed")
	}
	defer func() { _ = client.Close() }()
	logger.Log().Info("server connection set up")

	// setup auth
	aaa := auth.New(client)

	err = aaa.Login("login", "password")
	if err != nil {
		logger.Log().WithErr(err).Debug("login failed")
	} else {
		logger.Log().Debug("logged in")
	}
	time.Sleep(40 * time.Second)
	return

	//ui := client.New()
	//log.Print("start")
	//if err := ui.Run(); err != nil {
	//	panic(err)
	//}
	//log.Print("end")
}

func version() {
	_, _ = fmt.Fprintf(
		os.Stdout,
		`Build version: %s
Build date: %s
`, buildVersion, buildDate)
}
