package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ikirja/easy-web-metrics-go/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	ApiKey string
}

type EnvVars struct {
	Port   string
	ApiKey string
	DBName string
}

var envVars = EnvVars{
	Port:   "PORT",
	ApiKey: "API_KEY",
	DBName: "DB_NAME",
}

func main() {
	godotenv.Load()

	portString := os.Getenv(envVars.Port)
	if portString == "" {
		log.Fatal("PORT env variable not provided")
	}

	apiKey := os.Getenv(envVars.ApiKey)
	if apiKey == "" {
		log.Fatal("API_KEY env variable not provided")
	}

	dbName := os.Getenv(envVars.DBName)
	if dbName == "" {
		log.Fatal("DB_NAME env variable not provided")
	}

	err := database.Connect(dbName)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := ApiConfig{
		ApiKey: apiKey,
	}

	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(apiCfg.middlewareCheckApiKey)
		r.Get("/ping", handlerPing)

		r.Route("/metrics", func(r chi.Router) {
			r.Post("/visitor", handlerProcessVisitor)
		})
	})

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("[server]: starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
