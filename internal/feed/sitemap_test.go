package feed

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/store"
)

func TestGenerateSitemap(t *testing.T) {
	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "prod",
	}

	posts := []*store.Post{
		{
			Slug:        "test-post-1",
			Title:       "Test Post 1",
			Summary:     "First test post",
			HTML:        "<p>Test content</p>",
			RawMD:       "Test content",
			PublishedAt: "2025-09-12T10:00:00Z",
			UpdatedAt:   "2025-09-12T10:00:00Z",
			Draft:       false,
		},
		{
			Slug:        "test-post-2",
			Title:       "Test Post 2",
			Summary:     "Second test post",
			HTML:        "<p>Test content</p>",
			RawMD:       "Test content",
			PublishedAt: "2025-09-11T10:00:00Z",
			UpdatedAt:   "2025-09-11T10:00:00Z",
			Draft:       false,
		},
		{
			Slug:        "draft-post",
			Title:       "Draft Post",
			Summary:     "Draft content",
			HTML:        "<p>Draft content</p>",
			RawMD:       "Draft content",
			PublishedAt: "2025-09-10T10:00:00Z",
			UpdatedAt:   "2025-09-10T10:00:00Z",
			Draft:       true,
		},
	}

	sitemapXML, err := GenerateSitemap(posts, cfg)
	if err != nil {
		t.Fatalf("GenerateSitemap failed: %v", err)
	}

	// Parse the generated XML to validate structure
	var sitemap SitemapIndex
	err = xml.Unmarshal(sitemapXML, &sitemap)
	if err != nil {
		t.Fatalf("Failed to parse generated sitemap XML: %v", err)
	}

	// Verify sitemap metadata
	if sitemap.Xmlns != "http://www.sitemaps.org/schemas/sitemap/0.9" {
		t.Errorf("Expected sitemap namespace, got %s", sitemap.Xmlns)
	}

	// Should have:
	// - Home page (1)
	// - Published posts (2, excluding draft)
	// - Feed URL (1)
	// Total: 4 URLs
	expectedURLs := 4
	if len(sitemap.URLs) != expectedURLs {
		t.Errorf("Expected %d URLs, got %d", expectedURLs, len(sitemap.URLs))
	}

	// Verify home page URL
	homeURL := sitemap.URLs[0]
	if homeURL.Loc != "https://example.com/" {
		t.Errorf("Expected home URL 'https://example.com/', got %s", homeURL.Loc)
	}
	if homeURL.Priority != "1.0" {
		t.Errorf("Expected home priority '1.0', got %s", homeURL.Priority)
	}
	if homeURL.ChangeFreq != "daily" {
		t.Errorf("Expected home changefreq 'daily', got %s", homeURL.ChangeFreq)
	}

	// Verify post URLs
	postURLs := 0
	for _, url := range sitemap.URLs {
		if strings.Contains(url.Loc, "/p/") {
			postURLs++
			// Verify post URL properties
			if url.Priority != "0.8" {
				t.Errorf("Expected post priority '0.8', got %s", url.Priority)
			}
			if url.ChangeFreq != "monthly" {
				t.Errorf("Expected post changefreq 'monthly', got %s", url.ChangeFreq)
			}
			if url.LastMod == "" {
				t.Error("Expected post to have lastmod")
			}
		}
	}

	if postURLs != 2 {
		t.Errorf("Expected 2 post URLs, got %d", postURLs)
	}

	// Verify feed URL
	feedURLFound := false
	for _, url := range sitemap.URLs {
		if url.Loc == "https://example.com/feed.xml" {
			feedURLFound = true
			if url.Priority != "0.5" {
				t.Errorf("Expected feed priority '0.5', got %s", url.Priority)
			}
			break
		}
	}

	if !feedURLFound {
		t.Error("Expected feed URL in sitemap")
	}

	// Verify XML structure
	xmlString := string(sitemapXML)
	if !strings.Contains(xmlString, "<?xml") {
		t.Error("Expected XML declaration")
	}

	if !strings.Contains(xmlString, "<urlset") {
		t.Error("Expected urlset element")
	}
}

func TestGenerateSitemapWithDrafts(t *testing.T) {
	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "dev", // Dev mode should include drafts
	}

	posts := []*store.Post{
		{
			Slug:        "published-post",
			Title:       "Published Post",
			Summary:     "Published content",
			HTML:        "<p>Published content</p>",
			RawMD:       "Published content",
			PublishedAt: "2025-09-12T10:00:00Z",
			UpdatedAt:   "2025-09-12T10:00:00Z",
			Draft:       false,
		},
		{
			Slug:        "draft-post",
			Title:       "Draft Post",
			Summary:     "Draft content",
			HTML:        "<p>Draft content</p>",
			RawMD:       "Draft content",
			PublishedAt: "2025-09-11T10:00:00Z",
			UpdatedAt:   "2025-09-11T10:00:00Z",
			Draft:       true,
		},
	}

	sitemapXML, err := GenerateSitemap(posts, cfg)
	if err != nil {
		t.Fatalf("GenerateSitemap failed: %v", err)
	}

	var sitemap SitemapIndex
	err = xml.Unmarshal(sitemapXML, &sitemap)
	if err != nil {
		t.Fatalf("Failed to parse generated sitemap XML: %v", err)
	}

	// In dev mode, should include drafts
	// Home (1) + Published post (1) + Draft post (1) + Feed (1) = 4 URLs
	expectedURLs := 4
	if len(sitemap.URLs) != expectedURLs {
		t.Errorf("Expected %d URLs in dev mode, got %d", expectedURLs, len(sitemap.URLs))
	}
}

func TestGenerateSitemapWithNoPosts(t *testing.T) {
	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "prod",
	}

	var posts []*store.Post

	sitemapXML, err := GenerateSitemap(posts, cfg)
	if err != nil {
		t.Fatalf("GenerateSitemap failed: %v", err)
	}

	var sitemap SitemapIndex
	err = xml.Unmarshal(sitemapXML, &sitemap)
	if err != nil {
		t.Fatalf("Failed to parse generated sitemap XML: %v", err)
	}

	// Should still have home page and feed URL
	expectedURLs := 2
	if len(sitemap.URLs) != expectedURLs {
		t.Errorf("Expected %d URLs with no posts, got %d", expectedURLs, len(sitemap.URLs))
	}
}

func TestFormatSitemapDate(t *testing.T) {
	// Test with specific time
	testTime := "2025-09-12T10:30:45Z"
	parsed, err := time.Parse(time.RFC3339, testTime)
	if err != nil {
		t.Fatalf("Failed to parse test time: %v", err)
	}

	result := formatSitemapDate(parsed)
	
	// Should be in W3C datetime format
	if !strings.Contains(result, "2025-09-12T10:30:45") {
		t.Errorf("Expected formatted date to contain '2025-09-12T10:30:45', got %s", result)
	}
	
	// Should have timezone info
	if !strings.Contains(result, "+") && !strings.Contains(result, "-") && !strings.Contains(result, "Z") {
		t.Errorf("Expected timezone info in formatted date, got %s", result)
	}
}