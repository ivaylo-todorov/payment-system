package model

import (
	"errors"

	"github.com/google/uuid"
)

const (
	UserRoleAdmin    = "admin"
	UserRoleMerchant = "merchant"

	MerchantStatusActive   = "active"
	MerchantStatusInactive = "inactive"

	TransactionTypeAuthorize = "authorize"
	TransactionTypeCharge    = "charge"
	TransactionTypeRefund    = "refund"
	TransactionTypeReversal  = "reversal"

	TransactionStatusApproved = "approved"
	TransactionStatusReversed = "reversed"
	TransactionStatusRefunded = "refunded"
	TransactionStatusError    = "error"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrMerchantNotFound    = errors.New("merchant not found")
	ErrTransactionNotFound = errors.New("transaction not found")
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
	TransactionsAmount int64
}

type Transaction struct {
	Id            uuid.UUID
	ParentId      uuid.UUID
	MerchantId    uuid.UUID
	Type          string
	Amount        int64
	Status        string
	CustomerEmail string
	CustomerPhone string
}

type MerchantQuery struct {
}

type TransactionQuery struct {
}
