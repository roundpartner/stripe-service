package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"net/http"
)

type CustomerRequest struct {
	Account string `json:"account_id"`
	Email   string `json:"email"`
	Desc    string `json:"desc"`
}

func (rs *RestServer) GetCustomer(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	customer, err := customer.Get(id, nil)
	if err != nil {
		StripeError(w, err.Error())
		return
	}
	js, _ := json.Marshal(customer)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (rs *RestServer) NewCustomer(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	t := &CustomerRequest{}
	err := decoder.Decode(t)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}
	customerParams := &stripe.CustomerParams{
		Desc:  t.Desc,
		Email: t.Email,
	}
	customerParams.AddMeta("account_id", t.Account)
	customer, err := customer.New(customerParams)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	js, _ := json.Marshal(customer)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
