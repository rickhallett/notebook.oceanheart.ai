# notebook.oceanheart.ai

**A minimalist Go blog engine built for learning in public**

---

## Overview

A single-binary blog engine built with **Go**, **SQLite**, and modern web standards. Content is **Markdown-first** with YAML front matter, featuring syntax highlighting, SEO optimization, and Atom feeds. Designed for speed, simplicity, and easy deployment.

**Implementation Status:** Phase 04 Complete ✅
- ✅ Phase 00: Core Infrastructure (Go, SQLite/libSQL, config, health checks)
- ✅ Phase 01: Content Processing (Markdown, YAML front matter, syntax highlighting)
- ✅ Phase 02: HTTP Server & Routing (handlers, middleware)
- ✅ Phase 03: File-Based Templates, Static Assets, Admin Reload
- ✅ Phase 04: Feeds & SEO Features (Atom, sitemap, meta tags, external links)

---

## Features

### Core Functionality
* **Markdown posts** with YAML front matter (title, date, tags, summary, draft status)
* **Syntax highlighting** using Chroma with GitHub theme and line numbers
* **Draft system** - drafts only visible in development mode (`ENV=dev`)
* **SQLite + Turso (libSQL)** support with automatic migrations and content caching
* **Single binary deployment** - no external dependencies

### SEO & Feeds (Phase 04 Complete)
* **Atom 1.0 feed** at `/feed.xml` with latest 20 posts
* **XML sitemap** at `/sitemap.xml` with all published content
* **Complete SEO meta tags** - Open Graph, Twitter Cards, canonical URLs
* **Article metadata** - publication dates, modification times
* **External link processing** - automatic `target="_blank"` and `rel="noopener noreferrer"`

### Technical Features
* **Comprehensive middleware chain** - logging, gzip, security headers, caching
* **Goldmark renderer** with GitHub Flavored Markdown, footnotes
* **External link security** - automatic processing for safety
* **File-based templates** with dev-time hot reload
* **Static assets** served from `internal/view/assets` (no-cache in dev)
* **Extensive test coverage** - 39+ test functions across all packages
* **"Psychology twist"** - cognitive skills and bias tags with special styling

---

## Project Architecture

```
/cmd/notebook/main.go          # Application entry point, server setup
/internal/
  ├── config/                  # Configuration management
  │   ├── config.go           # Environment-based config loading
  │   └── config_test.go      # Configuration tests
  ├── http/                   # HTTP layer
  │   ├── handlers.go         # HTTP handlers (home, post, feed, sitemap)
  │   ├── handlers_test.go    # Handler tests
  │   ├── middleware.go       # Middleware chain (logging, gzip, security, caching)
  │   └── middleware_test.go  # Middleware tests
  ├── store/                  # Data persistence
  │   ├── sqlite.go           # SQLite/Turso (libSQL) operations, migrations
  │   └── sqlite_test.go      # Database tests
  ├── content/                # Content processing
  │   ├── loader.go           # Markdown file loading, front matter parsing
  │   ├── render.go           # Markdown → HTML rendering (Goldmark + Chroma)
  │   ├── links.go            # External link security processing
  │   └── *_test.go           # Comprehensive content tests
  ├── view/                   # File-based templates and assets
  │   ├── templates/          # HTML templates (layouts + pages)
  │   └── assets/             # CSS and related assets
  └── feed/                   # SEO & syndication
      ├── atom.go             # Atom 1.0 feed generation
      ├── sitemap.go          # XML sitemap generation
      └── *_test.go           # Feed and sitemap tests
/migrations/
  └── 001_init.sql            # Database schema
/content/                     # Markdown blog posts
go.mod                        # Go module dependencies
```

---

## Database Schema & Backends

Backends:
- Local: SQLite via `github.com/mattn/go-sqlite3` (default)
- Remote: Turso (libSQL) via `github.com/tursodatabase/libsql-client-go` when `DB_URL` and `DB_AUTH_TOKEN` are set

Schema (applied automatically via programmatic migration):

```sql
-- Posts table with content caching
CREATE TABLE posts (
  id INTEGER PRIMARY KEY,
  slug TEXT UNIQUE NOT NULL,           -- URL slug derived from filename
  title TEXT NOT NULL,                 -- From YAML front matter
  summary TEXT,                        -- Optional description
  html TEXT NOT NULL,                  -- Cached rendered HTML
  raw_md TEXT NOT NULL,                -- Original markdown content
  published_at DATETIME NOT NULL,      -- From front matter date
  updated_at DATETIME NOT NULL,        -- File modification time
  draft BOOLEAN NOT NULL DEFAULT 0     -- Draft status from front matter
);

-- Tags for categorization (ready for future implementation)
CREATE TABLE tags (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

-- Many-to-many post-tag relationships
CREATE TABLE post_tags (
  post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag_id  INTEGER NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);

-- Performance indexes
CREATE INDEX idx_posts_published ON posts(published_at DESC, draft);
CREATE INDEX idx_tags_name ON tags(name);
```

---

## Content Format

### Markdown with YAML Front Matter

Create posts in `/content/` with the filename pattern `YYYY-MM-DD-slug.md`:

````markdown
---
title: "Welcome to Oceanheart Notebook"
date: "2025-09-12"
tags: ["meta", "welcome", "cognitive-skill:reflection"]
summary: "Welcome to my learning-in-public blog."
draft: false
---

# Your Post Title

Your **markdown** content with [links](https://example.com) and code:

```go
func main() {
    fmt.Println("Hello, world!")
}
```
````

### Front Matter Fields
- **title**: Post title (required)
- **date**: Publication date in YYYY-MM-DD format (required)
- **tags**: Array of tags for categorization (optional)
- **summary**: Brief description for SEO and feeds (optional)
- **draft**: Boolean - `true` hides post in production (optional, defaults to `false`)

### Special Features
- **Syntax highlighting**: Powered by Chroma with GitHub theme and line numbers
- **External links**: Automatically processed for security (`target="_blank"`, `rel="noopener noreferrer"`)
- **Psychology tags**: `cognitive-skill:*` and `bias:*` tags render with special styling
- **GitHub Flavored Markdown**: Tables, task lists, strikethrough supported

---

## API Endpoints

### Public Routes
- **`GET /`** - Home page with post listings
- **`GET /p/{slug}`** - Individual post pages
- **`GET /tag/{name}`** - Tag filtering (template in place; data wiring TBD)
- **`GET /static/*`** - Static asset serving from `internal/view/assets`

### SEO & Syndication  
- **`GET /feed.xml`** - Atom 1.0 feed (latest 20 posts)
- **`GET /sitemap.xml`** - XML sitemap with all published content

### System
- **`GET /healthz`** - Health check endpoint returning JSON status
- **`GET /admin/reload`** - Reloads content from `CONTENT_DIR` and upserts to DB
  - Dev: allowed without auth
  - Prod: requires `RELOAD_TOKEN` via `X-Reload-Token` header or `?token=...`

### Response Features
- **Gzip compression** for text-based responses
- **Security headers** (CSP, XSS protection, etc.)
- **Cache headers** (5 minutes for dynamic content, 1 year for static)
- **SEO meta tags** on all pages (Open Graph, Twitter Cards, canonical URLs)

---

## Server Architecture

```go
// cmd/notebook/main.go - Application entry point
func main() {
    // Load configuration from environment
    cfg := config.LoadConfig()
    
    // Initialize SQLite database with automatic migrations
    db := store.MustOpen(cfg.DBPath)
    defer db.Close()
    
    // Load and cache markdown content
    loader := content.NewLoader(cfg.ContentDir)
    posts, err := loader.LoadAll()
    if err != nil {
        log.Printf("Warning: Failed to load content: %v", err)
    }
    
    // Cache posts in database for fast serving
    if len(posts) > 0 {
        db.UpsertPosts(posts)
        log.Printf("Loaded and cached %d posts", len(posts))
    }
    
    // Setup HTTP server with handlers
    server := httpserver.NewServer(db, cfg)
    mux := http.NewServeMux()
    
    // Route definitions
    mux.HandleFunc("/", server.HomeHandler)
    mux.HandleFunc("/p/", server.PostHandler)
    mux.HandleFunc("/tag/", server.TagHandler)
    mux.HandleFunc("/static/", server.StaticHandler)
    mux.HandleFunc("/feed.xml", server.FeedHandler)
    mux.HandleFunc("/sitemap.xml", server.SitemapHandler)
    mux.HandleFunc("/admin/reload", server.AdminReloadHandler)
    mux.HandleFunc("/healthz", healthCheckHandler)
    
    // Apply comprehensive middleware chain
    handler := httpserver.ChainMiddleware(mux,
        httpserver.SecurityHeadersMiddleware,
        httpserver.LoggingMiddleware,
        httpserver.GzipMiddleware,
        httpserver.CacheHeadersMiddleware,
    )
    
    // Start server
    srv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}
    log.Fatal(srv.ListenAndServe())
}
```

---

## Technology Stack

### Core Dependencies
```go
// go.mod - Key dependencies
require (
    github.com/mattn/go-sqlite3 v1.14.32                          // SQLite driver
    github.com/tursodatabase/libsql-client-go v0.0.x              // Turso/libSQL client
    github.com/yuin/goldmark v1.7.13                               // Markdown parser
    github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729...   // Syntax highlighting
    github.com/alecthomas/chroma/v2 v2.20.0                       // Code highlighter
    gopkg.in/yaml.v3 v3.0.1                                      // YAML front matter
    github.com/dlclark/regexp2 v1.11.5                           // Regex support for Chroma
)
```

### Architecture Decisions
- **Go standard library** for HTTP server (no external web frameworks)
- **SQLite/Turso (libSQL)** for persistence and caching (Turso optional; falls back to SQLite)
- **Goldmark** for extensible, CommonMark-compliant markdown processing
- **Chroma** for syntax highlighting with GitHub theme
- **File-based templates** with dev-time reparse for instant feedback
- **Comprehensive middleware** for production-ready HTTP handling

---

## Setup & Installation

### Prerequisites
- **Go 1.22.7** or later
- **CGO enabled** for SQLite support

### Local Development

```bash
# Clone and build
go mod tidy
go build -o notebook ./cmd/notebook

# Run in development mode
ENV=dev ./notebook

# Or run directly with go
ENV=dev go run ./cmd/notebook
```

### Environment Configuration

```bash
# Default values (all optional)
ENV=prod                                    # "dev" shows drafts
PORT=8080                                   # Server port
DB_PATH=./notebook.db                       # SQLite database file
CONTENT_DIR=./content                       # Markdown files directory
SITE_BASEURL=https://notebook.oceanheart.ai # Used in feeds/sitemaps
SITE_TITLE="Oceanheart Notebook"            # Site title in feeds/meta
RELOAD_TOKEN=                               # Optional: protects /admin/reload in prod

# Optional: Turso (libSQL) remote database
DB_URL=                                     # e.g. libsql://<db-name>-<org>.turso.io
DB_AUTH_TOKEN=                              # Turso auth token
```

### Production Build

```bash
# Build optimized binary
CGO_ENABLED=1 go build -ldflags "-s -w" -o notebook ./cmd/notebook

# Run in production
./notebook

# Container image (Dockerfile provided)
docker build -t notebook:latest .
docker run -p 8080:8080 --env ENV=prod notebook:latest

# Fly.io
fly launch   # once
fly deploy
```

### Quick Start

1. **Add content**: Create markdown files in `/content/` with YAML front matter
2. **Start server**: `go run ./cmd/notebook` 
3. **View blog**: Open http://localhost:8080
4. **Check feeds**: Visit http://localhost:8080/feed.xml and http://localhost:8080/sitemap.xml
5. (Optional) **Dev hot reload**: `scripts/dev.sh` (requires `watchexec`, `reflex`, or `entr`)

---

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific package tests  
go test -v ./internal/content
```

**Test Coverage**: 39+ test functions across all packages
- Configuration loading and environment handling
- Content loading, parsing, and rendering
- HTTP handlers and middleware
- Feed and sitemap generation
- External link processing
- Database operations

---

## Content Processing Pipeline

1. **File Loading**: Scan `/content/` for `*.md` files
2. **Front Matter Parsing**: Extract YAML metadata using `gopkg.in/yaml.v3`
3. **Markdown Rendering**: Convert to HTML using Goldmark with:
   - GitHub Flavored Markdown extensions
   - Syntax highlighting via Chroma
   - External link security processing
4. **Database Caching**: Store rendered HTML and metadata in SQLite
5. **HTTP Serving**: Serve cached content with appropriate headers

---

## Performance & Caching

- **SQLite caching**: Rendered HTML stored in database for fast serving
- **Gzip compression**: Automatic compression for text responses  
- **HTTP caching**: Appropriate cache headers (5 min dynamic, 1 year static)
- **Single binary**: No external dependencies, fast startup
- **Minimal memory footprint**: Efficient Go standard library usage

---

## Security Features

- **Security headers**: CSP, XSS protection, frame options, referrer policy
- **External link processing**: Automatic `target="_blank"` and `rel="noopener noreferrer"`
- **Draft protection**: Drafts only visible in development mode
- **Input sanitization**: Safe HTML rendering through Goldmark
- **SQLite safety**: Parameterized queries, no SQL injection risks

---

## Roadmap

### Next
- Tag filtering (wire DB relations to templates)
- Search functionality
- Theming refinements

### Future Phases
- HTMX progressive enhancement
- Admin interface
- Comment system
- Image handling
- Performance optimizations

---

## Development

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Update dependencies  
go mod tidy

# View module graph
go mod graph
```

---

## License

MIT © Oceanheart.ai
