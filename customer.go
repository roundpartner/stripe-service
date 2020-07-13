package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"log"
	"net/http"
	"sync"
)

type CustomerRequest struct {
	Account  string `json:"account"`
	User     string `json:"user"`
	Email    string `json:"email"`
	Desc     string `json:"desc"`
	Token    string `json:"token"`
	Discount string `json:"discount"`
}

type CustomersRequest struct {
	Limit string `json:"limit"`
	After string `json:"after"`
}

func (rs *RestServer) Customers(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	t := &CustomersRequest{Limit: "100"}
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
		list.Data = append(list.Data, i.Customer())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(list.Data)
	w.Write(js)
}

func (rs *RestServer) GetCustomer(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]
	customer, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}
	js, _ := json.Marshal(customer)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getCustomer(id string) (*stripe.Customer, error) {
	return customer.Get(id, nil)
}

func NewCustomerParam(t *CustomerRequest) *stripe.CustomerParams {
	customerParams := &stripe.CustomerParams{
		Description: &t.Desc,
		Email:       &t.Email,
	}
	customerParams.AddMetadata("account", t.Account)
	customerParams.AddMetadata("user", t.User)
	customerParams.AddMetadata("discount", t.Discount)
	return customerParams
}

func (rs *RestServer) NewCustomer(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	decoder := json.NewDecoder(req.Body)
	t := &CustomerRequest{}
	err := decoder.Decode(t)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	customerParams := NewCustomerParam(t)
	newCustomer, err := customer.New(customerParams)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	if "" == t.Token {
		js, _ := json.Marshal(newCustomer)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return
	}

	log.Printf("[INFO] [%s] Adding token to customer %s", ServiceName, newCustomer.ID)

	card, err := card.New(&stripe.CardParams{
		Customer: &newCustomer.ID,
		Token:    &t.Token,
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

func delete(id string) bool {
	_, err := customer.Del(id, nil)
	if nil != err {
		return false
	}
	return true
}

func (rs *RestServer) UpdateCustomer(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]

	decoder := json.NewDecoder(req.Body)
	t := &CustomerRequest{}
	err := decoder.Decode(t)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	customerParams := NewCustomerParam(t)
	updatedCustomer, err := customer.Update(id, customerParams)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	js, _ := json.Marshal(updatedCustomer)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (rs *RestServer) UpdateDiscount(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]
	coupon := params["coupon"]

	customerParams := &stripe.CustomerParams{}
	customerParams.Coupon = stripe.String(coupon)

	cust, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	if cust.Discount != nil {
		log.Printf("[INFO] [%s] Customer %s already has a discount so %s will not be applied", ServiceName, id, coupon)
		js, _ := json.Marshal(cust)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return
	}

	updatedCustomer, err := customer.Update(id, customerParams)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	js, _ := json.Marshal(updatedCustomer)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

type CustomerMeta struct {
	Account string `json:"account"`
	Token   string `json:"token"`
}

var customerMetaList map[string]*CustomerMeta

func (rs *RestServer) ReloadCustomers(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	loadCustomers()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

func loadCustomers() {
	customerMetaList = make(map[string]*CustomerMeta)
	params := &stripe.CustomerListParams{}
	list := customer.List(params)
	for list.Next() {
		if "" == list.Customer().Metadata["account"] {
			log.Printf("Customer %s does not have account set", list.Customer().ID)
			continue
		}
		if get(list.Customer().Metadata["account"]) != nil {
			log.Printf("Customer %s is a duplicate", list.Customer().ID)
			continue
		}
		if list.Customer().DefaultSource == nil {
			log.Printf("Customer %s has not default source", list.Customer().ID)
			continue
		}
		cm := &CustomerMeta{
			list.Customer().Metadata["account"],
			list.Customer().DefaultSource.ID,
		}
		add(cm)
	}
}

var customerMutex = struct {
	sync.RWMutex
	customers map[string]*CustomerMeta
}{customers: make(map[string]*CustomerMeta)}

func get(id string) *CustomerMeta {
	customerMutex.RLock()
	c, ok := customerMetaList[id]
	customerMutex.RUnlock()
	if ok == false {
		return nil
	}
	return c
}

func add(c *CustomerMeta) {
	customerMutex.Lock()
	customerMetaList[c.Account] = c
	customerMutex.Unlock()
}

func (rs *RestServer) GetCustomerSubscriptions(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]
	customer, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	if customer.Subscriptions == nil {
		customer.Subscriptions = &stripe.SubscriptionList{}
	}

	js, _ := json.Marshal(customer.Subscriptions.Data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (rs *RestServer) GetCustomerSession(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]
	plan := params["plan"]
	customer, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	sub := &stripe.CheckoutSessionSubscriptionDataParams{
		Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
			&stripe.CheckoutSessionSubscriptionDataItemsParams{
				Plan: stripe.String(plan),
			},
		},
	}

	session, err := rs.CreateSession(customer, sub)

	if err != nil {
		StripeError(w, err.Error())
		return
	}

	js, _ := json.Marshal(session)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
