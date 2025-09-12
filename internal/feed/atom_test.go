package feed

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/store"
)

func TestGenerateAtom(t *testing.T) {
	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "test",
	}

	posts := []*store.Post{
		{
			Slug:        "test-post-1",
			Title:       "Test Post 1",
			Summary:     "First test post",
			HTML:        "<p>This is the first test post.</p>",
			RawMD:       "This is the first test post.",
			PublishedAt: "2025-09-12T10:00:00Z",
			UpdatedAt:   "2025-09-12T10:00:00Z",
			Draft:       false,
		},
		{
			Slug:        "test-post-2",
			Title:       "Test Post 2",
			Summary:     "Second test post",
			HTML:        "<p>This is the second test post.</p>",
			RawMD:       "This is the second test post.",
			PublishedAt: "2025-09-11T10:00:00Z",
			UpdatedAt:   "2025-09-11T10:00:00Z",
			Draft:       false,
		},
		{
			Slug:        "draft-post",
			Title:       "Draft Post",
			Summary:     "This is a draft",
			HTML:        "<p>Draft content.</p>",
			RawMD:       "Draft content.",
			PublishedAt: "2025-09-10T10:00:00Z",
			UpdatedAt:   "2025-09-10T10:00:00Z",
			Draft:       true,
		},
	}

	atomXML, err := GenerateAtom(posts, cfg)
	if err != nil {
		t.Fatalf("GenerateAtom failed: %v", err)
	}

	// Parse the generated XML to validate structure
	var feed AtomFeed
	err = xml.Unmarshal(atomXML, &feed)
	if err != nil {
		t.Fatalf("Failed to parse generated Atom XML: %v", err)
	}

	// Verify feed metadata
	if feed.Title != "Test Blog" {
		t.Errorf("Expected title 'Test Blog', got %s", feed.Title)
	}

	if feed.Xmlns != "http://www.w3.org/2005/Atom" {
		t.Errorf("Expected Atom namespace, got %s", feed.Xmlns)
	}

	if feed.ID != "https://example.com/" {
		t.Errorf("Expected ID 'https://example.com/', got %s", feed.ID)
	}

	// Verify links
	if len(feed.Link) < 2 {
		t.Errorf("Expected at least 2 links, got %d", len(feed.Link))
	}

	// Verify entries (should exclude drafts in test mode)
	if len(feed.Entry) != 2 {
		t.Errorf("Expected 2 entries (excluding draft), got %d", len(feed.Entry))
	}

	// Verify first entry
	firstEntry := feed.Entry[0]
	if firstEntry.Title != "Test Post 1" {
		t.Errorf("Expected first entry title 'Test Post 1', got %s", firstEntry.Title)
	}

	if firstEntry.ID != "https://example.com/p/test-post-1" {
		t.Errorf("Expected first entry ID 'https://example.com/p/test-post-1', got %s", firstEntry.ID)
	}

	if firstEntry.Summary != "First test post" {
		t.Errorf("Expected first entry summary 'First test post', got %s", firstEntry.Summary)
	}

	// Verify XML structure
	xmlString := string(atomXML)
	if !strings.Contains(xmlString, "<?xml") {
		t.Error("Expected XML declaration")
	}

	if !strings.Contains(xmlString, "<feed") {
		t.Error("Expected feed element")
	}
}

func TestGenerateAtomWithManyPosts(t *testing.T) {
	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "prod",
	}

	// Create 25 posts to test the 20-post limit
	var posts []*store.Post
	for i := 0; i < 25; i++ {
		posts = append(posts, &store.Post{
			Slug:        fmt.Sprintf("post-%d", i),
			Title:       fmt.Sprintf("Post %d", i),
			Summary:     fmt.Sprintf("Summary %d", i),
			HTML:        fmt.Sprintf("<p>Content %d</p>", i),
			RawMD:       fmt.Sprintf("Content %d", i),
			PublishedAt: "2025-09-12T10:00:00Z",
			UpdatedAt:   "2025-09-12T10:00:00Z",
			Draft:       false,
		})
	}

	atomXML, err := GenerateAtom(posts, cfg)
	if err != nil {
		t.Fatalf("GenerateAtom failed: %v", err)
	}

	var feed AtomFeed
	err = xml.Unmarshal(atomXML, &feed)
	if err != nil {
		t.Fatalf("Failed to parse generated Atom XML: %v", err)
	}

	// Should be limited to 20 entries
	if len(feed.Entry) != 20 {
		t.Errorf("Expected 20 entries (limit), got %d", len(feed.Entry))
	}
}

func TestGenerateAtomWithNoPosts(t *testing.T) {
	cfg := &config.Config{
		SiteTitle:   "Test Blog",
		SiteBaseURL: "https://example.com",
		Environment: "prod",
	}

	var posts []*store.Post

	_, err := GenerateAtom(posts, cfg)
	if err == nil {
		t.Error("Expected error for empty posts list")
	}

	if !strings.Contains(err.Error(), "no posts available") {
		t.Errorf("Expected 'no posts available' error, got: %s", err.Error())
	}
}

func TestGenerateAtomDraftHandling(t *testing.T) {
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
			HTML:        "<p>Published content.</p>",
			RawMD:       "Published content.",
			PublishedAt: "2025-09-12T10:00:00Z",
			UpdatedAt:   "2025-09-12T10:00:00Z",
			Draft:       false,
		},
		{
			Slug:        "draft-post",
			Title:       "Draft Post",
			Summary:     "Draft content",
			HTML:        "<p>Draft content.</p>",
			RawMD:       "Draft content.",
			PublishedAt: "2025-09-11T10:00:00Z",
			UpdatedAt:   "2025-09-11T10:00:00Z",
			Draft:       true,
		},
	}

	atomXML, err := GenerateAtom(posts, cfg)
	if err != nil {
		t.Fatalf("GenerateAtom failed: %v", err)
	}

	var feed AtomFeed
	err = xml.Unmarshal(atomXML, &feed)
	if err != nil {
		t.Fatalf("Failed to parse generated Atom XML: %v", err)
	}

	// In dev mode, should include both published and draft posts
	if len(feed.Entry) != 2 {
		t.Errorf("Expected 2 entries in dev mode, got %d", len(feed.Entry))
	}
}

func TestFormatAtomDate(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool // whether it should parse successfully
	}{
		{"2025-09-12T10:00:00Z", true},
		{"2025-09-12T10:00:00+00:00", true},
		{"invalid-date", false}, // should fallback to current time
		{"", false},             // should fallback to current time
	}

	for _, tc := range testCases {
		result := formatAtomDate(tc.input)
		
		// Result should always be in RFC3339 format
		if len(result) < 19 { // minimum RFC3339 length
			t.Errorf("formatAtomDate(%s) returned invalid format: %s", tc.input, result)
		}
		
		// Should contain date components
		if !strings.Contains(result, "T") || !strings.Contains(result, ":") {
			t.Errorf("formatAtomDate(%s) doesn't look like RFC3339: %s", tc.input, result)
		}
	}
}