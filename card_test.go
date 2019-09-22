package main

import (
	"encoding/json"
	"github.com/roundpartner/stripe-service/util"
	"github.com/stripe/stripe-go"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteAllCards(t *testing.T) {
	t.Skipf("Preventing deleting all cards from test customer")
	err := DeleteAllCards("cus_C3MQXNRknd5e6p")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRestServer_UpdateCustomerCard(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_gb\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/customer/cus_C3MQXNRknd5e6p/card", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %d", rr.Code)
		t.Fail()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fail()
	}

	customer := &stripe.Customer{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(customer)
	if nil != err {
		t.Error(err.Error())
		t.Error(rr.Body.String())
		t.Fail()
	}
}

func TestUpdatingCustomerCardWithoutSource(t *testing.T) {
	body := strings.NewReader("{\"token\": \"\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/customer/cus_C3MQXNRknd5e6p/card", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("wrong error code returned: %d", rr.Code)
		t.Fail()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("content type was not set correctly")
		t.FailNow()
	}

	response := &messageJson{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(response)
	if nil != err {
		t.Error(err.Error())
		t.Error(rr.Body.String())
		t.FailNow()
	}

	if response.Error != "token is required in this request" {
		t.Errorf("Response was: %s", response.Error)
		t.Fail()
	}
}
