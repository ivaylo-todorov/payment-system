package model

type StoreSettings struct {
	ShowSQLQueries bool
	DummyDb        bool
}

type ApplicationSettings struct {
	StoreSettings StoreSettings
}
