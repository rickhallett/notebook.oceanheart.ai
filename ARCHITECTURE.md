# Architecture Documentation

**notebook.oceanheart.ai** - A minimalist blog engine with multi-language development support

Version: 0.1.0  
Last Updated: 2025-09-12  
Repository: https://github.com/rickhallett/notebook.oceanheart.ai

## Table of Contents

1. [System Overview](#system-overview)
2. [Technology Stack](#technology-stack)
3. [System Architecture](#system-architecture)
4. [Component Architecture](#component-architecture)
5. [Data Flow](#data-flow)
6. [Database Schema](#database-schema)
7. [HTTP API Reference](#http-api-reference)
8. [Content Processing Pipeline](#content-processing-pipeline)
9. [Multi-Language Development Environment](#multi-language-development-environment)
10. [Deployment Architecture](#deployment-architecture)
11. [Development Guidelines](#development-guidelines)

## System Overview

notebook.oceanheart.ai is a **learning-in-public** blog engine that combines simplicity with powerful content management features. The system is designed around the philosophy of fast, focused blogging with special emphasis on cognitive psychology tracking through specialized tagging.

### Key Design Principles

- **Single Binary Deployment**: Complete application compiles to one executable
- **File-First Content**: Markdown files with YAML front matter as source of truth
- **Database Caching**: SQLite provides performance layer over filesystem content
- **Progressive Enhancement**: HTMX over full SPA complexity
- **Multi-Language Support**: Go backend with optional Bun/TypeScript, Python, Ruby environments

### Core Features

- Markdown-to-HTML processing with syntax highlighting
- Psychology-aware tagging system (`cognitive-skill:*`, `bias:*`)
- Atom feeds and XML sitemaps
- External link security processing
- Development admin endpoints
- Multi-environment configuration

## Technology Stack

### Primary (Go Backend)

```
Go 1.22.7
├── goldmark v1.7.13          # Markdown processing
├── goldmark-highlighting     # Syntax highlighting  
├── chroma v2.20.0           # Code syntax themes
├── sqlite3 v1.14.32         # Database driver
└── yaml.v3                  # YAML front matter parsing
```

### Frontend Enhancement

```
Bun Runtime
├── TypeScript ^5            # Type safety
├── React JSX               # Optional component system
└── HTMX                    # Progressive enhancement (served via CDN)
```

### Additional Language Support

- **Python**: UV package manager, Python 3.11+
- **Ruby**: Ruby 3.0.0, Bundler

## System Architecture

The system follows a **layered architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────────┐
│                HTTP Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────┐ │
│  │ Middleware  │  │  Handlers   │  │Templates │ │
│  │             │  │             │  │          │ │
│  └─────────────┘  └─────────────┘  └──────────┘ │
├─────────────────────────────────────────────────┤
│                Business Logic                    │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────┐ │
│  │  Content    │  │    Feed     │  │  Config  │ │
│  │ Processing  │  │ Generation  │  │          │ │
│  └─────────────┘  └─────────────┘  └──────────┘ │
├─────────────────────────────────────────────────┤
│                Data Layer                        │
│  ┌─────────────┐  ┌─────────────┐               │
│  │   Store     │  │ Filesystem  │               │
│  │  (SQLite)   │  │ (Markdown)  │               │
│  └─────────────┘  └─────────────┘               │
└─────────────────────────────────────────────────┘
```

### Request Processing Flow

1. **HTTP Request** → Middleware Chain → Route Handler
2. **Content Loading** → Filesystem → Parse → Cache in SQLite
3. **Response Generation** → Template Rendering → HTTP Response

## Component Architecture

### Directory Structure

```
notebook.oceanheart.ai/
├── cmd/notebook/           # Application entry point
│   └── main.go            # Server initialization and routing
├── internal/              # Internal packages (cannot be imported externally)
│   ├── config/           # Configuration management
│   │   ├── config.go     # Environment variable loading
│   │   └── config_test.go
│   ├── content/          # Markdown processing and caching
│   │   ├── loader.go     # Filesystem content loading
│   │   ├── render.go     # Markdown to HTML conversion
│   │   ├── links.go      # External link processing
│   │   └── *_test.go     # Unit tests
│   ├── store/            # Data persistence layer
│   │   ├── sqlite.go     # SQLite operations and migrations
│   │   └── sqlite_test.go
│   ├── http/             # HTTP server components
│   │   ├── handlers.go   # Request handlers
│   │   ├── middleware.go # HTTP middleware chain
│   │   └── *_test.go
│   ├── feed/             # Feed generation (Atom, Sitemap)
│   │   ├── atom.go
│   │   ├── sitemap.go
│   │   └── *_test.go
│   └── view/             # Templates and static assets
│       ├── templates/    # HTML templates (future)
│       └── assets/       # CSS/JS assets (future)
├── content/              # Markdown blog posts
│   └── *.md             # Posts with YAML front matter
├── migrations/           # Database schema definitions
│   └── 001_init.sql     # Initial schema
├── docs/                 # Documentation and specifications
├── package.json          # Bun/TypeScript configuration
├── tsconfig.json         # TypeScript compiler options
├── pyproject.toml        # Python/UV configuration
├── Gemfile              # Ruby gems configuration
└── go.mod               # Go module dependencies
```

### Core Components

#### 1. Configuration Management (`internal/config/`)

**Purpose**: Environment-based configuration loading with sensible defaults

```go
type Config struct {
    Environment string  // "dev", "prod"
    DBPath      string  // SQLite database file path
    ContentDir  string  // Markdown files directory
    SiteBaseURL string  // Canonical URL for feeds/sitemaps
    SiteTitle   string  // Site branding
    Port        string  // HTTP server port
}
```

**Key Functions**:
- `LoadConfig()`: Environment variable loading with defaults
- `IsDev()`: Development mode detection
- `IsAdmin()`: Admin endpoint enablement

#### 2. Content Processing (`internal/content/`)

**Purpose**: Markdown file processing pipeline with caching

**loader.go**:
- `FrontMatter`: YAML metadata structure
- `Loader.LoadAll()`: Recursive directory scanning
- `Loader.ParseContent()`: Front matter + markdown parsing
- Slug generation from filenames
- Date parsing and validation

**render.go**:
- Goldmark markdown processor with extensions:
  - GitHub Flavored Markdown
  - Footnote support
  - Syntax highlighting (Chroma)
- `Renderer.Render()`: Markdown → HTML conversion
- `Renderer.GetStyle()`: CSS generation for syntax highlighting

**links.go**:
- `ProcessExternalLinks()`: Security attribute injection
- External link detection and `rel="noopener noreferrer"` addition
- Domain-based link classification

#### 3. Data Persistence (`internal/store/`)

**Purpose**: SQLite-based caching layer with automatic migrations

**Core Types**:
```go
type Post struct {
    ID          int
    Slug        string    // URL identifier
    Title       string    // Display title
    Summary     string    // Meta description
    HTML        string    // Rendered content
    RawMD       string    // Original markdown
    PublishedAt string    // RFC3339 timestamp
    UpdatedAt   string    // RFC3339 timestamp
    Draft       bool      // Visibility flag
}

type Tag struct {
    ID   int
    Name string
}
```

**Key Operations**:
- `MustOpen()`: Database initialization with migrations
- `UpsertPosts()`: Batch content caching
- `GetPostBySlug()`: Individual post retrieval
- `LinkPostTags()`: Many-to-many tag relationships

#### 4. HTTP Layer (`internal/http/`)

**handlers.go**:
- `HomeHandler`: Post listing with pagination support
- `PostHandler`: Individual post serving with draft filtering
- `TagHandler`: Tag-based filtering (Phase 03)
- `FeedHandler`: Atom feed generation
- `SitemapHandler`: XML sitemap generation

**middleware.go**:
- `LoggingMiddleware`: Request timing and status logging
- `GzipMiddleware`: Response compression
- `SecurityHeadersMiddleware`: XSS and clickjacking protection
- `CacheHeadersMiddleware`: Static asset and dynamic content caching
- `ChainMiddleware()`: Middleware composition utility

#### 5. Feed Generation (`internal/feed/`)

**atom.go**:
- Atom 1.0 feed specification compliance
- Post content inclusion with HTML sanitization
- Author and metadata handling

**sitemap.go**:
- XML sitemap generation for SEO
- URL priority and change frequency optimization
- Last modification date tracking

## Data Flow

### Content Loading Flow

```
Filesystem (.md files)
    ↓
YAML Front Matter + Markdown Parsing
    ↓
Goldmark Rendering (HTML + Syntax Highlighting)
    ↓
External Link Processing (Security Attributes)
    ↓
SQLite Cache Storage (Posts + Tags + Relations)
    ↓
HTTP Response (Templates + Cached HTML)
```

### Request Lifecycle

```
HTTP Request
    ↓
Middleware Chain:
├── Security Headers
├── Request Logging  
├── Gzip Compression
└── Cache Headers
    ↓
Route Handler:
├── Database Query (SQLite)
├── Template Rendering
└── Response Generation
    ↓
HTTP Response (HTML/XML/JSON)
```

### Content Update Flow

```
New/Modified .md File
    ↓
File System Monitoring (Manual in current version)
    ↓
Content Loader Re-parsing
    ↓
Database Cache Update
    ↓
Feed/Sitemap Regeneration
```

## Database Schema

### SQLite Schema (v001_init)

```sql
-- Posts table with metadata and caching
CREATE TABLE posts (
  id INTEGER PRIMARY KEY,
  slug TEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  summary TEXT,
  html TEXT NOT NULL,          -- Cached rendered HTML
  raw_md TEXT NOT NULL,        -- Original markdown
  published_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  draft BOOLEAN NOT NULL DEFAULT 0
);

-- Tags table for categorization
CREATE TABLE tags (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

-- Many-to-many relationship between posts and tags
CREATE TABLE post_tags (
  post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag_id  INTEGER NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);

-- Performance indexes
CREATE INDEX idx_posts_published ON posts(published_at DESC, draft);
CREATE INDEX idx_tags_name ON tags(name);

-- Migration tracking
CREATE TABLE schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Data Relationships

```
Posts (1) ←→ (M) Post_Tags (M) ←→ (1) Tags

Post:
- Contains rendered HTML cache
- Stores original markdown for re-processing
- Draft flag controls visibility
- Slug provides URL-friendly identifier

Tag:
- Normalized tag storage
- Special prefixes: "cognitive-skill:", "bias:"
- Automatic creation via GetOrCreateTag()
```

## HTTP API Reference

### Public Endpoints

| Method | Path | Description | Content-Type |
|--------|------|-------------|--------------|
| `GET` | `/` | Home page with post listings | `text/html` |
| `GET` | `/p/{slug}` | Individual post page | `text/html` |
| `GET` | `/tag/{name}` | Tag filtering page | `text/html` |
| `GET` | `/feed.xml` | Atom 1.0 feed | `application/atom+xml` |
| `GET` | `/sitemap.xml` | XML sitemap | `application/xml` |
| `GET` | `/healthz` | Health check | `application/json` |
| `GET` | `/static/{file}` | Static assets | varies |

### Response Headers

**Security Headers** (all responses):
```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
```

**Cache Headers**:
```http
# Static assets
Cache-Control: public, max-age=31536000

# Dynamic content
Cache-Control: public, max-age=300

# Feeds
Cache-Control: public, max-age=3600
```

### Error Handling

- **404**: Invalid routes, missing posts, drafts in production
- **500**: Database errors, template errors, feed generation failures
- **Health Check**: `{"status":"ok","service":"notebook.oceanheart.ai"}`

## Content Processing Pipeline

### Markdown File Format

```yaml
---
title: "Post Title"
date: "2025-09-12"
tags: ["go", "architecture", "cognitive-skill:abstraction"]
summary: "Post summary for SEO and listings"
draft: false
---

# Your markdown content here

Code blocks with syntax highlighting:

```go
func main() {
    fmt.Println("Hello, World!")
}
```
```

### Processing Steps

1. **File Discovery**: `filepath.Walk()` scans content directory
2. **Front Matter Parsing**: YAML extraction and validation
3. **Markdown Rendering**: 
   - Goldmark with GFM, footnotes, syntax highlighting
   - Chroma for code block styling
   - HTML sanitization
4. **Link Processing**: External link security attributes
5. **Slug Generation**: Filename to URL conversion
6. **Database Caching**: Upsert operations with tag relationships

### Special Tag Processing

**Cognitive Skill Tags**:
- Pattern: `cognitive-skill:*`
- Rendered as blue ribbons: `background: #e6f3ff; color: #0066cc`
- Examples: `abstraction`, `analysis`, `synthesis`

**Bias Awareness Tags**:
- Pattern: `bias:*`
- Rendered as red ribbons: `background: #ffe6e6; color: #cc0000`
- Examples: `confirmation`, `dunning-kruger`, `availability`

## Multi-Language Development Environment

### Go (Primary)

**Development Commands**:
```bash
# Run development server
go run ./cmd/notebook

# Build production binary
CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/notebook ./cmd/notebook

# Run tests
go test ./...

# Format code
gofmt -s -w .
```

**Key Dependencies**:
- goldmark: Fast markdown processor
- chroma: Syntax highlighting
- sqlite3: Database driver (CGO required)
- yaml.v3: YAML parsing

### Bun/TypeScript (Frontend)

**Development Commands**:
```bash
# Install dependencies
bun install

# Development server with hot reload
bun --hot ./index.ts

# Build assets
bun build <file.html|file.ts|file.css>

# Run tests
bun test
```

**Configuration** (`tsconfig.json`):
- Target: ESNext with bundler mode
- JSX: React JSX support
- Strict type checking enabled

### Python (UV)

**Environment**: Python 3.11+ managed by UV
```bash
# Run Python scripts
uv run python main.py
```

### Ruby

**Environment**: Ruby 3.0.0 with Bundler
```bash
# Install gems
bundle install
```

## Deployment Architecture

### Single Binary Deployment

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags "-s -w" -o notebook ./cmd/notebook

# Deploy
./notebook
```

### Environment Configuration

```bash
# Production environment
ENV=prod
DB_PATH=./notebook.db
CONTENT_DIR=./content
SITE_BASEURL=https://notebook.oceanheart.ai
SITE_TITLE="Oceanheart Notebook"
PORT=8080
```

### File System Layout

```
production/
├── notebook              # Binary executable
├── notebook.db          # SQLite database
├── content/             # Markdown files
│   ├── 2025-09-12-welcome.md
│   └── ...
└── migrations/          # Schema files (embedded in binary)
```

### Performance Characteristics

- **Cold Start**: ~50ms (SQLite connection + content loading)
- **Response Time**: ~5ms for cached content
- **Memory Usage**: ~10MB base + content cache
- **Disk Usage**: ~5MB binary + content + SQLite database

## Development Guidelines

### Code Organization

**Package Structure**:
- `cmd/`: Application entry points only
- `internal/`: Private packages not importable by external projects
- Packages organized by domain responsibility, not technical layer

**Naming Conventions**:
- Go: Standard Go naming (camelCase, PascalCase)
- Files: Snake case for multi-word files (`snake_case_test.go`)
- Database: Snake case for columns (`published_at`)
- URLs: Kebab case for slugs (`my-post-title`)

### Testing Strategy

**Test Organization**:
- Unit tests: `*_test.go` files alongside source
- Integration tests: `integration_test.go` files
- Test databases: Temporary files, cleaned up in `defer`

**Testing Patterns**:
```go
func TestFeature(t *testing.T) {
    // Setup temporary resources
    tempDB := "test_feature.db"
    defer os.Remove(tempDB)
    
    // Test implementation
    // ...
    
    // Assertions with clear error messages
    if got != want {
        t.Errorf("Expected %v, got %v", want, got)
    }
}
```

### Performance Guidelines

**Database**:
- Use prepared statements for repeated queries
- Implement batch operations for bulk updates
- Index on query patterns, not just foreign keys

**Caching Strategy**:
- SQLite serves as content cache layer
- HTTP cache headers for client-side caching
- Gzip compression for text responses

**Memory Management**:
- Stream large responses instead of buffering
- Close database connections and file handles
- Use buffered I/O for file operations

### Security Considerations

**Content Security**:
- External links get `rel="noopener noreferrer"`
- Input validation on all user-controlled data
- SQL injection prevention via prepared statements

**HTTP Security**:
- Security headers on all responses
- Content type validation
- Request size limits (future enhancement)

### Future Development

**Phase 03 Features** (Planned):
- HTMX-powered search and pagination
- File-based template system
- Tag filtering functionality
- Admin panel for content management
- Real-time content reloading

**Scalability Considerations**:
- Read replicas for SQLite (future)
- CDN integration for static assets
- Background job processing for content updates
- Metrics and monitoring integration

---

## Additional Resources

- **Repository**: https://github.com/rickhallett/notebook.oceanheart.ai
- **Go Documentation**: https://pkg.go.dev/notebook.oceanheart.ai
- **Issue Tracker**: GitHub Issues
- **License**: Check repository for license information

*This architecture document is maintained alongside code changes and updated as the system evolves.*