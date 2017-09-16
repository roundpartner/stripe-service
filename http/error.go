package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type message struct {
	Error string `json:"error"`
}

func InternalError(w http.ResponseWriter, msg string) {
	jsonError(w, msg, http.StatusInternalServerError)
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
