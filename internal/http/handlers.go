package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/store"
)

// Server holds the HTTP server dependencies
type Server struct {
	store *store.Store
	cfg   *config.Config
	tmpl  *template.Template
}

// NewServer creates a new HTTP server
func NewServer(store *store.Store, cfg *config.Config) *Server {
	s := &Server{
		store: store,
		cfg:   cfg,
	}
	s.loadTemplates()
	return s
}

// loadTemplates compiles all HTML templates
func (s *Server) loadTemplates() {
	// For Phase 02, we'll use simple inline templates
	// In Phase 03+, this will load from files
	s.tmpl = template.Must(template.New("base").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - {{.SiteTitle}}</title>
    <meta name="description" content="{{.Description}}">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; line-height: 1.6; }
        h1, h2, h3 { color: #333; }
        .header { border-bottom: 1px solid #eee; margin-bottom: 2rem; padding-bottom: 1rem; }
        .nav { margin-bottom: 2rem; }
        .nav a { margin-right: 1rem; text-decoration: none; color: #0066cc; }
        .post-list { list-style: none; padding: 0; }
        .post-item { margin-bottom: 1.5rem; padding-bottom: 1rem; border-bottom: 1px solid #f0f0f0; }
        .post-title { margin: 0 0 0.5rem 0; }
        .post-meta { color: #666; font-size: 0.9rem; }
        .tag { background: #f0f0f0; padding: 0.2rem 0.5rem; border-radius: 3px; font-size: 0.8rem; margin-right: 0.5rem; }
        .tag.cognitive-skill { background: #e6f3ff; color: #0066cc; }
        .tag.bias { background: #ffe6e6; color: #cc0000; }
        pre { background: #f8f8f8; padding: 1rem; overflow-x: auto; }
        code { background: #f0f0f0; padding: 0.2rem 0.4rem; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="header">
        <h1><a href="/" style="text-decoration: none; color: inherit;">{{.SiteTitle}}</a></h1>
        <div class="nav">
            <a href="/">Home</a>
            <a href="/healthz">Health</a>
        </div>
    </div>
    <main>{{.Content}}</main>
</body>
</html>
`))
}

// HomeHandler serves the home page with post listings
func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Only serve root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := s.store.GetAllPosts(s.cfg.IsDev())
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		log.Printf("Error loading posts: %v", err)
		return
	}

	// Build post list HTML
	var content strings.Builder
	content.WriteString("<ul class=\"post-list\">")
	
	if len(posts) == 0 {
		content.WriteString("<li>No posts found. Create some content in the /content directory!</li>")
	}

	for _, post := range posts {
		content.WriteString(fmt.Sprintf(`
		<li class="post-item">
			<h2 class="post-title"><a href="/p/%s">%s</a></h2>
			<div class="post-meta">%s</div>
			<p>%s</p>
		</li>`, post.Slug, post.Title, post.PublishedAt, post.Summary))
	}
	content.WriteString("</ul>")

	data := map[string]interface{}{
		"Title":       "Home",
		"SiteTitle":   s.cfg.SiteTitle,
		"Description": "Learning in public blog",
		"Content":     template.HTML(content.String()),
	}

	if err := s.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

// PostHandler serves individual post pages
func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path /p/{slug}
	slug := strings.TrimPrefix(r.URL.Path, "/p/")
	if slug == "" || slug == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	post, err := s.store.GetPostBySlug(slug)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error getting post %s: %v", slug, err)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	// Don't serve drafts in production
	if post.Draft && !s.cfg.IsDev() {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title":       post.Title,
		"SiteTitle":   s.cfg.SiteTitle,
		"Description": post.Summary,
		"Content":     template.HTML(post.HTML),
	}

	if err := s.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

// TagHandler serves tag filtering pages
func (s *Server) TagHandler(w http.ResponseWriter, r *http.Request) {
	// Extract tag from URL path /tag/{name}
	tagName := strings.TrimPrefix(r.URL.Path, "/tag/")
	if tagName == "" || tagName == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// For Phase 02, we'll implement a simple version
	// Phase 03+ will add proper tag filtering from database
	content := fmt.Sprintf(`
		<h2>Posts tagged: %s</h2>
		<p>Tag filtering will be implemented in Phase 03.</p>
		<p><a href="/">← Back to home</a></p>
	`, tagName)

	data := map[string]interface{}{
		"Title":       fmt.Sprintf("Tag: %s", tagName),
		"SiteTitle":   s.cfg.SiteTitle,
		"Description": fmt.Sprintf("Posts tagged with %s", tagName),
		"Content":     template.HTML(content),
	}

	if err := s.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

// StaticHandler serves static assets
func (s *Server) StaticHandler(w http.ResponseWriter, r *http.Request) {
	// For Phase 02, serve minimal static content
	// Phase 03+ will serve from filesystem
	if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("/* Static CSS will be added in Phase 03 */"))
		return
	}

	http.NotFound(w, r)
}