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

	if serviceName, isSet := os.LookupEnv("SERVICE_NAME"); isSet {
		ServiceName = serviceName
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR] [%s] %s", ServiceName, err)
		}
	}()

	initStripe()
	ListenAndServe(ServerConfig.Port)
}
