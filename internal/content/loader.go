package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"notebook.oceanheart.ai/internal/store"
)

// FrontMatter represents the YAML front matter of a markdown file
type FrontMatter struct {
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Summary string   `yaml:"summary"`
	Draft   bool     `yaml:"draft"`
}

// Loader handles loading and parsing markdown content from filesystem
type Loader struct {
	contentDir string
	renderer   *Renderer
}

// NewLoader creates a new content loader
func NewLoader(contentDir string) *Loader {
	return &Loader{
		contentDir: contentDir,
		renderer:   NewRenderer(),
	}
}

// LoadAll loads all markdown files from the content directory
func (l *Loader) LoadAll() ([]*store.Post, error) {
	var posts []*store.Post

	err := filepath.Walk(l.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		post, err := l.LoadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", path, err)
		}

		if post != nil {
			posts = append(posts, post)
		}

		return nil
	})

	return posts, err
}

// LoadFile loads and parses a single markdown file
func (l *Loader) LoadFile(filePath string) (*store.Post, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return l.ParseContent(string(content), filePath)
}

// ParseContent parses markdown content with front matter
func (l *Loader) ParseContent(content, filePath string) (*store.Post, error) {
	// Split front matter and content
	frontMatter, markdown, err := l.splitFrontMatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %w", err)
	}

	// Generate slug from filename
	slug := l.generateSlug(filePath)

	// Render markdown to HTML
	html, err := l.renderer.Render(markdown)
	if err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}

	// Parse date
	publishedAt, err := l.parseDate(frontMatter.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	post := &store.Post{
		Slug:        slug,
		Title:       frontMatter.Title,
		Summary:     frontMatter.Summary,
		HTML:        html,
		RawMD:       markdown,
		PublishedAt: publishedAt,
		UpdatedAt:   time.Now().Format(time.RFC3339),
		Draft:       frontMatter.Draft,
	}

	return post, nil
}

// splitFrontMatter separates YAML front matter from markdown content
func (l *Loader) splitFrontMatter(content string) (*FrontMatter, string, error) {
	// Check for front matter delimiter
	if !strings.HasPrefix(content, "---\n") {
		return &FrontMatter{}, content, nil
	}

	// Find the closing delimiter
	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("invalid front matter format")
	}

	// Parse YAML front matter
	var fm FrontMatter
	if err := yaml.Unmarshal([]byte(parts[0]), &fm); err != nil {
		return nil, "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &fm, parts[1], nil
}

// generateSlug creates a URL-friendly slug from filepath
func (l *Loader) generateSlug(filePath string) string {
	// Get filename without extension
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	slug := strings.TrimSuffix(base, ext)

	// Remove date prefix if present (e.g., "2025-09-12-")
	parts := strings.Split(slug, "-")
	if len(parts) >= 4 && len(parts[0]) == 4 && len(parts[1]) == 2 && len(parts[2]) == 2 {
		// Looks like YYYY-MM-DD prefix, remove it
		slug = strings.Join(parts[3:], "-")
	}

	return slug
}

// parseDate parses date string into RFC3339 format
func (l *Loader) parseDate(dateStr string) (string, error) {
	if dateStr == "" {
		return time.Now().Format(time.RFC3339), nil
	}

	// Try common date formats
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("unable to parse date: %s", dateStr)
}