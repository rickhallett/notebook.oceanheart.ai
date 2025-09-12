package http

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoggingMiddleware logs HTTP requests with timing information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap ResponseWriter to capture status code
		lw := &loggingWriter{ResponseWriter: w, statusCode: 200}
		
		// Process request
		next.ServeHTTP(lw, r)
		
		// Log request
		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, lw.statusCode, duration)
	})
}

// loggingWriter wraps http.ResponseWriter to capture status code
type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

// GzipMiddleware compresses responses when client supports it
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client accepts gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Skip compression for certain content types
		contentType := w.Header().Get("Content-Type")
		if strings.Contains(contentType, "image/") || 
		   strings.Contains(contentType, "video/") ||
		   strings.Contains(contentType, "application/zip") {
			next.ServeHTTP(w, r)
			return
		}

		// Create gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Set headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// Wrap response writer
		gzw := &gzipWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(gzw, r)
	})
}

// gzipWriter wraps http.ResponseWriter to compress response
type gzipWriter struct {
	http.ResponseWriter
	io.Writer
}

func (gw *gzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}

// SecurityHeadersMiddleware adds basic security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		next.ServeHTTP(w, r)
	})
}

// CacheHeadersMiddleware adds appropriate cache headers
func CacheHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cache static assets longer
		if strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
		} else {
			// Cache dynamic content briefly
			w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
		}
		
		next.ServeHTTP(w, r)
	})
}

// ChainMiddleware applies multiple middleware functions in order
func ChainMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}