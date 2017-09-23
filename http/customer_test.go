package http

import (
	"encoding/json"
	"github.com/roundpartner/stripe-service/util"
	"github.com/stripe/stripe-go"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetCustomer(t *testing.T) {
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customer/cus_BRsEJtkXRxHxPU", nil)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fail()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fail()
	}

}

func TestNewCustomer(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_gb\", \"account_id\": \"1\", \"email\": \"example@mailinator.com\", \"desc\": \"Added by go test\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/customer", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fail()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fail()
	}

	customer := &stripe.Customer{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(customer)
	if nil != err {
		t.Error(rr.Body.String())
		t.Fail()
	}
}
