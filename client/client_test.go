package client

import (
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestSubscription(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:57493").
		Get("/customer/cus_12345/subscription").
		Reply(http.StatusOK).
		BodyString(`[{"status":"active","days_until_due": 0,"current_period_end":1571246262,"plan":{"id":"plan_12345","nickname":"Plan"}}]`)

	subscription := Subscription("cus_12345")

	if !gock.IsDone() {
		t.Errorf("Mocked http was not called")
	}

	if len(subscription) != 1 {
		t.Fatalf("Unexpected total count: %d", len(subscription))
	}

	if subscription[0].Status != "active" {
		t.Errorf("Unexpected status for subscription")
	}

	if subscription[0].DaysUntilDue != 0 {
		t.Errorf("Unexpected due date for invoice")
	}

	if subscription[0].CurrentPeriodEnd != 1571246262 {
		t.Errorf("Unexpected current period end date for subscription")
	}

	if subscription[0].Plan.Name != "Plan" {
		t.Errorf("Unexpected plan")
	}

	if subscription[0].Plan.PlanId != "plan_12345" {
		t.Errorf("Unexpected plan id")
	}
}

func TestSession(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:57493").
		Post("/v2/customer/cus_12345/session").
		Reply(http.StatusOK).
		BodyString(`{"cancel_url":"https://example/cancel","client_reference_id":"","customer":{"address":{"city":"","country":"","line1":"","line2":"","postal_code":"","state":""},"balance":0,"created":0,"currency":"","default_source":null,"deleted":false,"delinquent":false,"description":"","discount":null,"email":"","id":"cus_12345","invoice_prefix":"","invoice_settings":null,"livemode":false,"metadata":null,"name":"","phone":"","preferred_locales":null,"shipping":null,"sources":null,"subscriptions":null,"tax_exempt":"","tax_ids":null,"account_balance":0,"tax_info":null,"tax_info_verification":null},"customer_email":"","deleted":false,"display_items":[{"amount":1000,"currency":"gbp","custom":null,"quantity":1,"plan":{"active":true,"aggregate_usage":"","amount":1000,"billing_scheme":"per_unit","created":1562755534,"currency":"gbp","deleted":false,"id":"plan_12345","interval":"month","interval_count":1,"livemode":false,"metadata":{},"nickname":"monthly","product":{"active":false,"attributes":null,"caption":"","created":0,"deactivate_on":null,"description":"","id":"prod_FPSAW8eylXpMIS","images":null,"livemode":false,"metadata":null,"name":"","package_dimensions":null,"shippable":false,"statement_descriptor":"","type":"","unit_label":"","url":"","updated":0},"tiers":null,"tiers_mode":"","transform_usage":null,"trial_period_days":0,"usage_type":"licensed"},"sku":null,"type":"plan"}],"id":"cs_test_1234","livemode":false,"locale":"","object":"checkout.session","payment_intent":null,"payment_method_types":["card"],"subscription":null,"submit_type":"","success_url":"https://example/success"}`)

	plan := []string{"plan_12345"}
	session := Session("cus_12345", plan)

	if !gock.IsDone() {
		t.Errorf("Mocked http was not called")
	}

	if session == nil {
		t.Fatalf("Session not returned")
	}

	if session.CustomerId != "cus_12345" {
		t.Errorf("Unexpected customer id")
	}

	if session.Plan[0].PlanId != "plan_12345" {
		t.Errorf("Unexpected plan id")
	}

	if session.Plan[0].Amount != 1000 {
		t.Errorf("Unexpected amount")
	}

	if session.Amount != 1000 {
		t.Errorf("Unexpected amount total")
	}
}

func TestSessionError(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:57493").
		Post("/v2/customer/cus_12345/session").
		Reply(http.StatusBadRequest).
		BodyString(`{"error":{"code":"resource_missing","status":400,"message":"No such plan","type":"invalid_request_error"}`)

	plan := []string{"plan_12345"}
	session := Session("cus_12345", plan)

	if !gock.IsDone() {
		t.Errorf("Mocked http was not called")
	}

	if session != nil {
		t.Errorf("Expected no session to be returned")
	}
}
