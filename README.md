# notebook.oceanheart.ai

**A minimalist Go + HTMX + SQLite blog for learning in public**

---

## 1) Overview

Single-binary blog engine built with **Go (net/http)**, **HTMX** (progressive enhancement), and **SQLite**.
Content is **Markdown-first**, rendered to HTML at request time and cached. Mobile-first layout, fast, and easy to extend.

---

## 2) Features

* **Markdown posts** with front matter (title, date, tags, summary).
* **HTMX** search/filter and lazy pagination (no SPA).
* **Tags & archives**, RSS/Atom feed, sitemap.xml.
* **Drafts** (only visible when `ENV=dev`).
* **Fast SQLite cache** (+ auto-migrate).
* **Code highlighting** (Chroma), smart typography.
* **Single binary** deploy; optional Docker.
* **“Psychology twist”** metadata: cognitive skills/bias tags rendered as ribbons.

---

## 3) Project structure

```
/cmd/notebook/main.go          # wire server, routes, flags
/internal/http/handlers.go     # handlers (home, post, feed, search, admin)
/internal/http/middleware.go   # logging, gzip, cache headers
/internal/store/sqlite.go      # SQLite access (posts, tags)
/internal/content/loader.go    # load/parse markdown, front matter
/internal/content/render.go    # markdown → HTML (goldmark + chroma)
/internal/view/templates/      # base.tmpl, post.tmpl, list.tmpl, partials
/internal/view/assets/         # css, minimal js, favicon
/migrations/*.sql              # schema
/content/                      # your markdown posts
Makefile
Dockerfile
README.md
```

---

## 4) Data model (SQLite)

```sql
-- migrations/001_init.sql
CREATE TABLE posts (
  id INTEGER PRIMARY KEY,
  slug TEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  summary TEXT,
  html TEXT NOT NULL,
  raw_md TEXT NOT NULL,
  published_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  draft BOOLEAN NOT NULL DEFAULT 0
);
CREATE TABLE tags (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);
CREATE TABLE post_tags (
  post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag_id  INTEGER NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);
CREATE INDEX idx_posts_published ON posts(published_at DESC, draft);
CREATE INDEX idx_tags_name ON tags(name);
```

---

## 5) Content format (Markdown + front matter)

Example: `content/2025-09-12-minimalism.md`

````md
---
title: "Minimalism in Go"
date: "2025-09-12"
tags: ["go", "architecture", "cognitive-skill:abstraction"]
summary: "Notes on stripping a web app to essentials."
draft: false
---

Your markdown here. ```go code blocks ``` will be highlighted.
````

**Notes**

* `cognitive-skill:*` or `bias:*` tags render as colored ribbons for the psych twist.

---

## 6) Handlers (routes)

* `GET /` – list page (HTMX paginated).
* `GET /p/:slug` – post page.
* `GET /tag/:name` – tag listing.
* `GET /search?q=...` – HTMX partial results.
* `GET /feed.xml` – Atom/RSS.
* `GET /sitemap.xml` – sitemap.
* `GET /healthz` – health check.
* **Admin-lite (optional)**

  * `POST /admin/reindex` – reload content from `/content`.
  * `POST /admin/flush-cache` – clear HTML cache.

---

## 7) Minimal server wiring (excerpt)

```go
// cmd/notebook/main.go
func main() {
  dbPath := env("DB_PATH", "./notebook.db")
  contentDir := env("CONTENT_DIR", "./content")
  envMode := env("ENV", "prod")

  store := store.MustOpen(dbPath)          // apply migrations
  loader := content.NewLoader(contentDir)  // reads *.md
  posts := loader.LoadAll()                // parse & render
  store.UpsertPosts(posts)                 // cache into SQLite

  mux := http.NewServeMux()
  mux.Handle("/", htmxCache(homeHandler(store)))
  mux.Handle("/p/", htmxCache(postHandler(store)))
  mux.Handle("/search", searchHandler(store))
  mux.Handle("/feed.xml", feedHandler(store))
  // admin
  mux.Handle("/admin/reindex", adminOnly(reindexHandler(store, loader, envMode)))

  srv := &http.Server{Addr: ":8080", Handler: gzip(mwLog(mux))}
  log.Fatal(srv.ListenAndServe())
}
```

---

## 8) HTMX search (partial)

```html
<form hx-get="/search" hx-target="#results" hx-push-url="true">
  <input name="q" placeholder="Search posts…" />
</form>
<div id="results" hx-trigger="load" hx-get="/search"></div>
```

Server returns the `partials/list_items.tmpl` snippet for HTMX swaps.

---

## 9) Build & run

### Local

```bash
# dev prerequisites: Go 1.22+
make dev        # go run ./cmd/notebook
# or
go run ./cmd/notebook
```

Default env:

```
ENV=dev
DB_PATH=./notebook.db
CONTENT_DIR=./content
SITE_BASEURL=https://notebook.oceanheart.ai
SITE_TITLE="Oceanheart Notebook"
```

### Production build

```bash
make build      # produces ./bin/notebook
./bin/notebook
```

### Docker

```Dockerfile
# Dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o notebook ./cmd/notebook

FROM gcr.io/distroless/static
WORKDIR /
COPY --from=build /app/notebook /notebook
COPY content /content
ENV DB_PATH=/data/notebook.db CONTENT_DIR=/content
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/notebook"]
```

```bash
docker build -t notebook .
docker run -p 8080:8080 -v $(pwd)/data:/data notebook
```

---

## 10) Theming & typography

* Edit `/internal/view/assets/app.css` (Tailwind-free, \~3KB CSS).
* Templates: `base.tmpl`, `post.tmpl`, `list.tmpl`.
* Add logo/colors to match Oceanheart.

---

## 11) Code highlighting & markdown

* Renderer: `goldmark` + `goldmark-highlighting` (Chroma).
* Smart punctuation: `goldmark-smartypants`.
* External links: `target=_blank` with `rel="noopener"`.

---

## 12) Caching strategy

* SQLite stores `html` and `raw_md`. On startup or **reindex**, if file mtime is newer than DB, re-render.
* Send `ETag`/`Last-Modified`; simple in-memory LRU for hot posts.

---

## 13) SEO & feeds

* `feed.xml` (Atom) with latest 20 posts.
* `sitemap.xml` generated from DB slugs.
* `<meta>` tags per post (title, description from `summary`, published time).

---

## 14) Security & ops

* `ENV=prod` hides drafts and disables admin endpoints by default.
* Simple rate limit on search (IP bucket).
* `/healthz` returns 200 + build hash.
* Logs: structured (level, method, path, ms).

---

## 15) Makefile (excerpt)

```makefile
dev:
\tENV=dev go run ./cmd/notebook
build:
\tCGO_ENABLED=0 go build -ldflags "-s -w" -o bin/notebook ./cmd/notebook
test:
\tgo test ./...
fmt:
\tgofmt -s -w .
reindex:
\tcurl -X POST localhost:8080/admin/reindex
```

---

## 16) Roadmap (tiny, focused)

* [ ] Pagination via HTMX “Load more”.
* [ ] Image pipeline (thumbs, alt text helper).
* [ ] Simple in-page footnotes.
* [ ] Per-post “cognitive ribbons” legend page.
* [ ] Optional PWA manifest + offline shell.
* [ ] Minimal comments (static file or PB-backed).

---

## 17) License

MIT © Oceanheart.ai

---

## 18) Deploy note

Point `notebook.oceanheart.ai` → your host (Fly/Railway/Render).
Binary reads posts from `/content`; push new posts via git deploy or rsync, then hit `/admin/reindex` (or restart).
