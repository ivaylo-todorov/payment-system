package db

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ivaylo-todorov/payment-system/model"
)

var db *sqLiteDb

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func setup() {
	_, err := os.Stat("payment_system.db")
	if err == nil {
		err := os.Remove("payment_system.db")
		if err != nil {
			log.Fatalf("Cannot remove old db file")
		}
	}

	db, err = NewDb(model.StoreSettings{})
	if err != nil {
		log.Fatalf("Cannot create new db: %s", err.Error())
	}
}

func shutdown() {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestCreateAdmin(t *testing.T) {
	expected := model.Admin{
		Name:        "name",
		Description: "description",
		Email:       RandomString(8),
	}

	a, err := db.CreateAdmin(expected)
	require.NoError(t, err)

	assert.NotZero(t, a.Id)

	actual := Admin{}
	err = db.Db().Raw("select * from admins where admin_id = ?", a.Id).
		Scan(&actual).Error

	require.NoError(t, err)

	user := User{}
	err = db.Db().Raw("select * from users where id = ?", actual.UserID).
		Scan(&user).Error

	require.NoError(t, err)

	assert.Equal(t, expected.Name, user.Name)
	assert.Equal(t, expected.Description, user.Description)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, model.UserRoleAdmin, user.Role)
}

func TestCreateMerchant(t *testing.T) {
	expected := model.Merchant{
		Name:        "name",
		Description: "description",
		Email:       RandomString(8),
		Status:      "status",
	}

	m, err := db.CreateMerchant(expected)
	require.NoError(t, err)

	assert.NotZero(t, m.Id)

	actual := Merchant{}
	err = db.Db().Raw("select * from merchants where merchant_id = ?", m.Id).
		Scan(&actual).Error

	require.NoError(t, err)

	assert.Equal(t, expected.Status, actual.Status)

	user := User{}
	err = db.Db().Raw("select * from users where id = ?", actual.UserID).
		Scan(&user).Error

	require.NoError(t, err)

	assert.Equal(t, expected.Name, user.Name)
	assert.Equal(t, expected.Description, user.Description)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, model.UserRoleMerchant, user.Role)
}

func TestUpdateMerchant(t *testing.T) {
	expected := model.Merchant{
		Name:        "name",
		Description: "description",
		Email:       RandomString(8),
		Status:      "status",
	}

	m, err := db.CreateMerchant(expected)
	require.NoError(t, err)

	assert.NotZero(t, m.Id)

	expected.Id = m.Id
	expected.Name = "new name"
	expected.Description = "new description"
	expected.Email = RandomString(8)
	expected.Status = "new status"

	m, err = db.UpdateMerchant(expected)
	require.NoError(t, err)

	actual := Merchant{}
	err = db.Db().Raw("select * from merchants where merchant_id = ?", m.Id).
		Scan(&actual).Error

	require.NoError(t, err)

	assert.Equal(t, expected.Status, actual.Status)

	user := User{}
	err = db.Db().Raw("select * from users where id = ?", actual.UserID).
		Scan(&user).Error

	require.NoError(t, err)

	assert.Equal(t, expected.Name, user.Name)
	assert.Equal(t, expected.Description, user.Description)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, model.UserRoleMerchant, user.Role)
}

func TestDeleteMerchant(t *testing.T) {
	expected := model.Merchant{
		Name:        "name",
		Description: "description",
		Email:       RandomString(8),
		Status:      "status",
	}

	m, err := db.CreateMerchant(expected)
	require.NoError(t, err)

	assert.NotZero(t, m.Id)

	err = db.DeleteMerchant(m.Id)
	require.NoError(t, err)

	actual := Merchant{}
	err = db.Db().Raw("select * from merchants where merchant_id = ? and deleted_at is null", m.Id).
		Scan(&actual).Error

	require.NoError(t, err)

	assert.NotEqual(t, nil, actual.DeletedAt)

	user := User{}
	err = db.Db().Raw("select * from users where id = ?", actual.UserID).
		Scan(&user).Error

	require.NoError(t, err)

	assert.NotEqual(t, nil, user.DeletedAt)
}

// TODO: add tests for transactions

func RandomString(n int) string {
	rand.Seed(time.Now().UnixMicro())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
