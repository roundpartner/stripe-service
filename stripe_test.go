package main

import (
	"github.com/stripe/stripe-go"
	"strings"
	"testing"
)

func TestInitStripe(t *testing.T) {
	if initStripe(); !strings.HasPrefix(stripe.Key, "sk_test_") {
		t.Fail()
	}
}

func TestGetStripeKey(t *testing.T) {
	if !strings.HasPrefix(getStripeKey(), "sk_test_") {
		t.Fail()
	}
}
