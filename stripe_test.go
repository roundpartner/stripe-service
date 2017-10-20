package main

import (
	"testing"
	"github.com/stripe/stripe-go"
	"strings"
)

func TestInitStrip(t *testing.T) {
	if initStripe(); !strings.HasPrefix(stripe.Key, "sk_test_") {
		t.Fail()
	}
}

func TestGetStripeKey(t *testing.T) {
	if !strings.HasPrefix(getStripeKey(), "sk_test_") {
		t.Fail()
	}
}
