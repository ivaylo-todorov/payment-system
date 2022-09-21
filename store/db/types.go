package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string
	Password string
	Role     string

	Name        string
	Description string
	Email       string `gorm:"unique;not null"`
}

type Admin struct {
	gorm.Model

	// foreigh key
	UserID uint
	User   User

	AdminId uuid.UUID `gorm:"type:uuid"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.AdminId == uuid.Nil {
		a.AdminId = uuid.New()
	}
	return nil
}

type Merchant struct {
	gorm.Model

	// foreigh key
	UserID uint
	User   User

	MerchantId uuid.UUID `gorm:"type:uuid"`
	Status     string
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.MerchantId == uuid.Nil {
		m.MerchantId = uuid.New()
	}
	return nil
}

type Transaction struct {
	gorm.Model

	// foreigh key
	MerchantID uint
	Merchant   Merchant

	TransactionId uuid.UUID `gorm:"type:uuid"`
	ParentId      uuid.UUID `gorm:"type:uuid"`
	Type          string
	Amount        uint64
	Status        string
	CustomerEmail string
	CustomerPhone string
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.TransactionId == uuid.Nil {
		t.TransactionId = uuid.New()
	}
	return nil
}
