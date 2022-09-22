package store

import (
	"github.com/google/uuid"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/store/db"
)

type Store interface {
	CreateAdmin(model.Admin) (model.Admin, error)

	CreateMerchant(model.Merchant) (model.Merchant, error)
	UpdateMerchant(model.Merchant) (model.Merchant, error)
	DeleteMerchant(uuid.UUID) error
	GetMerchants(model.MerchantQuery) ([]model.Merchant, error)

	CreateTransaction(model.Transaction) (model.Transaction, error)
	GetTransaction(uuid.UUID) (model.Transaction, error)
	GetTransactions(model.TransactionQuery) ([]model.Transaction, error)
}

func NewStore(settingss model.StoreSettings) (Store, error) {
	return db.NewDb(settingss)
}
