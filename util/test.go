package util

import (
	"errors"
	"os"
)

func GetTestKey() string {
	key := os.Getenv("STRIPE_KEY")

	if len(key) == 0 {
		panic("STRIPE_KEY environment variable is not set, but is needed to run tests!\n")
	}

	return key
}

func SetSubscriptionEnvironmentVariables() error {
	if err := os.Setenv("STRIPE_SUCCESS_URL", "https://example/success"); err != nil {
		return errors.New("unable to set environment")
	}
	if err := os.Setenv("STRIPE_CANCEL_URL", "https://example/cancel"); err != nil {
		return errors.New("unable to set environment")
	}
	return nil
}
