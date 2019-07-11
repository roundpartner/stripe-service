package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ShutdownGracefully(server *http.Server) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		signal.Stop(c)

		serviceAvailable = false

		log.Printf("[INFO] [%s] Server shutting down gracefully", ServiceName)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		err := server.Shutdown(ctx)
		if nil != err {
			log.Printf("[ERROR] [%s] Error shutting down server: %s", ServiceName, err.Error())
		}
	}()
}
