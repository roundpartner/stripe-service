package client

import (
	"bytes"
	"encoding/json"
	"github.com/stripe/stripe-go"
	"log"
	"net/http"
	"strings"
)

type SubscriptionItem struct {
	Status           string       `json:"status"`
	CustomerStatus   string       `json:"customer_status"`
	DaysUntilDue     int          `json:"days_until_due"`
	CurrentPeriodEnd int          `json:"current_period_end"`
	Plan             PlanItem     `json:"plan,omitempty"`
	Plans            []PlanItem   `json:"plans"`
	Items            *RawPlanItem `json:"items,omitempty"`
}

type RawPlanItem struct {
	Plans []SubscriptionItem `json:"data"`
}

type PlanItem struct {
	Id     string `json:"id,omitempty"`
	PlanId string `json:"plan_id"`
	Name   string `json:"nickname"`
	Amount int64  `json:"amount"`
}

type SessionItem struct {
	Id         string           `json:"session_id"`
	CustomerId string           `json:"customer_id"`
	Plan       map[int]PlanItem `json:"plan"`
	Amount     int64            `json:"amount"`
}

func Subscription(customer string) []*SubscriptionItem {
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

	var subscriptions []*SubscriptionItem

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&subscriptions); err != nil {
		log.Printf("[ERROR] Decode error: %s", err.Error())
		return nil
	}

	for key := range subscriptions {
		if subscriptions[key].Items == nil {
			log.Printf("[INFO] No plans found in subscription")
			continue
		}
		for subkey := range subscriptions[key].Items.Plans {
			subscriptions[key].Items.Plans[subkey].Plan.PlanId = subscriptions[key].Items.Plans[subkey].Plan.Id
			subscriptions[key].Items.Plans[subkey].Plan.Id = ""
			subscriptions[key].Plans = append(subscriptions[key].Plans, subscriptions[key].Items.Plans[subkey].Plan)
		}
		subscriptions[key].Items = nil
		subscriptions[key].CustomerStatus = strings.Title(subscriptions[key].Status)
	}

	return subscriptions
}

func Session(customer string, plan []string) *SessionItem {
	client := &http.Client{}
	url := "http://localhost:57493/v2/customer/" + customer + "/session"

	body, err := json.Marshal(plan)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
	if err := decoder.Decode(&session); err != nil {
		log.Printf("[ERROR] Decoding session response: %s", err.Error())
		return nil
	}

	if session.ID == "" {
		return nil
	}

	totalAmount := int64(0)

	planItems := map[int]PlanItem{}
	for index, plan := range session.DisplayItems {
		planItems[index] = PlanItem{
			PlanId: plan.Plan.ID,
			Name:   plan.Plan.Nickname,
			Amount: plan.Amount,
		}
		totalAmount += plan.Amount
	}

	return &SessionItem{
		Id:         session.ID,
		CustomerId: session.Customer.ID,
		Plan:       planItems,
		Amount:     totalAmount,
	}
}
