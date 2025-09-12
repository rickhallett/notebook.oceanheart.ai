package content

import (
	"strings"
	"testing"
)

func TestProcessExternalLinks(t *testing.T) {
	baseURL := "https://example.com"
	
	testCases := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "external link without attributes",
			html:     `<a href="https://external.com/page">External Link</a>`,
			expected: `<a href="https://external.com/page" target="_blank" rel="noopener noreferrer">External Link</a>`,
		},
		{
			name:     "internal link unchanged",
			html:     `<a href="https://example.com/page">Internal Link</a>`,
			expected: `<a href="https://example.com/page">Internal Link</a>`,
		},
		{
			name:     "relative link unchanged",
			html:     `<a href="/page">Relative Link</a>`,
			expected: `<a href="/page">Relative Link</a>`,
		},
		{
			name:     "anchor link unchanged",
			html:     `<a href="#section">Anchor Link</a>`,
			expected: `<a href="#section">Anchor Link</a>`,
		},
		{
			name:     "mailto link unchanged",
			html:     `<a href="mailto:test@example.com">Email Link</a>`,
			expected: `<a href="mailto:test@example.com">Email Link</a>`,
		},
		{
			name:     "external link with existing target",
			html:     `<a href="https://external.com/page" target="_self">External Link</a>`,
			expected: `<a href="https://external.com/page" target="_self" rel="noopener noreferrer">External Link</a>`,
		},
		{
			name:     "external link with existing rel",
			html:     `<a href="https://external.com/page" rel="bookmark">External Link</a>`,
			expected: `<a href="https://external.com/page" rel="bookmark" target="_blank">External Link</a>`,
		},
		{
			name:     "external link with both attributes",
			html:     `<a href="https://external.com/page" target="_blank" rel="noopener">External Link</a>`,
			expected: `<a href="https://external.com/page" target="_blank" rel="noopener">External Link</a>`,
		},
		{
			name:     "multiple links mixed",
			html:     `<a href="https://external1.com">Ext1</a> <a href="/internal">Int</a> <a href="https://external2.com">Ext2</a>`,
			expected: `<a href="https://external1.com" target="_blank" rel="noopener noreferrer">Ext1</a> <a href="/internal">Int</a> <a href="https://external2.com" target="_blank" rel="noopener noreferrer">Ext2</a>`,
		},
		{
			name:     "link with additional attributes",
			html:     `<a href="https://external.com/page" class="btn" id="link1">External Link</a>`,
			expected: `<a href="https://external.com/page" class="btn" id="link1" target="_blank" rel="noopener noreferrer">External Link</a>`,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ProcessExternalLinks(tc.html, baseURL)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestIsExternalLink(t *testing.T) {
	baseURL := "https://example.com"
	
	testCases := []struct {
		url      string
		expected bool
		name     string
	}{
		// Internal links
		{"https://example.com/page", false, "internal link same domain"},
		{"http://example.com/page", false, "internal link same domain http"},
		{"https://example.com/path/to/page", false, "internal link with path"},
		{"/relative/path", false, "relative path"},
		{"#anchor", false, "anchor link"},
		{"//cdn.example.com", false, "protocol-relative"},
		
		// External links
		{"https://external.com", true, "external https"},
		{"http://external.com", true, "external http"},
		{"https://subdomain.external.com", true, "external subdomain"},
		
		// Special protocols
		{"mailto:test@example.com", false, "mailto protocol"},
		{"tel:+1234567890", false, "tel protocol"},
		{"javascript:alert('test')", false, "javascript protocol"},
		
		// Edge cases
		{"relative-file.html", false, "relative file"},
		{"", false, "empty string"},
		{"https://", true, "incomplete url"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isExternalLink(tc.url, baseURL)
			if result != tc.expected {
				t.Errorf("isExternalLink(%q, %q) = %v, expected %v", tc.url, baseURL, result, tc.expected)
			}
		})
	}
}

func TestProcessExternalLinksWithDifferentBaseDomains(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
		html    string
		expectExternal bool
	}{
		{
			name:    "subdomain should be external",
			baseURL: "https://example.com",
			html:    `<a href="https://blog.example.com/post">Blog Link</a>`,
			expectExternal: true,
		},
		{
			name:    "www variant should be external",
			baseURL: "https://example.com",
			html:    `<a href="https://www.example.com/page">WWW Link</a>`,
			expectExternal: true,
		},
		{
			name:    "same domain different protocol",
			baseURL: "https://example.com",
			html:    `<a href="http://example.com/page">HTTP Link</a>`,
			expectExternal: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ProcessExternalLinks(tc.html, tc.baseURL)
			hasTargetBlank := strings.Contains(result, `target="_blank"`)
			hasRelAttributes := strings.Contains(result, `rel="noopener noreferrer"`)
			
			if tc.expectExternal {
				if !hasTargetBlank || !hasRelAttributes {
					t.Errorf("Expected external link processing, got: %s", result)
				}
			} else {
				if hasTargetBlank || hasRelAttributes {
					t.Errorf("Expected no external link processing, got: %s", result)
				}
			}
		})
	}
}

func TestProcessExternalLinksEmptyInput(t *testing.T) {
	result := ProcessExternalLinks("", "https://example.com")
	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
	
	result = ProcessExternalLinks("No links here", "https://example.com")
	if result != "No links here" {
		t.Errorf("Expected unchanged text, got %q", result)
	}
}

func TestProcessExternalLinksMalformedHTML(t *testing.T) {
	testCases := []struct {
		name string
		html string
	}{
		{
			name: "unclosed link tag",
			html: `<a href="https://external.com">Link without closing tag`,
		},
		{
			name: "malformed href",
			html: `<a href=https://external.com>Link without quotes</a>`,
		},
		{
			name: "multiple href attributes",
			html: `<a href="https://external.com" href="https://other.com">Multiple href</a>`,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic and should return some result
			result := ProcessExternalLinks(tc.html, "https://example.com")
			if len(result) == 0 {
				t.Errorf("Expected some output for malformed HTML, got empty string")
			}
		})
	}
}