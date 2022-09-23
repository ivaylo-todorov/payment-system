package model

import "time"

type StoreSettings struct {
	ShowSQLQueries bool
	DummyDb        bool
}

type ApplicationSettings struct {
	StoreSettings               StoreSettings
	TransactionCleanupFrequency time.Duration // in minutes
}
