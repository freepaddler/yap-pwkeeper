package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"yap-pwkeeper/internal/app/client"
	"yap-pwkeeper/internal/app/client/config"
	"yap-pwkeeper/internal/app/client/grpccli"
	"yap-pwkeeper/internal/app/client/memstore"
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
	defer func() {
		//log.SetOutput(os.Stderr)
		log.Printf("application exited with code: %d\n", exitCode)
		os.Exit(exitCode)
	}()
	// print version
	version()

	// get config
	conf := config.New()

	// version flag
	if conf.Version {
		return
	}

	// setup logging
	if conf.Log {
		if conf.Logfile != "" {
			f, err := os.Create(conf.Logfile)
			if err != nil {
				log.Printf("unable to open log file %s", conf.Logfile)
			}
			defer func() { _ = f.Close() }()
			log.SetOutput(f)
		}
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(io.Discard)
	}

	// setup grpc client
	log.Println("setup server connection...")
	grpcClient, err := grpccli.New(conf.Address)
	if err != nil {
		log.Printf("server connection setup failed: %s", err.Error())
		exitCode = 1
		return
	}
	defer func() { _ = grpcClient.Close() }()

	store := memstore.New(grpcClient)

	ui := client.New(
		client.WithDataStore(store),
	)
	log.Println("starting ui")

	if err := ui.Run(); err != nil {
		//log.SetOutput(os.Stderr)
		log.Println("ui terminated")
	}

}

func version() {
	_, _ = fmt.Fprintf(
		os.Stdout,
		`Build version: %s
Build date: %s
`, buildVersion, buildDate)
}
