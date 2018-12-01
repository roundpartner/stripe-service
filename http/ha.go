package http

import (
	"net/http"
	"github.com/gorilla/mux"
)

func AddCheck(router *mux.Router) {
	router.HandleFunc("/check", Check).Methods("GET")
}

func Check(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
