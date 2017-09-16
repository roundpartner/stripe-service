package main

import (
	"github.com/roundpartner/stripe-service/http"
)

func main() {
	initStripe()
	http.ListenAndServe()
}
