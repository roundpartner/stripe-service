FROM golang

COPY stripe-service /bin/stripe-service

CMD ["stripe-service"]