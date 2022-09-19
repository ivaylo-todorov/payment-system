package server

type Merchant struct {
	Id                 string `json:"uuid"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Email              string `json:"email"`
	Status             string `json:"status"`
	TransactionsAmount uint64 `json:"total_transaction_sum"`
}

type MerchantRequest struct {
	Merchant Merchant `json:"merchant"`
}

type MerchantResponse struct {
	Merchants []Merchant `json:"merchants"`
}

type Transaction struct {
	Id            string `json:"uuid"`
	Amount        uint64 `json:"amount"`
	Status        string `json:"status"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type TransactionRequest struct {
	Transaction Transaction `json:"transaction"`
}

type TransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
}
