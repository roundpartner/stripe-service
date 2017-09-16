package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"net/http"
)

func ListenAndServe() {
	rs := New()
	http.ListenAndServe(":57493", rs.router())
}

type RestServer struct {
}

func (rs *RestServer) router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/charge", rs.Charge).Methods("POST")
	return router
}

func New() *RestServer {
	return &RestServer{}
}

type ChargeRequest struct {
	Token  string `json:"token"`
	Amount uint64 `json:"amount"`
	Desc   string `json:"desc"`
}

func (rs *RestServer) Charge(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	t := &ChargeRequest{}
	decoder.Decode(t)
	defer req.Body.Close()

	token := t.Token
	params := &stripe.ChargeParams{
		Amount:   t.Amount,
		Currency: "gbp",
		Desc:     t.Desc,
	}
	params.SetSource(token)

	charge, err := charge.New(params)

	if err != nil {
		InternalError(w, err.Error())
		return
	}

	js, _ := json.Marshal(charge)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
