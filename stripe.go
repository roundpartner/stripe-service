package main

import (
	"github.com/stripe/stripe-go"
	"os"
)

func initStripe() {
	stripe.Key = getStripeKey()
}

func getStripeKey() string {
	key := os.Getenv("STRIPE_KEY")

	if len(key) == 0 {
		panic("STRIPE_KEY environment variable is not set, but is needed to start server!\n")
	}

	return key
}
