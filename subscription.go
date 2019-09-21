package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"log"
	"net/http"
	"os"
)

func (rs *RestServer) GetCustomerSessionV2(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s from %s", ServiceName, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]

	decoder := json.NewDecoder(req.Body)
	var t []string
	if err := decoder.Decode(&t); err != nil {
		StripeError(w, err.Error())
		return
	}

	if len(t) == 0 {
		StripeError(w, "no plans provided")
	}

	sub := &stripe.CheckoutSessionSubscriptionDataParams{
		Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{},
	}
	for _, plan := range t {
		sub.Items = append(sub.Items, &stripe.CheckoutSessionSubscriptionDataItemsParams{
			Plan: stripe.String(plan),
		})
	}

	customer, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
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

func (rs *RestServer) CreateSession(customer *stripe.Customer, sub *stripe.CheckoutSessionSubscriptionDataParams) (*stripe.CheckoutSession, error) {
	successUrl := os.Getenv("STRIPE_SUCCESS_URL")
	cancelUrl := os.Getenv("STRIPE_CANCEL_URL")

	stripeParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		SubscriptionData: sub,
		Customer:         &customer.ID,
		SuccessURL:       stripe.String(successUrl),
		CancelURL:        stripe.String(cancelUrl),
	}

	return session.New(stripeParams)
}
