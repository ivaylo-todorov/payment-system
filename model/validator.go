package model

import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

func ValidateAdminCreate(a Admin) error {
	if a.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if err := validateEmailString(a.Email); err != nil {
		return err
	}
	return nil
}

func ValidateMerchantCreate(m Merchant) error {
	if m.Status != MerchantStatusActive && m.Status != MerchantStatusInactive {
		return fmt.Errorf("invalid merchant status")
	}
	if m.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if err := validateEmailString(m.Email); err != nil {
		return err
	}
	return nil
}

func ValidateMerchantUpdate(m Merchant) error {
	if m.Id == uuid.Nil {
		return fmt.Errorf("missing merchant id")
	}

	if m.Status != "" {
		if m.Status != MerchantStatusActive && m.Status != MerchantStatusInactive {
			return fmt.Errorf("invalid merchant status")
		}
	}

	if m.Email != "" {
		if err := validateEmailString(m.Email); err != nil {
			return err
		}
	}
	return nil
}

func ValidateMerchantDelete(m Merchant) error {
	if m.Id == uuid.Nil {
		return fmt.Errorf("missing merchant id")
	}
	return nil
}

func ValidateTransactionCreate(t Transaction) error {
	if t.MerchantId == uuid.Nil {
		return fmt.Errorf("missing merchant id")
	}
	if err := validateEmailString(t.CustomerEmail); err != nil {
		return err
	}
	if t.Status != "" {
		return fmt.Errorf("transaction status must be empty")
	}

	if t.Type == TransactionTypeAuthorize {
		if t.Amount == 0 {
			return fmt.Errorf("zero transaction amount")
		}
		return nil
	}

	if t.Type == TransactionTypeCharge || t.Type == TransactionTypeRefund {
		if t.ParentId == uuid.Nil {
			return fmt.Errorf("missing reference transaction id")
		}
		if t.Amount == 0 {
			return fmt.Errorf("zeero transaction amount")
		}
		return nil
	}

	if t.Type == TransactionTypeReversal {
		if t.ParentId == uuid.Nil {
			return fmt.Errorf("missing reference transaction id")
		}
		if t.Amount == 0 {
			return fmt.Errorf("transaction amount should be zero")
		}
		return nil
	}

	return fmt.Errorf("invalid transaction type, %s", t.Type)
}

func validateEmailString(address string) error {
	_, err := mail.ParseAddress(address)
	return err
}
