package client

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"log"
	"net/http"
)

type SubscriptionItem struct {
	Status       string   `json:"status"`
	DaysUntilDue int      `json:"days_until_due"`
	Plan         PlanItem `json:"plan"`
}

type PlanItem struct {
	Name string `json:"nickname"`
}

type SessionItem struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	PlanId     string `json:"plan_id"`
}

func Subscription(customer string) []SubscriptionItem {
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

	var subscriptions []SubscriptionItem

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&subscriptions)
	return subscriptions
}

func Session(customer string, plan string) *SessionItem {
	client := &http.Client{}
	url := "http://localhost:57493/customer/" + customer + "/session/" + plan
	req, err := http.NewRequest("POST", url, nil)
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

	session := stripe.CheckoutSession{}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&session)
	return &SessionItem{
		Id:         session.ID,
		CustomerId: session.Customer.ID,
		PlanId:     session.DisplayItems[0].Plan.ID,
	}
}
