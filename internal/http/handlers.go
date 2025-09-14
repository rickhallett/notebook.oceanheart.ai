package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/content"
	"notebook.oceanheart.ai/internal/feed"
	"notebook.oceanheart.ai/internal/store"
	"notebook.oceanheart.ai/internal/view"
)

// Server holds the HTTP server dependencies
type Server struct {
	store *store.Store
	cfg   *config.Config
	view  *view.Manager
}

// NewServer creates a new HTTP server
func NewServer(store *store.Store, cfg *config.Config) *Server {
	s := &Server{store: store, cfg: cfg}
	s.view = view.NewManager("internal/view/templates", cfg.IsDev())
	return s
}

// loadTemplates compiles all HTML templates
func (s *Server) loadTemplates() {}

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

	data := map[string]interface{}{
		"Title":        "Home",
		"SiteTitle":    s.cfg.SiteTitle,
		"Description":  "Learning in public",
		"CanonicalURL": s.cfg.SiteBaseURL + "/",
		"BaseURL":      s.cfg.SiteBaseURL,
		"IsPost":       false,
		"Posts":        posts,
		"HasPosts":     len(posts) > 0,
	}

	// Set Content-Type header for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if contentHTML, err := s.view.RenderString("pages/home.content", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	} else {
		data["Content"] = template.HTML(contentHTML)
	}
	if err := s.view.Execute(w, "base", data); err != nil {
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
		"Title":        post.Title,
		"SiteTitle":    s.cfg.SiteTitle,
		"Description":  post.Summary,
		"CanonicalURL": s.cfg.SiteBaseURL + "/p/" + post.Slug,
		"BaseURL":      s.cfg.SiteBaseURL,
		"IsPost":       true,
		"PublishedAt":  post.PublishedAt,
		"UpdatedAt":    post.UpdatedAt,
		"Post":         post,
	}

	// Set Content-Type header for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if contentHTML, err := s.view.RenderString("pages/post.content", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	} else {
		data["Content"] = template.HTML(contentHTML)
	}
	if err := s.view.Execute(w, "base", data); err != nil {
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

	data := map[string]interface{}{
		"Title":        fmt.Sprintf("Tag: %s", tagName),
		"SiteTitle":    s.cfg.SiteTitle,
		"Description":  fmt.Sprintf("Posts tagged with %s", tagName),
		"CanonicalURL": s.cfg.SiteBaseURL + "/tag/" + tagName,
		"BaseURL":      s.cfg.SiteBaseURL,
		"IsPost":       false,
		"Tag":          tagName,
	}

	// Set Content-Type header for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if contentHTML, err := s.view.RenderString("pages/tag.content", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	} else {
		data["Content"] = template.HTML(contentHTML)
	}
	if err := s.view.Execute(w, "base", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

// StaticHandler serves static assets
func (s *Server) StaticHandler(w http.ResponseWriter, r *http.Request) {
	// Serve from internal/view/assets when file exists; otherwise allow CSS placeholder
	const base = "internal/view/assets"
	rel := strings.TrimPrefix(r.URL.Path, "/static/")
	clean := filepath.Clean(rel)
	// Prevent directory traversal
	if strings.Contains(clean, "..") {
		http.NotFound(w, r)
		return
	}
	fp := filepath.Join(base, clean)
	if st, err := os.Stat(fp); err == nil && !st.IsDir() {
		if s.cfg.IsDev() {
			w.Header().Set("Cache-Control", "no-store, max-age=0")
		}
		http.ServeFile(w, r, fp)
		return
	}
	if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("/* static asset not found; placeholder css */"))
		return
	}
	http.NotFound(w, r)
}

// FeedHandler serves the Atom feed
func (s *Server) FeedHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.store.GetAllPosts(s.cfg.IsDev())
	if err != nil {
		http.Error(w, "Failed to load posts for feed", http.StatusInternalServerError)
		log.Printf("Error loading posts for feed: %v", err)
		return
	}

	// Convert posts to pointers
	var postPtrs []*store.Post
	for i := range posts {
		postPtrs = append(postPtrs, &posts[i])
	}

	// Generate Atom feed
	atomXML, err := feed.GenerateAtom(postPtrs, s.cfg)
	if err != nil {
		http.Error(w, "Failed to generate feed", http.StatusInternalServerError)
		log.Printf("Error generating feed: %v", err)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	w.WriteHeader(http.StatusOK)
	w.Write(atomXML)
}

// SitemapHandler serves the XML sitemap
func (s *Server) SitemapHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.store.GetAllPosts(s.cfg.IsDev())
	if err != nil {
		http.Error(w, "Failed to load posts for sitemap", http.StatusInternalServerError)
		log.Printf("Error loading posts for sitemap: %v", err)
		return
	}

	// Convert posts to pointers
	var postPtrs []*store.Post
	for i := range posts {
		postPtrs = append(postPtrs, &posts[i])
	}

	// Generate sitemap
	sitemapXML, err := feed.GenerateSitemap(postPtrs, s.cfg)
	if err != nil {
		http.Error(w, "Failed to generate sitemap", http.StatusInternalServerError)
		log.Printf("Error generating sitemap: %v", err)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	w.WriteHeader(http.StatusOK)
	w.Write(sitemapXML)
}

// ChromaCSSHandler serves the CSS for syntax highlighting
func (s *Server) ChromaCSSHandler(w http.ResponseWriter, r *http.Request) {
	renderer := content.NewRenderer()
	css, err := renderer.GetStyle()
	if err != nil {
		http.Error(w, "Failed to generate CSS", http.StatusInternalServerError)
		log.Printf("Error generating Chroma CSS: %v", err)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	if s.cfg.IsDev() {
		w.Header().Set("Cache-Control", "no-store, max-age=0")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(css))
}

// AdminReloadHandler reloads content from the filesystem and upserts into the DB.
// Dev: always allowed. Prod: requires RELOAD_TOKEN via header or query param.
func (s *Server) AdminReloadHandler(w http.ResponseWriter, r *http.Request) {
	// Authorization
	token := r.Header.Get("X-Reload-Token")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if !s.cfg.AllowReload(token) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Load from content dir
	loader := content.NewLoader(s.cfg.ContentDir)
	posts, err := loader.LoadAll()
	if err != nil {
		log.Printf("reload: failed to load content: %v", err)
		http.Error(w, "failed to load content", http.StatusInternalServerError)
		return
	}

	// Upsert
	if err := s.store.UpsertPosts(posts); err != nil {
		log.Printf("reload: failed to upsert posts: %v", err)
		http.Error(w, "failed to upsert posts", http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"ok","reloaded":%d}`, len(posts))))
}
