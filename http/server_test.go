package http

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"github.com/roundpartner/stripe-service/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCharge(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_visa\", \"amount\": 720, \"desc\": \"example\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/charge", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fail()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fail()
	}

	charge := &stripe.Charge{}

	decoder := json.NewDecoder(rr.Body)
	decoder.Decode(charge)

	if 720 != charge.Amount {
		t.Fail()
	}

	if "succeeded" != charge.Status {
		t.Fail()
	}

	if false == charge.Paid {
		t.Fail()
	}

}

func TestChargeDecimalFails(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_visa\", \"amount\": 999.99, \"desc\": \"example\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/charge", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fail()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fail()
	}

	if "{\"error\":\"json: cannot unmarshal number 999.99 into Go struct field ChargeRequest.amount of type uint64\"}" != rr.Body.String() {
		t.Error(rr.Body.String())
		t.Fail()
	}

}

func TestChargeLowAmountFails(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_visa\", \"amount\": 29, \"desc\": \"example\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/charge", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fail()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fail()
	}

}
