package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"log"
)

func ListenAndServe() {
	rs := New()
	server := &http.Server{Addr: ":57493", Handler: rs.router()}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		signal.Stop(c)
		log.Println("Server shutting down gracefully")
		server.Shutdown(nil)
	}()

	log.Println("Server starting")
	err := server.ListenAndServe()
	if nil != err {
		log.Println(err.Error())
	}
}

type RestServer struct {
}

func (rs *RestServer) router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/charge", rs.Charge).Methods("POST")
	router.HandleFunc("/customer", rs.Customers).Methods("GET")
	router.HandleFunc("/customer/{id}", rs.GetCustomer).Methods("GET")
	router.HandleFunc("/customer", rs.NewCustomer).Methods("POST")
	router.HandleFunc("/customer/{id}/card", rs.UpdateCustomerCard).Methods("PUT")
	return router
}

func New() *RestServer {
	return &RestServer{}
}

type ChargeRequest struct {
	Trans    string `json:"trans_id"`
	Token    string `json:"token"`
	Amount   uint64 `json:"amount"`
	Desc     string `json:"desc"`
	Email    string `json:"receipt_email"`
	Business string `json:"business_name"`
	Customer string `json:"customer"`
}

func (rs *RestServer) Charge(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	t := &ChargeRequest{}
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
		Amount:   t.Amount,
		Currency: "gbp",
		Desc:     t.Desc,
		Email:    t.Email,
		Customer: t.Customer,
	}
	params.AddMeta("trans_id", t.Trans)
	params.AddMeta("business_name", t.Business)
	params.SetSource(token)

	charge, err := charge.New(params)

	if err != nil {
		StripeError(w, err.Error())
		return
	}

	js, _ := json.Marshal(charge)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
