[![Build Status](https://travis-ci.org/roundpartner/stripe-service.svg?branch=master)](https://travis-ci.org/roundpartner/stripe-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/roundpartner/stripe-service)](https://goreportcard.com/report/github.com/roundpartner/stripe-service)
# Stripe Micro Service
A Micro Service for Stripe Payments in Go

# Building
```bash
go build
```

# Usage
```bash
export STRIPE_KEY="your stripe key"
./stripe-service
```
## Charge
To take a single payment the charge end point provides this
```bash
curl -X POST\
    -d "{\"token\": \"tok_visa\", \"amount\": 1000, \"desc\": \"example\"}" \
    http://0.0.0.0:57493/charge
```
