package main

import (
	"log"

	"github.com/ivaylo-todorov/payment-system/server"
)

func main() {
	webServer := server.NewServer()
	log.Fatal((webServer.Start()))
}
