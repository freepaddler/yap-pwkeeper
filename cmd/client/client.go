package main

import (
	"fmt"
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
	if conf.Debug == 2 {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}

	f, err := os.Create("/tmp/pwkeeper.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	// setup grpc client
	log.Println("setup server connection...")
	grpcClient, err := grpccli.New(conf.Address)
	if err != nil {
		log.Printf("server connection setup failed: %s", err.Error())
		exitCode = 1
		return
	}
	defer func() { _ = grpcClient.Close() }()

	_ = grpcClient.Login("chu", "Victor")

	store := memstore.New(grpcClient)
	//err = store.Update()
	//if err != nil {
	//	log.Println(err)
	//}

	ui := client.New(
		client.WithAuthServer(grpcClient),
		client.WithDataStore(store),
		//client.WithDebug(conf.Debug),
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
