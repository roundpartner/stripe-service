package main

import (
	"github.com/roundpartner/stripe-service/http"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	initStripe()
	http.ListenAndServe()
}
