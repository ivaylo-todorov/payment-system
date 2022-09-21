package db

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ivaylo-todorov/payment-system/model"
)

func NewDb(settings model.StoreSettings) (*sqLiteDb, error) {
	gormConfig := &gorm.Config{}

	if settings.ShowSQLQueries {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), gormConfig)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(Admin{}, &Merchant{}, &User{}, &Transaction{})
	if err != nil {
		return nil, err
	}

	res := db.Exec("PRAGMA foreign_keys = ON", nil)
	if res.Error != nil {
		return nil, err
	}

	return &sqLiteDb{
		db: db,
	}, nil
}

type sqLiteDb struct {
	db *gorm.DB
}

func (s *sqLiteDb) CreateAdmin(a model.Admin) (model.Admin, error) {
	txFunc := func(tx *gorm.DB) error {

		user := User{
			Role:        model.UserRoleAdmin,
			Name:        a.Name,
			Description: a.Description,
			Email:       a.Email,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		admin := Admin{
			UserID: user.ID,
		}

		if err := tx.Create(&admin).Error; err != nil {
			return err
		}

		a.Id = admin.AdminId

		return nil
	}

	return a, s.db.Transaction(txFunc)
}

func (s *sqLiteDb) CreateMerchant(m model.Merchant) (model.Merchant, error) {
	txFunc := func(tx *gorm.DB) error {

		user := User{
			Role:        model.UserRoleMerchant,
			Name:        m.Name,
			Description: m.Description,
			Email:       m.Email,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		merchant := Merchant{
			UserID: user.ID,
			Status: model.MerchantStatusActive,
		}

		if err := tx.Create(&merchant).Error; err != nil {
			return err
		}

		m.Id = merchant.MerchantId
		m.Status = merchant.Status

		return nil
	}

	return m, s.db.Transaction(txFunc)
}

func (s *sqLiteDb) UpdateMerchant(m model.Merchant) (model.Merchant, error) {
	merchant := Merchant{}

	result := s.db.Where("merchant_id = ?", m.Id.String()).First(&merchant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Merchant{}, model.ErrMerchantNotFound
		}
		return model.Merchant{}, result.Error
	}

	txFunc := func(tx *gorm.DB) error {

		user := User{
			Model:       gorm.Model{ID: merchant.UserID},
			Name:        m.Name,
			Description: m.Description,
			Email:       m.Email,
		}

		userColumns := []string{}
		if m.Name != "" {
			userColumns = append(userColumns, "Name")
		}
		if m.Description != "" {
			userColumns = append(userColumns, "Description")
		}
		if m.Email != "" {
			userColumns = append(userColumns, "Email")
		}

		if err := tx.Model(&user).Select(userColumns).Updates(user).Error; err != nil {
			return err
		}

		merchant.Status = m.Status

		merchantColumns := []string{}
		if m.Status != "" {
			merchantColumns = append(merchantColumns, "Status")
		}

		if err := tx.Model(&merchant).Select(merchantColumns).Updates(&merchant).Error; err != nil {
			return err
		}

		return nil
	}

	return m, s.db.Transaction(txFunc)
}

func (s *sqLiteDb) DeleteMerchant(m model.Merchant) error {
	merchant := Merchant{}

	result := s.db.Where("merchant_id = ?", m.Id.String()).First(&merchant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.ErrMerchantNotFound
		}
		return result.Error
	}

	return s.db.Transaction(func(tx *gorm.DB) error {

		user := User{
			Model: gorm.Model{ID: merchant.UserID},
		}

		if err := tx.Delete(&user).Error; err != nil {
			return err
		}

		if err := tx.Delete(&merchant).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *sqLiteDb) GetMerchants(query model.MerchantQuery) ([]model.Merchant, error) {

	rows, err := s.db.Model(&Merchant{}).Joins("User").
		Joins(`left join (select merchant_id, sum(amount) as total_transaction_sum from transactions where transactions.type = "charge" and transactions.status = "approved" group by merchant_id) t on merchants.id = t.merchant_id`).
		Select("merchants.merchant_id, merchants.status, t.total_transaction_sum").Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	merchants := []model.Merchant{}

	var m struct {
		Merchant
		TotalTransactionSum uint64
	}

	for rows.Next() {
		err = s.db.ScanRows(rows, &m)
		if err != nil {
			return merchants, err
		}

		merchants = append(merchants, model.Merchant{
			Id:                 m.MerchantId,
			Name:               m.User.Name,
			Description:        m.User.Description,
			Email:              m.User.Email,
			Status:             m.Status,
			TransactionsAmount: m.TotalTransactionSum,
		})
	}

	return merchants, nil
}

func (s *sqLiteDb) CreateTransaction(t model.Transaction) (model.Transaction, error) {

	merchant := Merchant{}

	result := s.db.Where("merchant_id = ?", t.MerchantId.String()).First(&merchant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Transaction{}, model.ErrMerchantNotFound
		}
		return model.Transaction{}, result.Error
	}

	if t.Type == model.TransactionTypeAuthorize {
		return s.createAuthorizeTransaction(merchant.ID, t)
	} else if t.Type == model.TransactionTypeCharge {
		return s.createChargeTransaction(merchant.ID, t)
	} else if t.Type == model.TransactionTypeRefund {
		return s.createRefundTransaction(merchant.ID, t)
	} else if t.Type == model.TransactionTypeReversal {
		return s.createReversalTransaction(merchant.ID, t)
	}

	return model.Transaction{}, fmt.Errorf("invalid transaction type")
}

func (s *sqLiteDb) GetTransactions(query model.TransactionQuery) ([]model.Transaction, error) {

	rows, err := s.db.Joins("Merchant").Model(&Transaction{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []model.Transaction{}

	var t Transaction
	for rows.Next() {
		err = s.db.ScanRows(rows, &t)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, model.Transaction{
			Id:            t.TransactionId,
			MerchantId:    t.Merchant.MerchantId,
			Type:          t.Type,
			Amount:        t.Amount,
			Status:        t.Status,
			CustomerEmail: t.CustomerEmail,
			CustomerPhone: t.CustomerPhone,
		})
	}

	return transactions, nil
}

func (s *sqLiteDb) createAuthorizeTransaction(merchantId uint, t model.Transaction) (model.Transaction, error) {
	transaction := Transaction{
		MerchantID: merchantId,

		Type:          t.Type,
		Amount:        t.Amount,
		Status:        t.Status,
		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
	}

	err := s.db.Create(&transaction).Error
	if err != nil {
		return model.Transaction{}, err
	}

	t.Id = transaction.TransactionId

	return t, nil
}

func (s *sqLiteDb) createChargeTransaction(merchantId uint, t model.Transaction) (model.Transaction, error) {
	transaction := Transaction{
		MerchantID: merchantId,

		ParentId:      t.ParentId,
		Type:          t.Type,
		Amount:        t.Amount,
		Status:        t.Status,
		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
	}

	err := s.db.Create(&transaction).Error
	if err != nil {
		return model.Transaction{}, err
	}

	t.Id = transaction.TransactionId

	return t, nil
}

func (s *sqLiteDb) createRefundTransaction(merchantId uint, t model.Transaction) (model.Transaction, error) {
	transaction := Transaction{
		MerchantID: merchantId,

		ParentId:      t.ParentId,
		Type:          t.Type,
		Amount:        t.Amount,
		Status:        t.Status,
		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
	}

	txFunc := func(tx *gorm.DB) error {

		if err := tx.Model(&Transaction{}).Where("TransactionId = ?", t.ParentId).Update("Status", model.TransactionStatusRefunded).Error; err != nil {
			return err
		}

		err := tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		t.Id = transaction.TransactionId

		return nil
	}

	return t, s.db.Transaction(txFunc)
}

func (s *sqLiteDb) createReversalTransaction(merchantId uint, t model.Transaction) (model.Transaction, error) {
	transaction := Transaction{
		MerchantID: merchantId,

		ParentId:      t.ParentId,
		Type:          t.Type,
		Amount:        t.Amount,
		Status:        t.Status,
		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
	}

	txFunc := func(tx *gorm.DB) error {

		if err := tx.Model(&Transaction{}).Where("TransactionId = ?", t.ParentId).Update("Status", model.TransactionStatusReversed).Error; err != nil {
			return err
		}

		err := tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		t.Id = transaction.TransactionId

		return nil
	}

	return t, s.db.Transaction(txFunc)
}
