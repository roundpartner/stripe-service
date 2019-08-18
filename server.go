package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roundpartner/go/ha"
	"github.com/roundpartner/go/transaction"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"log"
	"net/http"
	"time"
)

func ListenAndServe(port int) {
	address := fmt.Sprintf(":%d", port)
	rs := New()
	server := &http.Server{
		Addr:         address,
		Handler:      rs.router(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	ha.GracefulShutdown(server, ServiceName)

	log.Printf("[INFO] [%s] Server starting on port %d", ServiceName, port)
	err := server.ListenAndServe()
	if nil != err {
		log.Printf("[INFO] [%s] %s", ServiceName, err.Error())
	}
}

type RestServer struct {
	SNSService *SNSService
}

func (rs *RestServer) router() *mux.Router {
	router := mux.NewRouter()
	ha.Check(router)
	router.HandleFunc("/charge", rs.Charge).Methods("POST")
	router.HandleFunc("/customer", rs.Customers).Methods("GET")
	router.HandleFunc("/customer/{id}", rs.GetCustomer).Methods("GET")
	router.HandleFunc("/customer/{id}", rs.UpdateCustomer).Methods("PUT")
	router.HandleFunc("/customer/{id}/session/{plan}", rs.GetCustomerSession).Methods("POST")
	router.HandleFunc("/customer/{id}/subscription", rs.GetCustomerSubscriptions).Methods("GET")
	router.HandleFunc("/customer", rs.NewCustomer).Methods("POST")
	router.HandleFunc("/customer/{id}/card", rs.UpdateCustomerCard).Methods("PUT")
	router.HandleFunc("/reload", rs.ReloadCustomers).Methods("GET")
	router.HandleFunc("/webhook", rs.WebHook).Methods("POST")
	return router
}

func New() *RestServer {
	return &RestServer{}
}

type ChargeRequest struct {
	Trans    string `json:"trans_id"`
	Token    string `json:"token"`
	Amount   int64  `json:"amount"`
	Desc     string `json:"desc"`
	Email    string `json:"receipt_email"`
	Business string `json:"business_name"`
	Customer string `json:"customer"`
	Callback string `json:"callback"`
	Currency string
}

func (rs *RestServer) Charge(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	t := &ChargeRequest{
		Currency: "gbp",
	}
	err := decoder.Decode(t)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	if t.Amount < 30 {
		BadRequest(w, "Amount must be at least 30 pence")
		return
	}

	token := t.Token
	params := &stripe.ChargeParams{
		Amount:      &t.Amount,
		Currency:    &t.Currency,
		Description: &t.Desc,
	}
	if "" != t.Customer {
		params.Customer = &t.Customer
	}
	if "" != t.Email {
		params.ReceiptEmail = &t.Email
	}
	params.AddMetadata("trans_id", t.Trans)
	params.AddMetadata("business_name", t.Business)
	if "" != token {
		params.SetSource(token)
	}

	charge, err := charge.New(params)

	if err != nil {
		StripeError(w, err.Error())
		if t.Callback != "" {
			transaction.CallbackTransactionFailed(t.Callback, t.Trans, err.Error())
		}
		return
	}

	if t.Callback != "" {
		transaction.CallbackTransactionSuccessful(t.Callback, t.Trans, "5", charge.FailureMessage, charge.Amount)
	}

	js, _ := json.Marshal(charge)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
