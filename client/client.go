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
	Id     string `json:"id,omitempty"`
	PlanId string `json:"plan_id"`
	Name   string `json:"nickname"`
	Amount int64  `json:"amount"`
}

type SessionItem struct {
	Id         string   `json:"session_id"`
	CustomerId string   `json:"customer_id"`
	Plan       PlanItem `json:"plan"`
}

func Subscription(customer string) []SubscriptionItem {
	client := &http.Client{}
	url := "http://localhost:57493/customer/" + customer + "/subscription"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] %s", "service returned non ok status")
		return nil
	}
	defer resp.Body.Close()

	var subscriptions []SubscriptionItem

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&subscriptions)

	for key := range subscriptions {
		subscriptions[key].Plan.PlanId = subscriptions[key].Plan.Id
		subscriptions[key].Plan.Id = ""
	}

	return subscriptions
}

func Session(customer string, plan string) *SessionItem {
	client := &http.Client{}
	url := "http://localhost:57493/customer/" + customer + "/session/" + plan
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] %s", "service returned non ok status")
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	session := stripe.CheckoutSession{}
	decoder.Decode(&session)
	return &SessionItem{
		Id:         session.ID,
		CustomerId: session.Customer.ID,
		Plan: PlanItem{
			PlanId: session.DisplayItems[0].Plan.ID,
			Name:   session.DisplayItems[0].Plan.Nickname,
			Amount: session.DisplayItems[0].Plan.Amount,
		},
	}
}
