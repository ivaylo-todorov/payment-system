package controller

import (
	"fmt"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/store"
)

type Controller interface {
	CreateAdmins([]model.Admin) ([]model.Admin, error)

	CreateMerchants([]model.Merchant) ([]model.Merchant, error)
	UpdateMerchant(model.Merchant) (model.Merchant, error)
	DeleteMerchant(model.Merchant) error
	GetMerchants(model.MerchantQuery) ([]model.Merchant, error)

	StartTransaction(model.Transaction) (model.Transaction, error)
	GetTransactions(model.TransactionQuery) ([]model.Transaction, error)
	DeleteTransactions(model.TransactionQuery) error
}

func NewController(settings model.ApplicationSettings, store store.Store) (*controller, error) {
	return &controller{
		Store: store,
	}, nil
}

type controller struct {
	Store store.Store
}

func (c *controller) CreateAdmins(input []model.Admin) ([]model.Admin, error) {
	result := []model.Admin{}
	for _, i := range input {
		if err := model.ValidateAdminCreate(i); err != nil {
			return result, err
		}

		a, err := c.Store.CreateAdmin(i)
		if err != nil {
			return result, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (c *controller) CreateMerchants(input []model.Merchant) ([]model.Merchant, error) {
	result := []model.Merchant{}
	for _, i := range input {
		if err := model.ValidateMerchantCreate(i); err != nil {
			return result, err
		}

		m, err := c.Store.CreateMerchant(i)
		if err != nil {
			return result, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (c *controller) UpdateMerchant(merchant model.Merchant) (model.Merchant, error) {
	if err := model.ValidateMerchantUpdate(merchant); err != nil {
		return merchant, err
	}

	return c.Store.UpdateMerchant(merchant)
}

func (c *controller) DeleteMerchant(merchant model.Merchant) error {
	if err := model.ValidateMerchantDelete(merchant); err != nil {
		return err
	}

	return c.Store.DeleteMerchant(merchant.Id)
}

func (c *controller) GetMerchants(query model.MerchantQuery) ([]model.Merchant, error) {
	return c.Store.GetMerchants(query)
}

func (c *controller) StartTransaction(transaction model.Transaction) (model.Transaction, error) {
	if err := model.ValidateTransactionCreate(transaction); err != nil {
		return transaction, err
	}

	if transaction.Type == model.TransactionTypeAuthorize {
		transaction.Status = model.TransactionStatusApproved
		return c.Store.CreateTransaction(transaction)
	}

	parent, err := c.Store.GetTransaction(transaction.ParentId)
	if err != nil {
		return transaction, err
	}

	if parent.Status != model.TransactionStatusApproved && parent.Status != model.TransactionStatusRefunded {
		transaction.Status = model.TransactionStatusError
		return c.Store.CreateTransaction(transaction)
	}

	if transaction.Type == model.TransactionTypeCharge {
		if parent.Type != model.TransactionTypeAuthorize {
			return transaction, fmt.Errorf("invalid reference transaction type, %s", parent.Type)
		}

		if parent.Status != model.TransactionStatusApproved {
			return transaction, fmt.Errorf("invalid reference transaction status, %s", parent.Status)
		}

		if transaction.Amount > parent.Amount {
			return transaction, fmt.Errorf("transaction amount bigger than authorized")
		}

		transaction.Status = model.TransactionStatusApproved
		return c.Store.CreateTransaction(transaction)
	}

	if transaction.Type == model.TransactionTypeRefund {
		if parent.Type != model.TransactionTypeCharge {
			return transaction, fmt.Errorf("invalid reference transaction type, %s", parent.Type)
		}

		if parent.Status != model.TransactionStatusApproved {
			return transaction, fmt.Errorf("invalid reference transaction status, %s", parent.Status)
		}

		if transaction.Amount != parent.Amount {
			return transaction, fmt.Errorf("transaction amount different than charged")
		}

		transaction.Status = model.TransactionStatusRefunded
		return c.Store.CreateTransaction(transaction)
	}

	if transaction.Type == model.TransactionTypeReversal {
		if parent.Type != model.TransactionTypeAuthorize {
			return transaction, fmt.Errorf("invalid reference transaction type, %s", parent.Type)
		}

		if parent.Status != model.TransactionStatusApproved {
			return transaction, fmt.Errorf("invalid reference transaction status, %s", parent.Status)
		}

		transaction.Status = model.TransactionStatusReversed
		return c.Store.CreateTransaction(transaction)
	}

	return transaction, fmt.Errorf("invalid transaction type")
}

func (c *controller) GetTransactions(query model.TransactionQuery) ([]model.Transaction, error) {
	return c.Store.GetTransactions(query)
}

func (c *controller) DeleteTransactions(query model.TransactionQuery) error {
	return c.Store.DeleteTransactions(query)
}
