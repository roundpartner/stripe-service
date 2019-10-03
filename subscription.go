package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go/subitem"
	"log"
	"net/http"
	"os"
)

func (rs *RestServer) GetCustomerSessionV2(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s from %s", ServiceName, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]

	plans, err := rs.DecodePlans(req)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	sub := &stripe.CheckoutSessionSubscriptionDataParams{
		Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{},
	}
	for _, plan := range plans {
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

func (rs *RestServer) DecodePlans(req *http.Request) ([]string, error) {
	decoder := json.NewDecoder(req.Body)
	var plans []string
	if err := decoder.Decode(&plans); err != nil {
		return plans, err
	}

	if len(plans) == 0 {
		return plans, errors.New("no plans provided")
	}

	return plans, nil
}

func (rs *RestServer) UpgradeSubscription(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s from %s", ServiceName, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]

	plans, err := rs.DecodePlans(req)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	customer, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	if customer.Subscriptions.TotalCount < 1 {
		StripeError(w, "no subscriptions found")
		return
	}

	if customer.Subscriptions.TotalCount > 1 {
		log.Printf("[INFO] [%s] %s has %d subscription found but expected one", ServiceName, customer.ID, customer.Subscriptions.TotalCount)
	}

	subscription, err := sub.Get(customer.Subscriptions.Data[0].ID, nil)

	subParams := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
		Items:             []*stripe.SubscriptionItemsParams{},
	}

	for _, plan := range plans {
		subParams.Items = append(subParams.Items, &stripe.SubscriptionItemsParams{
			Plan: stripe.String(plan),
		})
	}

	subscription, err = sub.Update(subscription.ID, subParams)

	if err != nil {
		StripeError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

func (rs *RestServer) DowngradeSubscription(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s from %s", ServiceName, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]

	plans, err := rs.DecodePlans(req)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	customer, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	if customer.Subscriptions.TotalCount < 1 {
		StripeError(w, "no subscriptions found")
		return
	}

	err = rs.RemovePlans(customer.ID, plans)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

func (rs *RestServer) RemovePlans(customer string, plans []string) error {
	for _, plan := range plans {
		subParams := &stripe.SubscriptionListParams{
			Customer: customer,
		}
		subItems := sub.List(subParams)
		for subItems.Next() {
			subItem := subItems.Current().(*stripe.Subscription)

			subItemParams := &stripe.SubscriptionItemListParams{
				Subscription: stripe.String(subItem.ID),
			}
			planItems := subitem.List(subItemParams)
			for planItems.Next() {
				planItem := planItems.Current().(*stripe.SubscriptionItem)
				if planItem.Plan.ID == plan {
					_, err := subitem.Del(planItem.ID, nil)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (rs *RestServer) CancelSubscription(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s from %s", ServiceName, req.URL.Path, req.RemoteAddr)

	params := mux.Vars(req)
	id := params["id"]

	customer, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	if customer.Subscriptions.TotalCount < 1 {
		StripeError(w, "no subscriptions found")
		return
	}

	if customer.Subscriptions.TotalCount > 1 {
		log.Printf("[INFO] [%s] %s has %d subscription found but expected one", ServiceName, customer.ID, customer.Subscriptions.TotalCount)
	}

	subscription, err := sub.Get(customer.Subscriptions.Data[0].ID, nil)

	subParams := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err = sub.Update(subscription.ID, subParams)

	if err != nil {
		StripeError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}
