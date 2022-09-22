package server

type Admin struct {
	Id          string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Email       string `json:"email"`
}

type AdminResponse struct {
	Error  string  `json:"error"`
	Admins []Admin `json:"admins"`
}

// TODO: omit empty for Name, Description, Email, Staatus in request
type Merchant struct {
	Id                 string `json:"uuid"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Email              string `json:"email"`
	Status             string `json:"status"`
	TransactionsAmount int64  `json:"total_transaction_sum"`
}

type MerchantRequest struct {
	Merchant Merchant `json:"merchant"`
}

type MerchantResponse struct {
	Error     string     `json:"error"`
	Merchants []Merchant `json:"merchants"`
}

// TODO: omit empty for ParentId
type Transaction struct {
	Id            string `json:"uuid"`
	ParentId      string `json:"parent_uuid"`
	MerchantId    string `json:"merchant_uuid"`
	Type          string `json:"type"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type TransactionRequest struct {
	Transaction Transaction `json:"transaction"`
}

type TransactionResponse struct {
	Error        string        `json:"error"`
	Transactions []Transaction `json:"transactions"`
}
