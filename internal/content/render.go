package content

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
)

// Renderer handles markdown to HTML conversion with syntax highlighting
type Renderer struct {
	md goldmark.Markdown
}

// NewRenderer creates a new markdown renderer with syntax highlighting
func NewRenderer() *Renderer {
	// Configure syntax highlighting
	highlighter := highlighting.NewHighlighting(
		highlighting.WithStyle("github"),
		highlighting.WithFormatOptions(
			html.WithLineNumbers(true),
			html.WithClasses(true),
		),
	)

	// Create goldmark instance with extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown
			extension.Footnote,   // Footnote support
			highlighter,          // Syntax highlighting
		),
		goldmark.WithRendererOptions(
			htmlrenderer.WithHardWraps(),
			htmlrenderer.WithXHTML(),
			htmlrenderer.WithUnsafe(), // Allow raw HTML
		),
	)

	return &Renderer{md: md}
}

// Render converts markdown to HTML
func (r *Renderer) Render(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetStyle returns CSS for syntax highlighting
func (r *Renderer) GetStyle() (string, error) {
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true))
	var buf bytes.Buffer
	if err := formatter.WriteCSS(&buf, style); err != nil {
		return "", err
	}
	return buf.String(), nil
}