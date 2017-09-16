package http

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCharge(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_visa\", \"amount\": 720, \"desc\": \"example\"}")
	stripe.Key = utils.GetTestKey()
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
