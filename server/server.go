package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/ivaylo-todorov/payment-system/model"
	"github.com/ivaylo-todorov/payment-system/model/controller"
	"github.com/ivaylo-todorov/payment-system/store"
)

type server struct {
	Cancel     context.CancelFunc
	Controller controller.Controller
}

func NewServer(settings model.ApplicationSettings) (*server, error) {
	store, err := store.NewStore(settings.StoreSettings)
	if err != nil {
		return nil, err
	}

	c, err := controller.NewController(settings, store)
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

func (s *server) StartTransactionsCleanup(interval time.Duration) error {
	ctx, cancel := context.WithCancel(context.Background())

	s.Cancel = cancel

	go func() {
		for {
			select {
			case <-time.After(interval * time.Minute):
				now := time.Now()
				err := s.Controller.DeleteTransactions(model.TransactionQuery{
					OlderThan: &now,
				})
				if err != nil {
					log.Printf("Cleaning up transactions failed, %s", err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (s *server) StopTransactionsCleanup() error {
	s.Cancel()
	return nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODOO: Authentication

		fn(w, r)
	}
}
