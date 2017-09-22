package http

import (
	"net/http"
	"github.com/stripe/stripe-go/customer"
	"encoding/json"
	"github.com/gorilla/mux"
)

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
