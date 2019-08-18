package main

import (
	"bytes"
	"github.com/roundpartner/stripe-service/util"
	"github.com/stripe/stripe-go"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebHook(t *testing.T) {
	stripe.Key = util.GetTestKey()
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhook", bytes.NewBufferString(`test web hook`))
	rs := New()
	rs.SNSService = NewSNSService()
	rs.router().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("wrong error code returned: %d", rr.Code)
	}

	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Errorf("wrong content type returned: %s", rr.Header().Get("Content-Type"))
		t.FailNow()
	}

	if buffer := <-rs.SNSService.Queue; buffer.String() != "test web hook" {
		t.Fatalf("Unexpeced message in queue %s", buffer.String())
	}
}

func TestNewSNSService(t *testing.T) {
	NewSNSService()
}

func TestPush(t *testing.T) {
	defer gock.Off()
	gock.New("https://sns.eu-west-2.amazonaws.com").
		Post("/").
		Reply(http.StatusOK).
		BodyString(`{}`)

	buf := bytes.NewBufferString("hello world")

	snsService := NewSNSService()
	err := snsService.Push(buf)
	if err != nil {
		t.Errorf("Unexpected error returned: %s", err.Error())
	}

	if !gock.IsDone() {
		t.Errorf("Mocked http was not called")
	}
}

func TestRun(t *testing.T) {
	defer gock.Off()
	gock.New("https://sns.eu-west-2.amazonaws.com").
		Post("/").
		Reply(http.StatusOK).
		BodyString(`{}`)

	snsService := NewSNSService()
	snsService.Run()

	buf := bytes.NewBufferString("hello world")
	snsService.Queue <- buf

	time.Sleep(time.Second)

	if !gock.IsDone() {
		t.Errorf("Mocked http was not called")
	}
}
