package content

import (
	"strings"
	"testing"
)

func TestMarkdownRendering(t *testing.T) {
	renderer := NewRenderer()

	markdown := `# Test Heading

This is a **bold** text and *italic* text.

Here's some inline ` + "`code`" + ` and a code block:

` + "```go\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```" + `

And a [link](https://example.com).
`

	html, err := renderer.Render(markdown)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check for basic HTML elements
	expected := []string{
		"<h1>Test Heading</h1>",
		"<strong>bold</strong>",
		"<em>italic</em>",
		"<code>code</code>",
		"<pre class=\"chroma\">", // Code blocks use chroma pre tags
		"<a href=\"https://example.com\">link</a>",
	}

	for _, exp := range expected {
		if !strings.Contains(html, exp) {
			t.Errorf("Expected HTML to contain %s, but it didn't. HTML was: %s", exp, html)
		}
	}
}

func TestCodeHighlighting(t *testing.T) {
	renderer := NewRenderer()

	markdown := "```go\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```"

	html, err := renderer.Render(markdown)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check for syntax highlighting classes
	if !strings.Contains(html, "chroma") {
		t.Error("Expected syntax highlighting classes in HTML")
	}

	// Should contain line numbers
	if !strings.Contains(html, "line") {
		t.Error("Expected line numbers in code blocks")
	}
}

func TestGitHubFlavoredMarkdown(t *testing.T) {
	renderer := NewRenderer()

	// Test strikethrough
	markdown := "~~strikethrough~~ text"
	html, err := renderer.Render(markdown)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "<del>strikethrough</del>") {
		t.Error("Expected GFM strikethrough to work")
	}
}

func TestGetStyle(t *testing.T) {
	renderer := NewRenderer()

	css, err := renderer.GetStyle()
	if err != nil {
		t.Fatalf("GetStyle failed: %v", err)
	}

	if len(css) == 0 {
		t.Error("Expected CSS content for syntax highlighting")
	}

	// Should contain chroma-related CSS
	if !strings.Contains(css, "chroma") {
		t.Error("Expected chroma-related CSS classes")
	}
}