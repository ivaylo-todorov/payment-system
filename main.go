package main

import (
	"log"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/server"
)

func main() {
	settings := model.ApplicationSettings{}
	webServer := server.NewServer(settings)
	log.Fatal((webServer.Start()))
}
