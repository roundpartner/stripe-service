package client

import (
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestSubscription(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:57493").
		Get("/customer/cus_12345/subscription").
		Reply(http.StatusOK)

	Subscription("cus_12345")

	if !gock.IsDone() {
		t.Errorf("Mocked http was not called")
	}
}
