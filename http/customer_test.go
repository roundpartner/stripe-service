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

	if "Description" != *c.Description {
		t.FailNow()
	}

	if "Email@Address.com" != *c.Email {
		t.FailNow()
	}

	if "123" != c.Metadata["account"] {
		t.FailNow()
	}

	if "456" != c.Metadata["user"] {
		t.FailNow()
	}

	if "78" != c.Metadata["discount"] {
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
		t.Errorf("wrong error code returned: %d", rr.Code)
		t.FailNow()
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.FailNow()
	}

	customer := &stripe.CustomerList{}
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&customer.Data)
	if nil != err {
		t.Error(err.Error())
		t.Error(rr.Body.String())
		t.FailNow()
	}

	if len(customer.Data) == 0 {
		t.Skipf("%d values returned instead of 1", len(customer.Data))
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
		t.Errorf("wrong error code returned: %d", rr.Code)
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
		t.Errorf("wrong error code returned: %d", rr.Code)
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
		t.Errorf("wrong error code returned: %d", rr.Code)
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

func TestUpdateCustomer(t *testing.T) {
	stripe.Key = util.GetTestKey()
	body := strings.NewReader("{\"email\": \"gotest@mailinator.com\"}")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/customer/cus_C3MQXNRknd5e6p", body)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong error code returned: %d", rr.Code)
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
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

func TestChargeCustomer(t *testing.T) {
	body := strings.NewReader("{\"customer\": \"cus_C3MQXNRknd5e6p\", \"amount\": 720, \"desc\": \"example\", \"trans_id\": \"tnx_1234\", \"business_name\": \"RoundPartner\"}")
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

func TestReloadCustomers(t *testing.T) {
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reload", nil)
	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("wrong error code returned: %d", rr.Code)
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
		t.FailNow()
	}
}
