package main

import (
	"encoding/json"
	"log"
	"net/http"

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
