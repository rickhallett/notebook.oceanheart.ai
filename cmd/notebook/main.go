package main

import (
	"fmt"
	"log"
	"net/http"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/store"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := store.MustOpen(cfg.DBPath)
	defer db.Close()

	// Setup HTTP server
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/healthz", healthCheckHandler)

	// TODO: Add other handlers in Phase 02

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s (env=%s)", addr, cfg.Environment)
	
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(fmt.Errorf("server failed: %w", err))
	}
}

// healthCheckHandler returns basic health check response
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"notebook.oceanheart.ai"}`))
}