# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**notebook.oceanheart.ai** is a minimalist blog engine with multi-language development support:
- **Primary**: Go-based blog engine with HTMX + SQLite
- **Frontend**: Bun + TypeScript/React support
- **Additional**: Python (UV) and Ruby environments

## Development Commands

### Go (Primary Blog Engine)
```bash
# Development server
go run ./cmd/notebook

# Build binary
CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/notebook ./cmd/notebook

# Test
go test ./...

# Format code
gofmt -s -w .
```

### Bun/TypeScript (Frontend)
```bash
# Install dependencies
bun install

# Run development server with hot reload
bun --hot ./index.ts

# Build frontend assets
bun build <file.html|file.ts|file.css>

# Run tests
bun test
```

### Python (UV)
```bash
# Python environment managed by UV
uv run python main.py
```

### Ruby
```bash
# Ruby 3.0.0 environment
bundle install  # if needed
```

## Architecture

### Go Blog Engine Structure
- `cmd/notebook/main.go` - Server entry point, routing, configuration
- `internal/http/` - HTTP handlers and middleware
- `internal/store/` - SQLite data access layer
- `internal/content/` - Markdown processing and rendering
- `internal/view/` - Templates and static assets
- `migrations/*.sql` - Database schema migrations
- `content/` - Markdown blog posts with front matter

### Data Flow
1. Markdown files in `/content` with YAML front matter
2. Content loader parses and renders to HTML (goldmark + chroma highlighting)
3. SQLite caches rendered content with metadata
4. HTMX provides progressive enhancement for search/pagination
5. Templates render server-side HTML with partial updates

### Key Technologies
- **Backend**: Go net/http, SQLite, goldmark (markdown), chroma (syntax highlighting)
- **Frontend**: HTMX (no SPA), minimal CSS (~3KB), progressive enhancement
- **Content**: Markdown + YAML front matter, tag system with psychology twist (cognitive-skill:*, bias:*)
- **Deployment**: Single binary, Docker support, static asset serving

## Environment Variables

```bash
ENV=dev                                    # Show drafts, enable admin endpoints
DB_PATH=./notebook.db                      # SQLite database location
CONTENT_DIR=./content                      # Markdown files directory
SITE_BASEURL=https://notebook.oceanheart.ai # Base URL for feeds/sitemaps
SITE_TITLE="Oceanheart Notebook"          # Site title
```

## Key Routes & Features

- `GET /` - Home page with HTMX pagination
- `GET /p/:slug` - Individual post pages
- `GET /tag/:name` - Tag filtering
- `GET /search?q=...` - HTMX search with partial results
- `GET /feed.xml` - Atom/RSS feed
- `POST /admin/reindex` - Reload content (dev mode)
- `POST /admin/flush-cache` - Clear HTML cache (dev mode)

## Content Format

Blog posts use markdown with YAML front matter:
```markdown
---
title: "Post Title"
date: "2025-09-12"
tags: ["go", "architecture", "cognitive-skill:abstraction"]
summary: "Post summary for SEO and listings"
draft: false
---

Your markdown content here with ```go code blocks``` for highlighting.
```

Special tag prefixes:
- `cognitive-skill:*` - Rendered as colored ribbons
- `bias:*` - Psychology-themed tags with visual treatment

## Development Notes

- Project uses Bun instead of Node.js for TypeScript/frontend work
- Go code follows standard project layout with internal packages
- SQLite provides both storage and caching layer
- HTMX eliminates need for complex frontend framework
- Single binary deployment with embedded assets
- Content is file-system based but cached in database for performance