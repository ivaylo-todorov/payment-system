package model

import (
	"github.com/google/uuid"
)

type Admin struct {
	Id          uuid.UUID
	Name        string
	Description string
	Email       string
}

type Merchant struct {
	Id          uuid.UUID
	Name        string
	Description string
	Email       string

	Status             string
	TransactionsAmount uint64
}

type Transaction struct {
	Id            uuid.UUID
	MerchantId    uuid.UUID
	Amount        uint64
	Status        string
	CustomerEmail string
	CustomerPhone string
}
