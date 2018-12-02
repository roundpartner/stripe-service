package main

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

func TestThatServiceIsDown(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/check", nil)

	rs := New()
	rs.router().ServeHTTP(rr, req)

	serviceAvailable = false

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Service did not return ok no content status")
		t.FailNow()
	}
}

func TestGetMetrics(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)

	rs := New()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("Service did not return ok no content status")
		t.FailNow()
	}
}
