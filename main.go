package main

import (
	"log"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/server"
)

func main() {
	settings := model.ApplicationSettings{
		StoreSettings: model.StoreSettings{
			ShowSQLQueries: false,
		},
	}

	webServer, err := server.NewServer(settings)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal((webServer.Start()))
}
