package main

import (
	"fmt"
	"log"
	"net/http"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/content"
	httpserver "notebook.oceanheart.ai/internal/http"
	"notebook.oceanheart.ai/internal/store"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := store.MustOpen(cfg.DBPath)
	defer db.Close()

	// Load content from filesystem and cache in database
	loader := content.NewLoader(cfg.ContentDir)
	posts, err := loader.LoadAll()
	if err != nil {
		log.Printf("Warning: Failed to load content: %v", err)
		posts = []*store.Post{} // Continue with empty posts
	}

	// Cache posts in database
	if len(posts) > 0 {
		if err := db.UpsertPosts(posts); err != nil {
			log.Printf("Warning: Failed to cache posts: %v", err)
		} else {
			log.Printf("Loaded and cached %d posts", len(posts))
		}
	}

	// Create HTTP server with handlers
	server := httpserver.NewServer(db, cfg)

	// Setup routes
	mux := http.NewServeMux()

	// Main routes
	mux.HandleFunc("/", server.HomeHandler)
	mux.HandleFunc("/p/", server.PostHandler)
	mux.HandleFunc("/tag/", server.TagHandler)
	mux.HandleFunc("/static/", server.StaticHandler)
	mux.HandleFunc("/static/chroma.css", server.ChromaCSSHandler)

	// SEO and feed routes
	mux.HandleFunc("/feed.xml", server.FeedHandler)
	mux.HandleFunc("/sitemap.xml", server.SitemapHandler)

	// Admin/reload endpoint (dev allowed; prod requires token)
	mux.HandleFunc("/admin/reload", server.AdminReloadHandler)

	// Health check endpoint
	mux.HandleFunc("/healthz", healthCheckHandler)

	// Apply middleware chain
	handler := httpserver.ChainMiddleware(mux,
		httpserver.SecurityHeadersMiddleware,
		httpserver.LoggingMiddleware,
		httpserver.GzipMiddleware,
		httpserver.CacheHeadersMiddleware,
	)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s (env=%s)", addr, cfg.Environment)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
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
