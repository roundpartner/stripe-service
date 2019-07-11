package main

import (
	"flag"
	"github.com/artyom/autoflags"
	"log"
	"os"
)

var ServiceName = "stripe"

var ServerConfig = struct {
	Port int `flag:"port,port number to listen on"`
}{
	Port: 57493,
}

func main() {
	log.SetOutput(os.Stdout)
	autoflags.Define(&ServerConfig)
	flag.Parse()

	serviceName, isSet := os.LookupEnv("SERVICE_NAME")
	if isSet {
		ServiceName = serviceName
	}

	initStripe()
	ListenAndServe(ServerConfig.Port)
}
