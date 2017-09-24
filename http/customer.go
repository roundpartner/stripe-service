package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"net/http"
	"github.com/stripe/stripe-go/card"
)

type CustomerRequest struct {
	Account string `json:"account"`
	Email   string `json:"email"`
	Desc    string `json:"desc"`
	Token   string `json:"token"`
}

type CustomersRequest struct {
	Limit string `json:"limit"`
	After string `json:"after"`
}

func (rs *RestServer) Customers(w http.ResponseWriter, req *http.Request) {
	t := &CustomersRequest{Limit: "10"}
	if req.ContentLength > 0 {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(t)
		if err != nil {
			BadRequest(w, err.Error())
			return
		}
	}

	params := &stripe.CustomerListParams{}
	if "" != t.After {
		params.Filters.AddFilter("starting_after", "", t.After)
	}
	if "" != t.Limit {
		params.Filters.AddFilter("limit", "", t.Limit)
	}
	i := customer.List(params)
	list := stripe.CustomerList{}
	for i.Next() {
		list.Values = append(list.Values, i.Customer())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(list.Values)
	w.Write(js)
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

	card, err := card.New(&stripe.CardParams{
		Customer: customer.ID,
		Token: t.Token,
	})

	if err != nil {
		StripeError(w, err.Error())
		return
	}

	js, _ := json.Marshal(card.Customer)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
