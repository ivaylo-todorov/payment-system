package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ivaylo-todorov/payment-system/model"
)

const (
	AdminRecordSize    = 3
	MerchantRecordSize = 3
)

func ConvertAdminFromModel(a model.Admin) Admin {
	return Admin{
		Id:          a.Id.String(),
		Name:        a.Name,
		Description: a.Description,
		Email:       a.Email,
	}
}

func ConvertMerchantFromModel(m model.Merchant) Merchant {
	return Merchant{
		Id:                 m.Id.String(),
		Name:               m.Name,
		Description:        m.Description,
		Email:              m.Email,
		Status:             m.Status,
		TransactionsAmount: m.TransactionsAmount,
	}
}

func ConvertMerchantToModel(m Merchant) (model.Merchant, error) {

	var err error
	var id uuid.UUID
	if m.Id != "" {
		id, err = uuid.Parse(m.Id)
		if err != nil {
			return model.Merchant{}, err
		}
	}

	return model.Merchant{
		Id:          id,
		Name:        m.Name,
		Description: m.Description,
		Email:       m.Email,
	}, nil
}

func ConvertTransactionFromModel(t model.Transaction) Transaction {
	return Transaction{
		Id:            t.Id.String(),
		MerchantId:    t.MerchantId.String(),
		Type:          t.Type,
		Amount:        t.Amount,
		Status:        t.Status,
		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
	}
}

func ConvertTransactionToModel(t Transaction) (model.Transaction, error) {

	var err error
	var tid uuid.UUID
	if t.Id != "" {
		tid, err = uuid.Parse(t.Id)
		if err != nil {
			return model.Transaction{}, err
		}
	}

	var mid uuid.UUID
	if t.MerchantId != "" {
		mid, err = uuid.Parse(t.MerchantId)
		if err != nil {
			return model.Transaction{}, err
		}
	}

	return model.Transaction{
		Id:            tid,
		MerchantId:    mid,
		Type:          t.Type,
		Amount:        t.Amount,
		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
	}, nil
}

func ConvertCsvToAdmins(data []byte) ([]model.Admin, error) {
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		return nil, err
	}

	admins := []model.Admin{}

	for _, r := range records {

		if len(r) != MerchantRecordSize {
			return nil, fmt.Errorf("invalid admin csv data")
		}

		admins = append(admins, model.Admin{
			Name:        strings.TrimSpace(r[0]),
			Description: strings.TrimSpace(r[1]),
			Email:       strings.TrimSpace(r[2]),
		})
	}

	return admins, nil
}

func ConvertCsvToMerchants(data []byte) ([]model.Merchant, error) {
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		return nil, err
	}

	merchants := []model.Merchant{}

	for _, r := range records {

		if len(r) != AdminRecordSize {
			return nil, fmt.Errorf("invalid merchant csv data")
		}

		merchants = append(merchants, model.Merchant{
			Name:        strings.TrimSpace(r[0]),
			Description: strings.TrimSpace(r[1]),
			Email:       strings.TrimSpace(r[2]),
		})
	}

	return merchants, nil
}
