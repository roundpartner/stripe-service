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
    -d "{\"token\": \"tok_gb\", \"amount\": 1000, \"desc\": \"example\"}" \
    http://0.0.0.0:57493/charge
```
## Customer
### List
```bash
curl -X GET \
    -d "{\"limit\":\"10\"}" \
    http://0.0.0.0:57493/customer
```
### Get
The customer id will return the customer details
```bash
curl http://0.0.0.0:57493/customer/cus_BUoP6KtXPL3ajU
```
### Add
```bash
curl -X POST \
    -d "{\"token\": \"tok_gb\", \"account\": \"1\", \"email\": \"example@mailinator.com\", \"desc\": \"Added by go test\"}" \
    http://0.0.0.0:57493/customer
```
### New Default Card
```bash
curl -X PUT \
    -d "{\"token\": \"tok_mastercard_debit\"}" \
    http://0.0.0.0:57493/customer/cus_BUoP6KtXPL3ajU/card
```
### Reload
```bash
curl http://0.0.0.0:57493/reload
```
### Subscriptions
Get customer subscriptions
```bash
curl http://0.0.0.0:57493/customer/cus_DOQj7OGOt6mX1n/subscription
```
### Sessions
```bash
curl -X POST http://0.0.0.0:57493/customer/cus_BUoP6KtXPL3ajU/session/plan_FPSDCc5aQKEEP3
```
