package controller

import (
	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/store"
)

type Controller interface {
	CreateAdmins([]model.Admin) ([]model.Admin, error)

	CreateMerchants([]model.Merchant) ([]model.Merchant, error)
	UpdateMerchant(model.Merchant) (model.Merchant, error)
	DeleteMerchant(model.Merchant) error
	GetMerchants() ([]model.Merchant, error)

	StartTransaction(model.Transaction) (model.Transaction, error)
	GetTransactions() ([]model.Transaction, error)
}

func NewController(settings model.ApplicationSettings) *controller {
	return &controller{
		Store: store.NewStore(settings),
	}
}

type controller struct {
	Store store.Store
}

func (m *controller) CreateAdmins(input []model.Admin) ([]model.Admin, error) {
	result := []model.Admin{}
	for _, i := range input {
		a, err := m.Store.CreateAdmin(i)
		if err != nil {
			return result, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (m *controller) CreateMerchants(input []model.Merchant) ([]model.Merchant, error) {
	result := []model.Merchant{}
	for _, i := range input {
		m, err := m.Store.CreateMerchant(i)
		if err != nil {
			return result, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (m *controller) UpdateMerchant(merchant model.Merchant) (model.Merchant, error) {
	return m.Store.UpdateMerchant(merchant)
}

func (m *controller) DeleteMerchant(merchant model.Merchant) error {
	return m.Store.DeleteMerchant(merchant)
}

func (m *controller) GetMerchants() ([]model.Merchant, error) {
	return m.Store.GetMerchants()
}

func (m *controller) StartTransaction(transaction model.Transaction) (model.Transaction, error) {
	return m.Store.StartTransaction(transaction)
}

func (m *controller) GetTransactions() ([]model.Transaction, error) {
	return m.Store.GetTransactions()
}
