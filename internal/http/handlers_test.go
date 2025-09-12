package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/store"
)

func TestHomeHandler(t *testing.T) {
	// Create test database
	tempDB := "test_handlers.db"
	defer os.Remove(tempDB)

	db := store.MustOpen(tempDB)
	defer db.Close()

	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		Environment: "test",
	}

	server := NewServer(db, cfg)

	// Test home page with empty database
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.HomeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !contains(body, "Test Blog") {
		t.Error("Expected page to contain site title")
	}

	if !contains(body, "No posts found") {
		t.Error("Expected message for empty database")
	}

	// Test with a post
	post := &store.Post{
		Slug:        "test-post",
		Title:       "Test Post",
		Summary:     "A test post",
		HTML:        "<p>Test content</p>",
		RawMD:       "Test content",
		PublishedAt: "2025-09-12T10:00:00Z",
		UpdatedAt:   "2025-09-12T10:00:00Z",
		Draft:       false,
	}

	err := db.UpsertPost(post)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Test home page with content
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	server.HomeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body = w.Body.String()
	if !contains(body, "Test Post") {
		t.Error("Expected page to contain post title")
	}
}

func TestPostHandler(t *testing.T) {
	// Create test database
	tempDB := "test_post_handler.db"
	defer os.Remove(tempDB)

	db := store.MustOpen(tempDB)
	defer db.Close()

	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		Environment: "prod",
	}

	server := NewServer(db, cfg)

	// Insert test post
	post := &store.Post{
		Slug:        "test-post",
		Title:       "Test Post",
		Summary:     "A test post",
		HTML:        "<h1>Test Content</h1><p>This is a test.</p>",
		RawMD:       "# Test Content\n\nThis is a test.",
		PublishedAt: "2025-09-12T10:00:00Z",
		UpdatedAt:   "2025-09-12T10:00:00Z",
		Draft:       false,
	}

	err := db.UpsertPost(post)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Test existing post
	req := httptest.NewRequest("GET", "/p/test-post", nil)
	w := httptest.NewRecorder()

	server.PostHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !contains(body, "Test Post") {
		t.Error("Expected page to contain post title")
	}

	if !contains(body, "Test Content") {
		t.Error("Expected page to contain post content")
	}

	// Test non-existent post
	req = httptest.NewRequest("GET", "/p/nonexistent", nil)
	w = httptest.NewRecorder()

	server.PostHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	// Test draft post in production (should 404)
	draftPost := &store.Post{
		Slug:        "draft-post",
		Title:       "Draft Post",
		Summary:     "A draft post",
		HTML:        "<p>Draft content</p>",
		RawMD:       "Draft content",
		PublishedAt: "2025-09-12T10:00:00Z",
		UpdatedAt:   "2025-09-12T10:00:00Z",
		Draft:       true,
	}

	err = db.UpsertPost(draftPost)
	if err != nil {
		t.Fatalf("Failed to insert draft post: %v", err)
	}

	req = httptest.NewRequest("GET", "/p/draft-post", nil)
	w = httptest.NewRecorder()

	server.PostHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for draft in prod, got %d", w.Code)
	}
}

func TestTagHandler(t *testing.T) {
	// Create test database
	tempDB := "test_tag_handler.db"
	defer os.Remove(tempDB)

	db := store.MustOpen(tempDB)
	defer db.Close()

	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		Environment: "test",
	}

	server := NewServer(db, cfg)

	// Test tag page
	req := httptest.NewRequest("GET", "/tag/golang", nil)
	w := httptest.NewRecorder()

	server.TagHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !contains(body, "golang") {
		t.Error("Expected page to contain tag name")
	}

	if !contains(body, "Phase 03") {
		t.Error("Expected placeholder message for Phase 03")
	}

	// Test invalid tag path
	req = httptest.NewRequest("GET", "/tag/", nil)
	w = httptest.NewRecorder()

	server.TagHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for empty tag, got %d", w.Code)
	}
}

func TestStaticHandler(t *testing.T) {
	db := store.MustOpen(":memory:")
	defer db.Close()

	cfg := &config.Config{SiteTitle: "Test"}
	server := NewServer(db, cfg)

	// Test CSS request
	req := httptest.NewRequest("GET", "/static/app.css", nil)
	w := httptest.NewRecorder()

	server.StaticHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for CSS, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/css" {
		t.Error("Expected CSS content type")
	}

	// Test non-existent static file
	req = httptest.NewRequest("GET", "/static/nonexistent.js", nil)
	w = httptest.NewRecorder()

	server.StaticHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent file, got %d", w.Code)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || len(substr) == 0 || 
		   indexOfSubstring(s, substr) >= 0)
}

func indexOfSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}