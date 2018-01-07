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

func TestNewCustomerParam(t *testing.T) {
	cr := &CustomerRequest{
		Desc:     "Description",
		Email:    "Email@Address.com",
		Account:  "123",
		User:     "456",
		Discount: "78",
	}
	c := NewCustomerParam(cr)

	if "Description" != c.Desc {
		t.FailNow()
	}

	if "Email@Address.com" != c.Email {
		t.FailNow()
	}

	if "123" != c.Meta["account"] {
		t.FailNow()
	}

	if "456" != c.Meta["user"] {
		t.FailNow()
	}

	if "78" != c.Meta["discount"] {
		t.FailNow()
	}
}

func TestCustomers(t *testing.T) {
	body := strings.NewReader("{\"limit\": \"100\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customer", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %s", rr.Code)
		t.FailNow()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.FailNow()
	}

	customer := &stripe.CustomerList{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&customer.Values)
	if nil != err {
		t.Error(err.Error())
		t.Error(rr.Body.String())
		t.FailNow()
	}

	if len(customer.Values) == 0 {
		t.Skipf("%d values returned instead of 1", len(customer.Values))
	}

	t.Log(rr.Body.String())
}

func TestGetCustomer(t *testing.T) {
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customer/cus_BRsEJtkXRxHxPU", nil)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %s", rr.Code)
		t.FailNow()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.FailNow()
	}

	customer := &stripe.Customer{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(customer)
	if nil != err {
		t.Error(err.Error())
		t.Error(rr.Body.String())
		t.FailNow()
	}

}

func TestNewCustomer(t *testing.T) {
	body := strings.NewReader("{\"token\": \"tok_gb\", \"account\": \"123\", \"user\": \"456\", \"email\": \"gotest@mailinator.com\", \"desc\": \"Added by go test\", \"discount\": \"30\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/customer", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %s", rr.Code)
		t.FailNow()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.FailNow()
	}

	customer := &stripe.Customer{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(customer)
	if nil != err {
		t.Error(err.Error())
		t.Error(rr.Body.String())
		t.FailNow()
	}

	delete(customer.ID)
}

func TestNewCustomerWithoutCard(t *testing.T) {
	body := strings.NewReader("{\"account\": \"123\", \"user\": \"456\", \"email\": \"gotest@mailinator.com\", \"desc\": \"Added by go test\", \"discount\": \"30\"}")
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/customer", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %s", rr.Code)
		t.FailNow()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.FailNow()
	}

	customer := &stripe.Customer{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(customer)
	if nil != err {
		t.Error(err.Error())
		t.Error(rr.Body.String())
		t.FailNow()
	}

	delete(customer.ID)
}

func TestReloadCustomers(t *testing.T) {
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reload", nil)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %s", rr.Code)
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
		t.FailNow()
	}
}
