package main

import (
	"bytes"
	"github.com/roundpartner/stripe-service/util"
	"github.com/stripe/stripe-go"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebHook(t *testing.T) {
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhook", nil)
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

func TestNewSNSService(t *testing.T) {
	NewSNSService()
}

func TestPush(t *testing.T) {
	buf := bytes.NewBufferString("hello world")

	snsService := NewSNSService()
	err := snsService.Push(buf)
	if err != nil {
		t.Errorf("Unexpected error returned: %s", err.Error())
	}
}
