package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ikirja/easy-web-metrics-go/internal/database"
	"github.com/ikirja/easy-web-metrics-go/internal/metrics"
)

func handlerPing(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Message string `json:"message"`
	}{
		Message: "ping",
	}
	responseWithJson(w, 200, body)
}

func handlerProcessVisitor(w http.ResponseWriter, r *http.Request) {
	vr, err := metrics.ProcessVisitor(r)
	if err != nil {
		responseWithError(w, 400, "error proccessing visitor")
		return
	}
	responseWithJson(w, 200, vr)
}

func handlerGetVisitors(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Limit int64 `json:"limit"`
		Skip  int64 `json:"skip"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responseWithError(w, 400, "error parsing body")
		return
	}
	if body.Limit == 0 || body.Limit > 5000 {
		responseWithError(w, 400, "limit can't be 0 or more than 5000")
		return
	}

	v, err := metrics.GetVisitors(body.Limit, body.Skip)
	if err != nil {
		responseWithError(w, 400, "error getting visitors")
		return
	}

	result := struct {
		Visitors []database.VisitorDB `json:"visitors"`
	}{
		Visitors: v,
	}
	responseWithJson(w, 200, result)
}

func responseWithError(w http.ResponseWriter, code int, message string) {
	data, err := json.Marshal(struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   true,
		Message: message,
	})
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func responseWithJson(w http.ResponseWriter, code int, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}
