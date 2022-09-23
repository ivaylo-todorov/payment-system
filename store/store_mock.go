package store

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/ivaylo-todorov/payment-system/model"
)

var (
	AuthorizeTransactionOneUuid = uuid.MustParse("43367211-6797-48f0-a07f-e8ab4b640a60")
	ChargeTransactionUuid       = uuid.MustParse("1111feb5-0980-4f04-b73b-cc02ba3d5706")
	RefundTransactionUuid       = uuid.MustParse("3713af41-3018-4541-9d10-1e5b97fa5023")

	AuthorizeTransactionTwoUuid = uuid.MustParse("3c6729ac-5467-4f8b-88c4-acd7088ae34d")
	ReversalTransactionUuid     = uuid.MustParse("589d04ab-8a07-4f51-ad9c-5ba67ff0c1b9")

	AdminUuid       = uuid.MustParse("d45f5664-0bc1-44a9-9a43-fddb2462d3c3")
	MerchantOneUuid = uuid.MustParse("c15760c1-bb8d-4717-98f9-feb182950259")
	MerchantTwoUuid = uuid.MustParse("78e64a3b-ffb5-42be-a491-0d26bc73b3b5")
)

var adminMock = model.Admin{
	Id: AdminUuid,
}

var merchantMock = map[uuid.UUID]model.Merchant{
	MerchantOneUuid: {
		Id: MerchantOneUuid,
	},
	MerchantTwoUuid: {
		Id: MerchantTwoUuid,
	},
}

var transactionMock = map[uuid.UUID]model.Transaction{
	AuthorizeTransactionOneUuid: {
		Id: AuthorizeTransactionOneUuid,
	},
	ChargeTransactionUuid: {
		Id: ChargeTransactionUuid,
	},
	RefundTransactionUuid: {
		Id: RefundTransactionUuid,
	},
	AuthorizeTransactionTwoUuid: {
		Id: AuthorizeTransactionTwoUuid,
	},
	ReversalTransactionUuid: {
		Id: ReversalTransactionUuid,
	},
}

var createdMerchants = []model.Merchant{}
var createdTransactions = []model.Transaction{}

// TODO: use github.com/stretchr/testify/mock

func NewMockStore() (*mockStore, error) {
	return &mockStore{}, nil
}

type mockStore struct {
}

func (s *mockStore) CreateAdmin(a model.Admin) (model.Admin, error) {
	adminMock.Name = a.Name
	adminMock.Description = a.Description
	adminMock.Email = a.Email
	return adminMock, nil
}

func (s *mockStore) CreateMerchant(m model.Merchant) (model.Merchant, error) {
	if m.Id == uuid.Nil {
		return model.Merchant{}, fmt.Errorf("merchant id is required for testing")
	}

	if _, ok := merchantMock[m.Id]; !ok {
		return model.Merchant{}, fmt.Errorf("only mock merchant ids allowed")
	}

	merchantMock[m.Id] = model.Merchant{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Email:       m.Email,
		Status:      m.Status,
	}

	createdMerchants = append(createdMerchants, merchantMock[m.Id])

	return merchantMock[m.Id], nil
}

func (s *mockStore) UpdateMerchant(m model.Merchant) (model.Merchant, error) {

	if _, ok := merchantMock[m.Id]; !ok {
		return model.Merchant{}, fmt.Errorf("only mock merchant ids allowed")
	}

	name := merchantMock[m.Id].Name
	if m.Name != "" {
		name = m.Name
	}
	description := merchantMock[m.Id].Description
	if m.Description != "" {
		description = m.Description
	}
	email := merchantMock[m.Id].Email
	if m.Email != "" {
		email = m.Email
	}
	status := merchantMock[m.Id].Status
	if m.Status != "" {
		status = m.Status
	}

	merchantMock[m.Id] = model.Merchant{
		Id:          m.Id,
		Name:        name,
		Description: description,
		Email:       email,
		Status:      status,
	}

	return merchantMock[m.Id], nil
}

func (s *mockStore) DeleteMerchant(uuid.UUID) error {
	return nil
}

func (s *mockStore) GetMerchants(model.MerchantQuery) ([]model.Merchant, error) {
	result := []model.Merchant{}

	for _, m := range createdMerchants {
		result = append(result, m)
	}

	return result, nil
}

func (s *mockStore) CreateTransaction(t model.Transaction) (model.Transaction, error) {
	if t.Id == uuid.Nil {
		return model.Transaction{}, fmt.Errorf("transaction id is required for testing")
	}

	if _, ok := transactionMock[t.Id]; !ok {
		return model.Transaction{}, fmt.Errorf("only mock transaction ids allowed")
	}

	transactionMock[t.Id] = model.Transaction{
		Id:            t.Id,
		ParentId:      t.ParentId,
		MerchantId:    t.MerchantId,
		Type:          t.Type,
		Amount:        t.Amount,
		Status:        t.Status,
		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
	}

	createdTransactions = append(createdTransactions, transactionMock[t.Id])

	return transactionMock[t.Id], nil
}

func (s *mockStore) GetTransaction(id uuid.UUID) (model.Transaction, error) {
	for _, t := range createdTransactions {
		if t.Id == id {
			return t, nil
		}
	}
	return model.Transaction{}, fmt.Errorf("transaction not found")
}

func (s *mockStore) GetTransactions(model.TransactionQuery) ([]model.Transaction, error) {
	result := []model.Transaction{}

	for _, t := range createdTransactions {
		result = append(result, t)
	}

	return result, nil
}

func (s *mockStore) DeleteTransactions(query model.TransactionQuery) error {
	return nil
}
