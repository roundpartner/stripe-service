package main

import (
	"flag"
	"github.com/artyom/autoflags"
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
	ListenAndServe(ServerConfig.Port)
}
