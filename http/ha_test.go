package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestThatServiceIsUp(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/check", nil)

	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("Service did not return ok no content status")
		t.FailNow()
	}
}
