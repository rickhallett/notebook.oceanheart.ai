package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	loader := NewLoader("./test")
	
	content := `---
title: "Test Post"
date: "2025-09-12"
tags: ["go", "testing", "cognitive-skill:analysis"]
summary: "A test post for validation"
draft: false
---

# Test Content

This is a **test** markdown post with some ` + "`code`" + ` highlighting.

` + "```go\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```"

	post, err := loader.ParseContent(content, "test-post.md")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got %s", post.Title)
	}

	if post.Slug != "test-post" {
		t.Errorf("Expected slug 'test-post', got %s", post.Slug)
	}

	if post.Summary != "A test post for validation" {
		t.Errorf("Expected summary, got %s", post.Summary)
	}

	if post.Draft != false {
		t.Errorf("Expected draft=false, got %v", post.Draft)
	}

	if len(post.HTML) == 0 {
		t.Error("Expected HTML content to be generated")
	}

	if len(post.RawMD) == 0 {
		t.Error("Expected RawMD to be preserved")
	}
}

func TestGenerateSlug(t *testing.T) {
	loader := NewLoader("./test")

	tests := []struct {
		filePath string
		expected string
	}{
		{"2025-09-12-my-post.md", "my-post"},
		{"simple-post.md", "simple-post"},
		{"./content/2025-01-01-new-year.md", "new-year"},
		{"no-date-post.md", "no-date-post"},
	}

	for _, test := range tests {
		result := loader.generateSlug(test.filePath)
		if result != test.expected {
			t.Errorf("generateSlug(%s) = %s, expected %s", test.filePath, result, test.expected)
		}
	}
}

func TestParseDate(t *testing.T) {
	loader := NewLoader("./test")

	tests := []struct {
		input    string
		hasError bool
	}{
		{"2025-09-12", false},
		{"2025-09-12T10:30:00Z", false},
		{"", false}, // Should default to now
		{"invalid-date", true},
	}

	for _, test := range tests {
		result, err := loader.parseDate(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input %s, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			}
			if len(result) == 0 {
				t.Errorf("Expected date result for input %s", test.input)
			}
		}
	}
}

func TestLoadAllFromDirectory(t *testing.T) {
	// Create temporary test directory
	tempDir := "test_content"
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test markdown file
	testContent := `---
title: "Integration Test"
date: "2025-09-12"
tags: ["test"]
summary: "Test file"
draft: false
---

# Test

This is a test.
`

	err = os.WriteFile(filepath.Join(tempDir, "test-post.md"), []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test loading
	loader := NewLoader(tempDir)
	posts, err := loader.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}

	if posts[0].Title != "Integration Test" {
		t.Errorf("Expected title 'Integration Test', got %s", posts[0].Title)
	}
}