package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type InfoResponse struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Hostname    string `json:"hostname"`
	Timestamp   string `json:"timestamp"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	env     := getEnv("APP_ENV", "dev")
	version := getEnv("APP_VERSION", "0.0.1")
	port    := getEnv("APP_PORT", "8080")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	})

	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(InfoResponse{
			Version:     version,
			Environment: env,
			Hostname:    hostname,
			Timestamp:   time.Now().Format(time.RFC3339),
		})
	})

	log.Printf("Servidor iniciado — env=%s version=%s port=%s", env, version, port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}