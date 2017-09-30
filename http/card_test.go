package http

import (
	"testing"
	"strings"
	"github.com/stripe/stripe-go"
	"net/http/httptest"
	"net/http"
	"github.com/roundpartner/stripe-service/util"
	"encoding/json"
)

func TestRestServer_UpdateCustomerCard(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_mastercard\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/customer/cus_BUn5fLay0cMkMw/card", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %s", rr.Code)
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
