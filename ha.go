package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func AddCheck(router *mux.Router) {
	router.HandleFunc("/check", Check).Methods("GET")
}

func Check(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
