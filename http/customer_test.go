package http

import (
	"testing"
	"github.com/stripe/stripe-go"
	"net/http/httptest"
	"net/http"
	"github.com/roundpartner/stripe-service/util"
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
