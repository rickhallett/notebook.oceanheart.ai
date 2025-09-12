package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/feed"
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
    {{if .PublishedAt}}<meta name="article:published_time" content="{{.PublishedAt}}">{{end}}
    {{if .UpdatedAt}}<meta name="article:modified_time" content="{{.UpdatedAt}}">{{end}}
    <meta property="og:title" content="{{.Title}} - {{.SiteTitle}}">
    <meta property="og:description" content="{{.Description}}">
    <meta property="og:type" content="{{if .IsPost}}article{{else}}website{{end}}">
    <meta property="og:url" content="{{.CanonicalURL}}">
    <meta name="twitter:card" content="summary">
    <meta name="twitter:title" content="{{.Title}} - {{.SiteTitle}}">
    <meta name="twitter:description" content="{{.Description}}">
    <link rel="canonical" href="{{.CanonicalURL}}">
    <link rel="alternate" type="application/atom+xml" title="{{.SiteTitle}} Feed" href="{{.BaseURL}}/feed.xml">
    <style>
        /* Hugo Mini Theme CSS - Base Typography & Colors */
        html[theme='dark-mode'] {
          filter: invert(1) hue-rotate(180deg);
        }

        body {
          line-height: 1;
          font: normal 15px/1.5em 'Helvetica Neue', Helvetica, Arial, sans-serif;
          color: #404040;
          line-height: 1.75;
          letter-spacing: 0.008em;
        }

        a {
          text-decoration: none;
          color: #5badf0;
        }

        a:hover {
          color: #0366d6;
        }

        /* Typography */
        h1, h2, h3 {
          font-weight: 400;
          color: #404040;
        }

        p {
          margin-block-start: 1.5em;
          margin-block-end: 1.5em;
        }

        p, pre {
          word-break: normal;
          overflow-wrap: anywhere;
        }

        /* Code Typography */
        p code {
          font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, Courier, monospace;
          font-size: inherit;
          background-color: rgba(0, 0, 0, 0.06);
          padding: 0 2px;
          border: 1px solid rgba(0, 0, 0, 0.08);
          border-radius: 2px 2px;
          line-height: inherit;
          word-wrap: break-word;
          text-indent: 0;
        }

        pre code {
          font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, Courier, monospace;
        }

        .highlight pre {
          padding: 7px;
          overflow-x: auto;
        }

        .highlight {
          max-width: 100%;
          overflow-x: auto;
        }

        /* Markdown Content Styling */
        blockquote {
          margin-top: 5px;
          margin-bottom: 5px;
          padding-left: 1em;
          margin-left: 0px;
          border-left: 3px solid #eee;
          color: #757575;
        }

        hr {
          display: block;
          border: none;
          height: 2px;
          margin: 40px auto;
          background: #eee;
        }

        table {
          width: 100%;
          margin: 40px 0;
          border-collapse: collapse;
          line-height: 1.5em;
        }

        th, td {
          text-align: left;
          padding-right: 20px;
          vertical-align: top;
        }

        table td, td {
          border-spacing: none;
          border-style: solid;
          padding: 10px 15px;
          border-width: 1px 0 0 0;
        }

        thead th, th {
          text-align: left;
          padding: 10px 15px;
          height: 20px;
          font-size: 13px;
          font-weight: bold;
          color: #444;
          cursor: default;
          white-space: nowrap;
          border: 1px solid #dadadc;
        }

        tr>td {
          border: 1px solid #dadadc;
        }

        tr:nth-child(odd)>td {
          background: #fcfcfc;
        }

        .markdown-image img {
          max-width: 100%;
        }

        .anchor { 
          font-size: 100%; 
          visibility: hidden; 
          color: silver;
        }

        h1:hover a, h2:hover a, h3:hover a, h4:hover a { 
          visibility: visible
        }

        /* Layout System */
        .main {
          width: 100%;
          margin: 0 auto;
        }

        /* Navigation */
        nav.navigation {
          padding: 20px 20px 0;
          background: #fff;
          background: rgba(255, 255, 255, 0.9);
          margin: 0 auto;
          text-align: right;
          z-index: 100;
        }

        nav.navigation a {
          top: 8px;
          right: 6px;
          padding: 8px 12px;
          color: #5badf0;
          font-size: 13px;
          line-height: 1.35;
          border-radius: 3px;
        }

        nav.navigation a:hover {
          color: #0366d6;
        }

        nav.navigation a.button {
          background: #5badf0;
          color: #fff;
          margin-left: 12px;
        }

        /* Profile/Header */
        .profile {
          margin: 60px auto 0 auto;
          text-align: center;
        }

        .profile .avatar {
          display: inline-block;
          width: 80px;
          height: 80px;
          border-radius: 50%;
        }

        .profile h1 {
          font-weight: 400;
          letter-spacing: 0px;
          font-size: 20px;
          color: #404040;
          margin-bottom: 0;
          margin-top: 0;
        }

        .profile h2 {
          font-size: 20px;
          font-weight: 300;
          color: #757575;
          margin-top: 0;
        }

        /* Post List Layout (Home Page) */
        #list-page {
          max-width: 580px;
          margin: 0 auto;
          padding: 0 24px;
        }

        #list-page .item {
          margin: 12px 0;
        }

        #list-page .title {
          display: inline-block;
          color: #404040;
          font-size: 20px;
          font-weight: 400;
          margin: 0;
          width: 80%;
        }

        #list-page .title a {
          color: #404040;
          display: block;
        }

        #list-page .title a:hover {
          color: #0366d6;
        }

        #list-page .date {
          width: 20%;
          float: right;
          text-align: right;
          position: relative;
          top: 1px;
          color: #bbb;
        }

        #list-page .summary {
          color: #757575;
          margin-top: 12px;
          word-break: normal;
          overflow-wrap: anywhere;
          margin-bottom: 36px;
        }

        #list-page .pagination {
          margin: 48px 0;
          width: 100%;
          height: 32px;
          margin-top: 48px;
        }

        #list-page .pagination .pre {
          float: left;
        }

        #list-page .pagination .next {
          float: right;
        }

        /* Single Post Layout */
        #single {
          max-width: 680px;
          margin: 60px auto 0 auto;
          padding: 0 64px;
        }

        #single .title {
          text-align: center;
          font-size: 32px;
          font-weight: 400;
          line-height: 48px;
        }

        #single .tip {
          text-align: center;
          color: #8c8c8c;
          margin-top: 18px;
          font-size: 14px;
        }

        #single .tip .split {
          margin: 0 4px;
        }

        #single .content {
          margin-top: 36px;
        }

        /* Tags */
        #single .tags {
          margin-top: 24px;
        }

        #single .tags a, .tags a {
          background: #f2f2f2;
          padding: 4px 7px;
          color: #757575;
          font-size: 14px;
          margin-right: 3px;
        }

        #single .tags a:hover, .tags a:hover {
          color: #0366d6;
        }

        /* Table of Contents */
        .toc {
          margin: auto;
          background: #f8f8f8;
          border-radius: 0;
          padding: 10px 7px;
          margin-top: 36px;
        }

        .toc details summary {
          cursor: zoom-in;
          margin-inline-start: 14px;
          font-weight: 500;
        }

        .toc details[open] summary {
          cursor: zoom-out;
        }

        .toc #TableOfContents {
          margin-left: 10px;
        }

        .toc ul {
          padding-inline-start: 24px;
        }

        /* Footer */
        #footer {
          margin-top: 100px;
          margin-bottom: 100px;
          text-align: center;
          color: #bbbbbb;
          font-size: 14px;
        }

        #footer .copyright {
          margin: 20px auto;
          font-size: 15px;
        }

        .powerby {
          margin: 20px auto;
          font-size: 13px;
        }

        #footer .split {
          cursor: pointer;
        }

        #footer .split:hover path {
          fill: #ff3356;
          transition: 0.7s ease-out;
          cursor: pointer;
        }

        #social a {
          margin: 0 4px;
        }

        /* Responsive Design */
        @media (max-width: 700px) {
          nav.navigation {
            padding: 20px 10px 0 0;
            background: #fff;
            background: rgba(255, 255, 255, 0.9);
            margin: 0 auto;
            text-align: right;
            z-index: 100;
          }
          
          nav.navigation a {
            top: 8px;
            right: 6px;
            padding: 8px 8px;
            color: #5badf0;
            font-size: 13px;
            line-height: 1.35;
            border-radius: 3px;
          }

          #single {
            padding: 0 18px;
            margin: 20px auto 0 auto;
          }
          
          #single .title {
            font-size: 24px;
            line-height: 32px;
          }
        }

        @media (max-width: 324px) {
          nav.navigation a.button {
            display: none;
          }
        }

        /* Legacy support for existing classes */
        .header {
          border-bottom: 1px solid #eee;
          margin-bottom: 2rem;
          padding-bottom: 1rem;
        }

        .nav {
          margin-bottom: 2rem;
        }

        .nav a {
          margin-right: 1rem;
          text-decoration: none;
          color: #5badf0;
        }

        .post-list {
          list-style: none;
          padding: 0;
        }

        .post-item {
          margin-bottom: 1.5rem;
          padding-bottom: 1rem;
          border-bottom: 1px solid #f0f0f0;
        }

        .post-title {
          margin: 0 0 0.5rem 0;
        }

        .post-meta {
          color: #757575;
          font-size: 0.9rem;
        }

        .tag {
          background: #f2f2f2;
          padding: 4px 7px;
          color: #757575;
          font-size: 14px;
          margin-right: 3px;
        }

        .tag:hover {
          color: #0366d6;
        }
    </style>
</head>
<body>
    <nav class="navigation">
        <a href="/">Home</a>
        <a href="/healthz">Health</a>
        <a href="/feed.xml" class="button">Subscribe</a>
    </nav>

    <main class="main">{{.Content}}</main>
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

	// Build Hugo Mini style content
	var content strings.Builder
	
	// Add profile header if configured
	content.WriteString(`<header class="profile">
		<h1>` + s.cfg.SiteTitle + `</h1>
		<h2>Learning in public blog</h2>
	</header>`)
	
	// Add post list in Hugo Mini style
	content.WriteString(`<div id="list-page">`)

	if len(posts) == 0 {
		content.WriteString(`<section class="item">
			<div class="title">No posts found. Create some content in the /content directory!</div>
		</section>`)
	}

	for _, post := range posts {
		content.WriteString(fmt.Sprintf(`
		<section class="item">
			<div>
				<h1 class="title"><a href="/p/%s">%s</a></h1>
				<div class="date">%s</div>
			</div>
			<div class="summary">%s</div>
		</section>`, post.Slug, post.Title, post.PublishedAt, post.Summary))
	}
	content.WriteString("</div>")

	data := map[string]interface{}{
		"Title":        "Home",
		"SiteTitle":    s.cfg.SiteTitle,
		"Description":  "Learning in public blog",
		"Content":      template.HTML(content.String()),
		"CanonicalURL": s.cfg.SiteBaseURL + "/",
		"BaseURL":      s.cfg.SiteBaseURL,
		"IsPost":       false,
	}

	// Set Content-Type header for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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

	// Wrap content in Hugo Mini single post structure
	var singleContent strings.Builder
	singleContent.WriteString(`<section id="single">`)
	singleContent.WriteString(fmt.Sprintf(`<h1 class="title">%s</h1>`, post.Title))
	singleContent.WriteString(fmt.Sprintf(`<div class="tip">
		%s
		<span class="split">·</span>
		<span>Updated: %s</span>
	</div>`, post.PublishedAt, post.UpdatedAt))
	singleContent.WriteString(`<div class="content">`)
	singleContent.WriteString(post.HTML)
	singleContent.WriteString(`</div>`)
	singleContent.WriteString(`</section>`)

	data := map[string]interface{}{
		"Title":        post.Title,
		"SiteTitle":    s.cfg.SiteTitle,
		"Description":  post.Summary,
		"Content":      template.HTML(singleContent.String()),
		"CanonicalURL": s.cfg.SiteBaseURL + "/p/" + post.Slug,
		"BaseURL":      s.cfg.SiteBaseURL,
		"IsPost":       true,
		"PublishedAt":  post.PublishedAt,
		"UpdatedAt":    post.UpdatedAt,
	}

	// Set Content-Type header for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
		"Title":        fmt.Sprintf("Tag: %s", tagName),
		"SiteTitle":    s.cfg.SiteTitle,
		"Description":  fmt.Sprintf("Posts tagged with %s", tagName),
		"Content":      template.HTML(content),
		"CanonicalURL": s.cfg.SiteBaseURL + "/tag/" + tagName,
		"BaseURL":      s.cfg.SiteBaseURL,
		"IsPost":       false,
	}

	// Set Content-Type header for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
