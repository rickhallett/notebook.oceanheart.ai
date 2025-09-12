package http

import (
	"bytes"
	"compress/gzip"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap with logging middleware
	loggedHandler := LoggingMiddleware(handler)

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	loggedHandler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check log output
	logOutput := buf.String()
	if !strings.Contains(logOutput, "GET /test 200") {
		t.Errorf("Expected log to contain request info, got: %s", logOutput)
	}
}

func TestGzipMiddleware(t *testing.T) {
	// Create test handler
	testContent := "This is test content that should be compressed by gzip middleware"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testContent))
	})

	// Wrap with gzip middleware
	gzipHandler := GzipMiddleware(handler)

	// Test with gzip support
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	w := httptest.NewRecorder()

	gzipHandler.ServeHTTP(w, req)

	// Check headers
	if w.Header().Get("Content-Encoding") != "gzip" {
		t.Error("Expected Content-Encoding: gzip header")
	}

	if w.Header().Get("Vary") != "Accept-Encoding" {
		t.Error("Expected Vary: Accept-Encoding header")
	}

	// Check compressed content
	if w.Body.Len() == 0 {
		t.Error("Expected compressed content")
	}

	// Decompress and verify content
	gz, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gz.Close()

	var decompressed bytes.Buffer
	_, err = decompressed.ReadFrom(gz)
	if err != nil {
		t.Fatalf("Failed to decompress content: %v", err)
	}

	if decompressed.String() != testContent {
		t.Errorf("Decompressed content mismatch: got %s, expected %s", 
			decompressed.String(), testContent)
	}

	// Test without gzip support
	req = httptest.NewRequest("GET", "/test", nil)
	// No Accept-Encoding header
	w = httptest.NewRecorder()

	gzipHandler.ServeHTTP(w, req)

	if w.Header().Get("Content-Encoding") != "" {
		t.Error("Expected no Content-Encoding header without gzip support")
	}

	if w.Body.String() != testContent {
		t.Error("Expected uncompressed content when gzip not supported")
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with security headers middleware
	secureHandler := SecurityHeadersMiddleware(handler)

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	secureHandler.ServeHTTP(w, req)

	// Check security headers
	expectedHeaders := map[string]string{
		"X-Content-Type-Options":           "nosniff",
		"X-Frame-Options":                  "DENY",
		"X-XSS-Protection":                 "1; mode=block",
		"Referrer-Policy":                  "strict-origin-when-cross-origin",
	}

	for header, expectedValue := range expectedHeaders {
		if w.Header().Get(header) != expectedValue {
			t.Errorf("Expected %s: %s, got %s", 
				header, expectedValue, w.Header().Get(header))
		}
	}
}

func TestCacheHeadersMiddleware(t *testing.T) {
	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with cache headers middleware
	cacheHandler := CacheHeadersMiddleware(handler)

	// Test static asset
	req := httptest.NewRequest("GET", "/static/app.css", nil)
	w := httptest.NewRecorder()

	cacheHandler.ServeHTTP(w, req)

	if w.Header().Get("Cache-Control") != "public, max-age=31536000" {
		t.Errorf("Expected long cache for static assets, got %s", 
			w.Header().Get("Cache-Control"))
	}

	// Test dynamic content
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	cacheHandler.ServeHTTP(w, req)

	if w.Header().Get("Cache-Control") != "public, max-age=300" {
		t.Errorf("Expected short cache for dynamic content, got %s", 
			w.Header().Get("Cache-Control"))
	}
}

func TestChainMiddleware(t *testing.T) {
	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	// Create middleware that adds headers
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-1", "value1")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-2", "value2")
			next.ServeHTTP(w, r)
		})
	}

	// Chain middleware
	chained := ChainMiddleware(handler, middleware1, middleware2)

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	chained.ServeHTTP(w, req)

	// Check that both middleware ran
	if w.Header().Get("X-Test-1") != "value1" {
		t.Error("Expected X-Test-1 header from middleware1")
	}

	if w.Header().Get("X-Test-2") != "value2" {
		t.Error("Expected X-Test-2 header from middleware2")
	}

	if w.Body.String() != "test" {
		t.Error("Expected response body to be preserved")
	}
}