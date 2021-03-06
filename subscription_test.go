package main

import (
	"bytes"
	"encoding/json"
	"github.com/roundpartner/stripe-service/util"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestServer_GetCustomerSessionV2(t *testing.T) {
	stripe.Key = util.GetTestKey()
	if err := util.SetSubscriptionEnvironmentVariables(); err != nil {
		t.Fatalf("Unable to setup test environment: %s", err.Error())
	}

	body := `["plan_FPSDCc5aQKEEP3", "plan_FrDrMXuQmKGoIP"]`
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v2/customer/cus_C3MQXNRknd5e6p/session", bytes.NewBufferString(body))
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %d", rr.Code)
		t.Errorf("body: %s", rr.Body.String())
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
		t.FailNow()
	}

	session := &stripe.CheckoutSession{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&session)
	if err != nil {
		t.Fatalf("Unable to decode session data")
	}

	if session.Customer.ID != "cus_C3MQXNRknd5e6p" {
		t.Errorf("Unexpected customer returned")
	}

	if len(session.DisplayItems) != 2 {
		t.Errorf("Unexpected number of items returned %d instead of 2", len(session.DisplayItems))
	}

}

func TestRestServer_UpgradeSubscription(t *testing.T) {
	stripe.Key = util.GetTestKey()
	if err := util.SetSubscriptionEnvironmentVariables(); err != nil {
		t.Fatalf("Unable to setup test environment: %s", err.Error())
	}

	body := `["plan_FrDrMXuQmKGoIP"]`

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/customer/cus_FspJN56DogyJmY/subscription", bytes.NewBufferString(body))
	rs := New()
	if err := rs.RemovePlans("cus_FspJN56DogyJmY", []string{"plan_FrDrMXuQmKGoIP"}); err != nil {
		t.Fatalf("Unable to remove plan: %s", err.Error())
	}

	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("wrong error code returned: %d", rr.Code)
		t.Errorf("body: %s", rr.Body.String())
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
		t.FailNow()
	}
}

func TestRestServer_DowngradeSubscription(t *testing.T) {
	stripe.Key = util.GetTestKey()
	if err := util.SetSubscriptionEnvironmentVariables(); err != nil {
		t.Fatalf("Unable to setup test environment: %s", err.Error())
	}

	body := `["plan_FrDrMXuQmKGoIP"]`

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/customer/cus_FspJN56DogyJmY/subscription", bytes.NewBufferString(body))
	rs := New()

	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("wrong error code returned: %d", rr.Code)
		t.Errorf("body: %s", rr.Body.String())
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
		t.FailNow()
	}
}

func TestRestServer_CancelSubscription(t *testing.T) {
	stripe.Key = util.GetTestKey()
	if err := util.SetSubscriptionEnvironmentVariables(); err != nil {
		t.Fatalf("Unable to setup test environment: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/customer/cus_FspJN56DogyJmY/cancel", nil)
	rs := New()

	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("wrong error code returned: %d", rr.Code)
		t.Errorf("body: %s", rr.Body.String())
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
		t.FailNow()
	}

	subItem, err := sub.Get("sub_FspWjqCT5ep2VF", nil)
	if err != nil {
		t.Fatalf("Error getting subscription: %s", err.Error())
	}

	if !subItem.CancelAtPeriodEnd {
		t.Fatalf("Subscription not canceled")
	}

	subParams := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}
	_, err = sub.Update("sub_FspWjqCT5ep2VF", subParams)
	if err != nil {
		t.Fatalf("Error reactivating subscription: %s", err.Error())
	}
}
