package main

import (
	"log"

	"yap-pwkeeper/internal/app/client"
)

func main() {
	ui := client.New()
	log.Print("start")
	if err := ui.Run(); err != nil {
		panic(err)
	}
	log.Print("end")
}
