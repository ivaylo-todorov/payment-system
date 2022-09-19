package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
}

func NewServer() *server {
	return &server{}
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// TODO: controller

	w.WriteHeader(http.StatusCreated)

	w.Write(body)
}

func (s *server) getMerchants(w http.ResponseWriter, r *http.Request) {
	log.Printf("got GET /merchants request\n")

	// TODO: controller

	response := &MerchantResponse{
		Merchants: []Merchant{
			{
				Id:   "1",
				Name: "Merchant One",
			},
			{
				Id:   "2",
				Name: "Merchant Two",
			},
		},
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// TODO: controller

	w.WriteHeader(http.StatusCreated)

	w.Write(body)
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

	// TODO: controller

	request.Merchant.Id = vars["id"]

	response := MerchantResponse{
		Merchants: []Merchant{
			request.Merchant,
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

	// TODO: controller

	response := MerchantResponse{
		Merchants: []Merchant{
			request.Merchant,
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

func (s *server) postTransaction(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got POST /transactions request\n")

	defer r.Body.Close()

	var request TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode request payload: %v", err), http.StatusBadRequest)
		return
	}

	// TODO: controller

	response := TransactionResponse{
		Transactions: []Transaction{
			request.Transaction,
		},
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response payload: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // http.StatusCreated
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResp)
}

func (s *server) getTransactions(w http.ResponseWriter, r *http.Request) {
	log.Printf("got GET /transactions request\n")

	// TODO: controller

	response := &TransactionResponse{
		Transactions: []Transaction{
			{
				Id:     "1",
				Status: "approved",
				Amount: 100,
			},
			{
				Id:     "2",
				Status: "refunded",
				Amount: 200,
			},
		},
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
