package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/source"
	"log"
	"net/http"
)

type CardRequest struct {
	Token string `json:"token"`
}

func (rs *RestServer) UpdateCustomerCard(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s from %s", ServiceName, req.URL.Path, req.RemoteAddr)
	params := mux.Vars(req)
	id := params["id"]

	decoder := json.NewDecoder(req.Body)
	cardRequest := &CardRequest{}
	err := decoder.Decode(cardRequest)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	if "" == cardRequest.Token {
		BadRequest(w, "token is required in this request")
		return
	}

	c, err := getCustomer(id)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	crd, err := card.New(&stripe.CardParams{
		Customer: &c.ID,
		Token:    &cardRequest.Token,
	})

	if err != nil {
		StripeError(w, err.Error())
		return
	}

	customerParams := &stripe.CustomerParams{
		DefaultSource: &crd.ID,
	}
	customer, err := customer.Update(
		id,
		customerParams,
	)
	if err != nil {
		StripeError(w, err.Error())
		return
	}

	js, _ := json.Marshal(customer)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func DeleteAllCards(id string) error {
	c, err := getCustomer(id)
	if err != nil {
		return err
	}
	for _, card := range c.Sources.Data {
		params := &stripe.SourceObjectDetachParams{
			Customer: stripe.String(id),
		}
		_, err = source.Detach(card.ID, params)
		if err != nil {
			return err
		}
	}
	return nil
}
