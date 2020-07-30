FROM alpine

COPY stripe-service stripe-service

RUN ./stripe-service