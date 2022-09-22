package controller_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/model/controller"
	"github.com/ivaylo-todorov/payment-system/store"
)

func TransactionStatus(status string) types.GomegaMatcher {
	return &representJSONMatcher{
		expected: status,
	}
}

type representJSONMatcher struct {
	expected, actual string
}

func (m *representJSONMatcher) Match(actual interface{}) (success bool, err error) {
	t, ok := actual.(model.Transaction)
	if !ok {
		return false, fmt.Errorf("expected input of type model.Transaction")
	}
	m.actual = t.Status
	return t.Status == m.expected, nil
}

func (m *representJSONMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected transaction status: %s, received status: %s", m.expected, m.actual)
}

func (m *representJSONMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected transaction status: %s, maches received", m.expected)
}

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model Controller Suite")
}

var _ = Describe("Model Controller", func() {
	s, err := store.NewMockStore()
	Expect(err).To(BeNil())

	c, err := controller.NewController(model.ApplicationSettings{}, s)
	Expect(err).To(BeNil())

	Context("initially", func() {
		It("has 0 merchants", func() {
			Expect(c.GetMerchants(model.MerchantQuery{})).To(HaveLen(0))
		})
		It("has 0 transactions", func() {
			Expect(c.GetTransactions(model.TransactionQuery{})).To(HaveLen(0))
		})
	})

	Context("when a new admin is added", func() {
		admin := model.Admin{
			Id:          store.AdminUuid,
			Name:        "admin_one",
			Description: "some admin",
			Email:       "admin_one@email.com",
		}

		It("with empty name", func() {
			a := admin
			a.Name = ""
			_, err := c.CreateAdmins([]model.Admin{a})
			Expect(err).Should(HaveOccurred())
		})

		It("with invalid email", func() {
			a := admin
			a.Email = "my_email@"
			_, err := c.CreateAdmins([]model.Admin{a})
			Expect(err).Should(HaveOccurred())
		})

		It("successfully", func() {
			admins := []model.Admin{admin}
			Expect(c.CreateAdmins(admins)).Should(Equal(admins))
		})
	})

	Context("when a new merchant is added", func() {
		merchant := model.Merchant{
			Id:     store.MerchantOneUuid,
			Name:   "merchant_one",
			Email:  "merchant_one@email.com",
			Status: model.MerchantStatusActive,
		}

		It("with empty name", func() {
			m := merchant
			m.Name = ""
			_, err := c.CreateMerchants([]model.Merchant{m})
			Expect(err).Should(HaveOccurred())
		})

		It("with invalid email", func() {
			m := merchant
			m.Email = "my_email@"
			_, err := c.CreateMerchants([]model.Merchant{m})
			Expect(err).Should(HaveOccurred())
		})

		It("successfully", func() {
			merchants := []model.Merchant{merchant}
			Expect(c.CreateMerchants(merchants)).Should(Equal(merchants))
		})

		It("has 1 merchants", func() {
			Expect(c.GetMerchants(model.MerchantQuery{})).To(HaveLen(1))
		})
	})

	Context("when a second merchant is added", func() {
		merchant := model.Merchant{
			Id:     store.MerchantTwoUuid,
			Name:   "merchant_two",
			Email:  "merchant_twoe@email.com",
			Status: model.MerchantStatusInactive,
		}

		It("with invalid status", func() {
			m := merchant
			m.Status = "wrong status"
			_, err := c.CreateMerchants([]model.Merchant{m})
			Expect(err).Should(HaveOccurred())
		})

		It("successfully", func() {
			merchants := []model.Merchant{merchant}
			Expect(c.CreateMerchants(merchants)).Should(Equal(merchants))
		})

		It("has w merchants", func() {
			Expect(c.GetMerchants(model.MerchantQuery{})).To(HaveLen(2))
		})
	})

	Context("when the first merchant is updated", func() {

		It("and merchant id is empty", func() {
			m := model.Merchant{
				Name: "merchant_one_updated",
			}
			_, err := c.UpdateMerchant(m)
			Expect(err).Should(HaveOccurred())
		})

		// TODO: check Name, Description, Email update
	})

	Context("when an authorize transaction is created", func() {
		transaction := model.Transaction{
			Id:            store.AuthorizeTransactionOneUuid,
			MerchantId:    store.MerchantOneUuid,
			Type:          model.TransactionTypeAuthorize,
			Amount:        100,
			CustomerEmail: "customer_one@email.com",
		}

		It("with invalid type", func() {
			t := transaction
			t.Type = "type"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with empty merchant id", func() {
			t := transaction
			t.Id = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with invalid customer email", func() {
			t := transaction
			t.CustomerEmail = "customer"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with zero amount", func() {
			t := transaction
			t.Amount = 0
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with negative amount", func() {
			t := transaction
			t.Amount = -1
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("successfully", func() {
			Expect(c.StartTransaction(transaction)).Should(TransactionStatus(model.TransactionStatusApproved))
		})

		It("has 1 transactions", func() {
			Expect(c.GetTransactions(model.TransactionQuery{})).To(HaveLen(1))
		})
	})

	Context("when a charge transaction is created", func() {
		transaction := model.Transaction{
			Id:            store.ChargeTransactionUuid,
			ParentId:      store.AuthorizeTransactionOneUuid,
			MerchantId:    store.MerchantOneUuid,
			Type:          model.TransactionTypeCharge,
			Amount:        100,
			CustomerEmail: "customer_one@email.com",
		}

		It("with invalid type", func() {
			t := transaction
			t.Type = "type"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with empty merchant id", func() {
			t := transaction
			t.Id = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with invalid customer email", func() {
			t := transaction
			t.CustomerEmail = "customer"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with zero amount", func() {
			t := transaction
			t.Amount = 0
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with negative amount", func() {
			t := transaction
			t.Amount = -1
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with empty parent id", func() {
			t := transaction
			t.ParentId = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with bigger amount", func() {
			t := transaction
			t.Amount = 1111
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("successfully", func() {
			Expect(c.StartTransaction(transaction)).Should(TransactionStatus(model.TransactionStatusApproved))
		})

		It("has 2 transactions", func() {
			Expect(c.GetTransactions(model.TransactionQuery{})).To(HaveLen(2))
		})
	})

	Context("when a charge transaction is created", func() {
		transaction := model.Transaction{
			Id:            store.RefundTransactionUuid,
			ParentId:      store.ChargeTransactionUuid,
			MerchantId:    store.MerchantOneUuid,
			Type:          model.TransactionTypeRefund,
			Amount:        100,
			CustomerEmail: "customer_one@email.com",
		}

		It("with invalid type", func() {
			t := transaction
			t.Type = "type"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with empty merchant id", func() {
			t := transaction
			t.Id = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with invalid customer email", func() {
			t := transaction
			t.CustomerEmail = "customer"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with zero amount", func() {
			t := transaction
			t.Amount = 0
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with negative amount", func() {
			t := transaction
			t.Amount = -1
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with empty parent id", func() {
			t := transaction
			t.ParentId = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with wrong parent type", func() {
			t := transaction
			t.ParentId = store.AuthorizeTransactionOneUuid
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with different amount", func() {
			t := transaction
			t.Amount = 1111
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("successfully", func() {
			Expect(c.StartTransaction(transaction)).Should(TransactionStatus(model.TransactionStatusRefunded))
		})

		It("has 3 transactions", func() {
			Expect(c.GetTransactions(model.TransactionQuery{})).To(HaveLen(3))
		})
	})

	Context("when a second authorize transaction is created", func() {
		transaction := model.Transaction{
			Id:            store.AuthorizeTransactionTwoUuid,
			MerchantId:    store.MerchantTwoUuid,
			Type:          model.TransactionTypeAuthorize,
			Amount:        200,
			CustomerEmail: "customer_one@email.com",
		}

		It("with inactive merchant", func() {
			t := transaction
			t.Id = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		Context("then the merchant is activeted", func() {
			m := model.Merchant{
				Id:     store.MerchantTwoUuid,
				Status: model.MerchantStatusActive,
			}
			r, err := c.UpdateMerchant(m)
			Expect(err).Should(Succeed())
			Expect(r.Status).To(Equal(model.MerchantStatusActive))

			It("successfully", func() {
				Expect(c.StartTransaction(transaction)).Should(TransactionStatus(model.TransactionStatusApproved))
			})

			It("has 4 transactions", func() {
				Expect(c.GetTransactions(model.TransactionQuery{})).To(HaveLen(4))
			})
		})
	})

	Context("when a reversal transaction is created", func() {
		transaction := model.Transaction{
			Id:            store.ReversalTransactionUuid,
			ParentId:      store.AuthorizeTransactionTwoUuid,
			MerchantId:    store.MerchantTwoUuid,
			Type:          model.TransactionTypeReversal,
			Amount:        0,
			CustomerEmail: "customer_one@email.com",
		}

		It("with invalid type", func() {
			t := transaction
			t.Type = "type"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with empty merchant id", func() {
			t := transaction
			t.Id = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with invalid customer email", func() {
			t := transaction
			t.CustomerEmail = "customer"
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with non zero amount", func() {
			t := transaction
			t.Amount = 120
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("with empty parent id", func() {
			t := transaction
			t.ParentId = uuid.Nil
			_, err := c.StartTransaction(t)
			Expect(err).Should(HaveOccurred())
		})

		It("successfully", func() {
			Expect(c.StartTransaction(transaction)).Should(TransactionStatus(model.TransactionStatusReversed))
		})

		It("has 5 transactions", func() {
			Expect(c.GetTransactions(model.TransactionQuery{})).To(HaveLen(5))
		})
	})

	Context("when the first merchant is deleted", func() {

		It("and merchant id is empty", func() {
			m := model.Merchant{
				Name: "merchant_one_updated",
			}
			Expect(c.DeleteMerchant(m)).Should(HaveOccurred())
		})

		It("successfully", func() {
			m := model.Merchant{
				Id: store.MerchantOneUuid,
			}
			Expect(c.DeleteMerchant(m)).Should(Succeed())
		})
	})

})
