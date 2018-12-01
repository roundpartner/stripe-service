package main

import (
	"flag"
	"github.com/artyom/autoflags"
	"github.com/roundpartner/stripe-service/http"
	"log"
	"os"
)

var ServerConfig = struct {
	Port int `flag:"port,port number to listen on"`
}{
	Port: 57493,
}

func main() {
	log.SetOutput(os.Stdout)
	autoflags.Define(&ServerConfig)
	flag.Parse()
	initStripe()
	http.ListenAndServe(ServerConfig.Port)
}
