package store

import (
	"github.com/google/uuid"

	"github.com/ivaylo-todorov/payment-system/model"
)

type Store interface {
	CreateAdmin(model.Admin) (model.Admin, error)

	CreateMerchant(model.Merchant) (model.Merchant, error)
	UpdateMerchant(model.Merchant) (model.Merchant, error)
	DeleteMerchant(model.Merchant) error
	GetMerchants() ([]model.Merchant, error)

	StartTransaction(model.Transaction) (model.Transaction, error)
	GetTransactions() ([]model.Transaction, error)
}

func NewStore(settings model.ApplicationSettings) Store {
	return &dummyStore{}
}

type dummyStore struct {
}

func (s *dummyStore) CreateAdmin(a model.Admin) (model.Admin, error) {
	return model.Admin{
		Id:          uuid.New(),
		Name:        a.Name,
		Description: a.Description,
		Email:       a.Email,
	}, nil
}

func (s *dummyStore) CreateMerchant(m model.Merchant) (model.Merchant, error) {
	return model.Merchant{
		Id:          uuid.New(),
		Name:        m.Name,
		Description: m.Description,
		Email:       m.Email,
		Status:      "active",
	}, nil
}

func (s *dummyStore) UpdateMerchant(m model.Merchant) (model.Merchant, error) {
	return model.Merchant{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Email:       m.Email,
	}, nil
}

func (s *dummyStore) DeleteMerchant(model.Merchant) error {
	return nil
}

func (s *dummyStore) GetMerchants() ([]model.Merchant, error) {
	return []model.Merchant{
		{
			Id:   uuid.New(),
			Name: "Merchant One",
		},
		{
			Id:   uuid.New(),
			Name: "Merchant Two",
		},
	}, nil
}

func (s *dummyStore) StartTransaction(transaction model.Transaction) (model.Transaction, error) {
	transaction.Id = uuid.New()
	transaction.MerchantId = uuid.New()
	return transaction, nil
}

func (s *dummyStore) GetTransactions() ([]model.Transaction, error) {
	return []model.Transaction{
		{
			Id:         uuid.New(),
			MerchantId: uuid.New(),
			Status:     "approved",
			Amount:     100,
		},
		{
			Id:         uuid.New(),
			MerchantId: uuid.New(),
			Status:     "refunded",
			Amount:     200,
		},
	}, nil
}
