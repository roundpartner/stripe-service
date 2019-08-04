package client

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"log"
	"net/http"
)

func Subscription(customer string) *stripe.SubscriptionList {
	client := &http.Client{}
	url := "http://localhost:57493/customer/" + customer + "/subscription"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil
	}
	defer resp.Body.Close()

	subscriptions := &stripe.SubscriptionList{}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(subscriptions)
	return subscriptions
}
