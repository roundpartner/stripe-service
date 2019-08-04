package client

import (
	"encoding/json"
	"log"
	"net/http"
)

type SubscriptionItem struct {
	Status string `json:"status"`
	DaysUntilDue int `json:"days_until_due"`
	Plan PlanItem `json:"plan"`
}

type PlanItem struct {
	Name string `json:"nickname"`
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
