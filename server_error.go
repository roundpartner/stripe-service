package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type message struct {
	Error string `json:"error"`
}

type messageJson struct {
	Error interface{} `json:"error"`
}

func StripeError(w http.ResponseWriter, msg string) {
	log.Printf("stripe error: %s", msg)
	b := bytes.NewBufferString(msg).Bytes()
	var i interface{}
	json.Unmarshal(b, &i)
	js, _ := json.Marshal(messageJson{Error: i})
	writeJsonError(w, js, http.StatusBadRequest)
}

func InternalError(w http.ResponseWriter, msg string) {
	log.Printf("server error: %s", msg)
	jsonError(w, msg, http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, msg string) {
	log.Printf("client error: %s", msg)
	jsonError(w, msg, http.StatusBadRequest)
}

func jsonError(w http.ResponseWriter, error string, code int) {
	js, err := json.Marshal(message{Error: error})
	if err != nil {
		js = bytes.NewBufferString("{\"error\":\"Marshal Error\"}").Bytes()
		writeJsonError(w, js, http.StatusInternalServerError)
		return
	}
	writeJsonError(w, js, code)
}

func writeJsonError(w http.ResponseWriter, error []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write(error)
}
