package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/model/controller"
)

type server struct {
	Controller controller.Controller
}

func NewServer(settings model.ApplicationSettings) (*server, error) {
	c, err := controller.NewController(settings)
	if err != nil {
		return nil, err
	}

	return &server{
		Controller: c,
	}, nil
}

func (s *server) Start() error {

	r := mux.NewRouter()

	r.HandleFunc("/", makeHandler(s.root)).Methods("GET")
	r.HandleFunc("/admins", makeHandler(s.createAdmins)).Methods("POST")
	r.HandleFunc("/merchants", makeHandler(s.getMerchants)).Methods("GET")
	r.HandleFunc("/merchants", makeHandler(s.createMerchants)).Methods("POST")
	r.HandleFunc("/merchants/{id}", makeHandler(s.updateMerchant)).Methods("POST")
	r.HandleFunc("/merchants", makeHandler(s.deleteMerchants)).Methods("DELETE")
	r.HandleFunc("/transactions", makeHandler(s.getTransactions)).Methods("GET")
	r.HandleFunc("/transactions", makeHandler(s.postTransaction)).Methods("POST")

	err := http.ListenAndServe(":8080", r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
		return nil
	}

	fmt.Printf("error starting server: %s\n", err)

	return err
}

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	log.Printf("got GET / request\n")
	w.Write([]byte("A Payment System!"))
}

func (s *server) createAdmins(w http.ResponseWriter, r *http.Request) {
	log.Printf("got POST /admins request\n")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	admins, err := ConvertCsvToAdmins(body)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	response := &AdminResponse{}

	res, err := s.Controller.CreateAdmins(admins)
	if err != nil {
		response.Error = err.Error()
	}

	for _, a := range res {
		response.Admins = append(response.Admins, ConvertAdminFromModel(a))
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response payload: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResp)
}

func (s *server) getMerchants(w http.ResponseWriter, r *http.Request) {
	log.Printf("got GET /merchants request\n")

	query := model.MerchantQuery{}

	merchants, err := s.Controller.GetMerchants(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get merchants: %v", err), http.StatusInternalServerError)
		return
	}

	response := &MerchantResponse{}

	for _, m := range merchants {
		response.Merchants = append(response.Merchants, ConvertMerchantFromModel(m))
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response payload: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResp)
}

func (s *server) createMerchants(w http.ResponseWriter, r *http.Request) {
	log.Printf("got POST /merchants request\n")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	merchants, err := ConvertCsvToMerchants(body)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	response := &MerchantResponse{}

	res, err := s.Controller.CreateMerchants(merchants)
	if err != nil {
		response.Error = err.Error()
	}

	for _, m := range res {
		response.Merchants = append(response.Merchants, ConvertMerchantFromModel(m))
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response payload: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResp)
}

func (s *server) updateMerchant(w http.ResponseWriter, r *http.Request) {
	log.Printf("got POST /merchants/{id} request\n")

	defer r.Body.Close()

	var request MerchantRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode request payload: %v", err), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	request.Merchant.Id = vars["id"]

	merchant, err := ConvertMerchantToModel(request.Merchant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	merchant, err = s.Controller.UpdateMerchant(merchant)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not delete merchant: %v", err), http.StatusInternalServerError)
		return
	}

	response := MerchantResponse{
		Merchants: []Merchant{
			ConvertMerchantFromModel(merchant),
		},
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response payload: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResp)
}

func (s *server) deleteMerchants(w http.ResponseWriter, r *http.Request) {
	log.Printf("got DELETE /merchants request\n")

	defer r.Body.Close()

	var request MerchantRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode request payload: %v", err), http.StatusBadRequest)
		return
	}

	merchant, err := ConvertMerchantToModel(request.Merchant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Controller.DeleteMerchant(merchant)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not delete merchant: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) postTransaction(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got POST /transactions request\n")

	defer r.Body.Close()

	var request TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode request payload: %v", err), http.StatusBadRequest)
		return
	}

	transaction, err := ConvertTransactionToModel(request.Transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err = s.Controller.StartTransaction(transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not start transaction: %v", err), http.StatusInternalServerError)
		return
	}

	response := TransactionResponse{
		Transactions: []Transaction{
			ConvertTransactionFromModel(transaction),
		},
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response payload: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResp)
}

func (s *server) getTransactions(w http.ResponseWriter, r *http.Request) {
	log.Printf("got GET /transactions request\n")

	query := model.TransactionQuery{}

	transactions, err := s.Controller.GetTransactions(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	response := &TransactionResponse{}

	for _, t := range transactions {
		response.Transactions = append(response.Transactions, ConvertTransactionFromModel(t))
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response payload: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResp)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODOO: Authentication

		fn(w, r)
	}
}
