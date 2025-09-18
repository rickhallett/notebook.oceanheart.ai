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

	if !contains(body, "Posts tagged:") {
		t.Error("Expected posts tagged heading")
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

func TestSEOMetaTags(t *testing.T) {
	// Create test database
	tempDB := "test_seo.db"
	defer os.Remove(tempDB)

	db := store.MustOpen(tempDB)
	defer db.Close()

	cfg := &config.Config{
		SiteTitle:   "SEO Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "test",
	}

	server := NewServer(db, cfg)

	// Add a test post
	post := &store.Post{
		Slug:        "seo-test-post",
		Title:       "SEO Test Post",
		Summary:     "This is a test post for SEO meta tags",
		HTML:        "<p>Test content</p>",
		RawMD:       "Test content",
		PublishedAt: "2025-09-12T10:00:00Z",
		UpdatedAt:   "2025-09-12T10:00:00Z",
		Draft:       false,
	}

	if err := db.UpsertPost(post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Test home page SEO meta tags
	t.Run("HomePage SEO", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		server.HomeHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()

		// Check basic meta tags
		if !contains(body, `<meta name="viewport" content="width=device-width, initial-scale=1.0">`) {
			t.Error("Expected viewport meta tag")
		}

		if !contains(body, `<meta name="description" content="Learning in public">`) {
			t.Error("Expected description meta tag")
		}

		// Check Open Graph tags
		if !contains(body, `<meta property="og:title" content="Home - SEO Test Blog">`) {
			t.Error("Expected Open Graph title")
		}

		if !contains(body, `<meta property="og:type" content="website">`) {
			t.Error("Expected Open Graph type")
		}

		if !contains(body, `<meta property="og:url" content="https://example.com/">`) {
			t.Error("Expected Open Graph URL")
		}

		// Check canonical URL
		if !contains(body, `<link rel="canonical" href="https://example.com/">`) {
			t.Error("Expected canonical URL")
		}
	})

	// Test post page SEO meta tags
	t.Run("PostPage SEO", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/p/seo-test-post", nil)
		w := httptest.NewRecorder()

		server.PostHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()

		// Check basic meta tags
		if !contains(body, `<meta name="description" content="This is a test post for SEO meta tags">`) {
			t.Error("Expected description meta tag with post summary")
		}

		// Check Open Graph tags
		if !contains(body, `<meta property="og:title" content="SEO Test Post - SEO Test Blog">`) {
			t.Error("Expected Open Graph title with post title")
		}

		if !contains(body, `<meta property="og:type" content="article">`) {
			t.Error("Expected Open Graph type article for posts")
		}

		if !contains(body, `<meta property="og:url" content="https://example.com/p/seo-test-post">`) {
			t.Error("Expected Open Graph URL with post URL")
		}

		if !contains(body, `<meta property="og:description" content="This is a test post for SEO meta tags">`) {
			t.Error("Expected Open Graph description with post summary")
		}

		// Check article-specific meta tags
		if !contains(body, `<meta name="article:published_time" content="2025-09-12T10:00:00Z">`) {
			t.Error("Expected article published time")
		}

		// Check canonical URL
		if !contains(body, `<link rel="canonical" href="https://example.com/p/seo-test-post">`) {
			t.Error("Expected canonical URL for post")
		}

		// Check Twitter Card tags
		if !contains(body, `<meta name="twitter:card" content="summary">`) {
			t.Error("Expected Twitter card type")
		}

		if !contains(body, `<meta name="twitter:title" content="SEO Test Post - SEO Test Blog">`) {
			t.Error("Expected Twitter title")
		}

		if !contains(body, `<meta name="twitter:description" content="This is a test post for SEO meta tags">`) {
			t.Error("Expected Twitter description")
		}
	})

	// Note: Tag page SEO testing skipped - tags functionality not yet implemented in store
}

func TestFeedHandler(t *testing.T) {
	// Create test database
	tempDB := "test_feed.db"
	defer os.Remove(tempDB)

	db := store.MustOpen(tempDB)
	defer db.Close()

	cfg := &config.Config{
		SiteTitle:   "Feed Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "test",
	}

	server := NewServer(db, cfg)

	// Add a test post
	post := &store.Post{
		Slug:        "feed-test-post",
		Title:       "Feed Test Post",
		Summary:     "This is a test post for feed generation",
		HTML:        "<p>Test content</p>",
		RawMD:       "Test content",
		PublishedAt: "2025-09-12T10:00:00Z",
		UpdatedAt:   "2025-09-12T10:00:00Z",
		Draft:       false,
	}

	if err := db.UpsertPost(post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Test feed generation
	req := httptest.NewRequest("GET", "/feed.xml", nil)
	w := httptest.NewRecorder()

	server.FeedHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if !contains(contentType, "application/atom+xml") {
		t.Errorf("Expected Atom XML content type, got: %s", contentType)
	}

	body := w.Body.String()

	// Check feed structure
	if !contains(body, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML declaration")
	}

	if !contains(body, `<feed xmlns="http://www.w3.org/2005/Atom">`) {
		t.Error("Expected Atom feed root element")
	}

	if !contains(body, `<title>Feed Test Blog</title>`) {
		t.Error("Expected feed title")
	}

	if !contains(body, `<entry>`) {
		t.Error("Expected feed entry")
	}
}

func TestSitemapHandler(t *testing.T) {
	// Create test database
	tempDB := "test_sitemap.db"
	defer os.Remove(tempDB)

	db := store.MustOpen(tempDB)
	defer db.Close()

	cfg := &config.Config{
		SiteTitle:   "Sitemap Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "test",
	}

	server := NewServer(db, cfg)

	// Add a test post
	post := &store.Post{
		Slug:        "sitemap-test-post",
		Title:       "Sitemap Test Post",
		Summary:     "This is a test post for sitemap generation",
		HTML:        "<p>Test content</p>",
		RawMD:       "Test content",
		PublishedAt: "2025-09-12T10:00:00Z",
		UpdatedAt:   "2025-09-12T10:00:00Z",
		Draft:       false,
	}

	if err := db.UpsertPost(post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Test sitemap generation
	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	server.SitemapHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if !contains(contentType, "application/xml") {
		t.Errorf("Expected XML content type, got: %s", contentType)
	}

	body := w.Body.String()

	// Check sitemap structure
	if !contains(body, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML declaration")
	}

	if !contains(body, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`) {
		t.Error("Expected sitemap urlset root element")
	}

	if !contains(body, `<loc>https://example.com/</loc>`) {
		t.Error("Expected home page URL in sitemap")
	}

	if !contains(body, `<loc>https://example.com/p/sitemap-test-post</loc>`) {
		t.Error("Expected post URL in sitemap")
	}

	if !contains(body, `<loc>https://example.com/feed.xml</loc>`) {
		t.Error("Expected feed URL in sitemap")
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
