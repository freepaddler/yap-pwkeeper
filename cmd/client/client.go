package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/credentials/insecure"

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

	// print version
	version()

	// get config
	conf := config.New()

	// version flag
	if conf.Version {
		return
	}

	defer func() {
		//log.SetOutput(os.Stderr)
		log.Printf("application exited with code: %d\n", exitCode)
		os.Exit(exitCode)
	}()

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

	// enable tls connection to server
	tlsCredentials := insecure.NewCredentials()
	if conf.TlsCaCertFile != "" {
		var err error
		tlsCredentials, err = grpccli.LoadCACertificate(conf.TlsCaCertFile, conf.TlsInsecure)
		if err != nil {
			log.Printf("unable to load CA certificate: %s", err)
			return
		}
	}

	// setup grpc client
	log.Println("setup server connection...")
	grpcClient, err := grpccli.New(conf.Address,
		grpccli.WithTransportCredentials(tlsCredentials),
		grpccli.WithTimeouts(5*time.Second, 30*time.Second),
		grpccli.WithTokenRefresh(2*time.Minute, 5*time.Second),
	)
	if err != nil {
		log.Printf("server connection setup failed: %s", err.Error())
		exitCode = 1
		return
	}
	defer func() { _ = grpcClient.Close() }()

	store := memstore.New(grpcClient)

	ui := client.New(
		client.WithDataStore(store),
		client.WithMouse(conf.UseMouse),
	)
	log.Println("starting ui")

	// disable log output
	if conf.Logfile == "" {
		log.SetOutput(io.Discard)
	}

	if err := ui.Run(); err != nil {
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
