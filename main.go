package main

import (
	"log"
	"time"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/server"
)

func main() {
	settings := model.ApplicationSettings{
		StoreSettings: model.StoreSettings{
			ShowSQLQueries: false,
		},
		TransactionCleanupFrequency: time.Duration(60),
	}

	webServer, err := server.NewServer(settings)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = webServer.StartTransactionsCleanup(settings.TransactionCleanupFrequency)
	if err != nil {
		log.Fatal(err)
	}

	err = webServer.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = webServer.StopTransactionsCleanup()
	if err != nil {
		log.Fatal(err)
	}
}
