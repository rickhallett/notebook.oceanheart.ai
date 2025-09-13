## Blog Workflow Guide

### Local Setup
- Prereqs: Go 1.22.x (CGO enabled for SQLite).
- Start dev server:
  ```bash
  ENV=dev PORT=3010 \
  DB_PATH=./notebook.dev.db \
  CONTENT_DIR=./content \
  SITE_BASEURL=http://notebook.lvh.me:3010 \
  go run ./cmd/notebook
  ```
- Visit: `http://notebook.lvh.me:3010` (drafts visible in dev).

### Create a Post
- Location: `content/`
- Filename pattern: `YYYY-MM-DD-slug.md` (slug derives from filename after date).
- Front matter (YAML):
  ```yaml
  ---
  title: "Post Title"       # required
  date: "2025-09-12"        # YYYY-MM-DD
  tags: ["tag1", "tag2"]   # optional
  summary: "Short SEO blurb" # optional
  draft: false               # optional (default false)
  ---
  ```

### Edit a Post
- Edit the markdown body and/or front matter.
- Keep the filename base the same to preserve the slug/URL.
- `updated_at` is derived from file modification time; `date` controls publication date.

### Preview & Reload
- Dev shows drafts and should reflect changes quickly.
- To reload content without restarting the server (dev):
  ```bash
  curl http://notebook.lvh.me:3010/admin/reload
  ```
- Verify at `/` and `/p/<slug>`; check code blocks, links, and summaries.

### Drafts
- `draft: true` hides posts in production.
- Drafts are visible when running with `ENV=dev`.

### Publish to Production
- One-command publish from repo root (rsync + restart):
  ```bash
  export NOTEBOOK_HOST=<server>
  # optional: NOTEBOOK_USER, NOTEBOOK_REMOTE_DIR, NOTEBOOK_SERVICE
  ./scripts/publish.sh
  ```
- Or reload content remotely without restart (set `RELOAD_TOKEN` in `/etc/notebook.env` on the server first):
  ```bash
  curl -H "X-Reload-Token: $RELOAD_TOKEN" https://notebook.oceanheart.ai/admin/reload
  ```

### Deployment Notes
- Production behind Caddy at `https://notebook.oceanheart.ai`.
- Server env (`/etc/notebook.env`) example:
  ```bash
  ENV=prod
  PORT=8080
  DB_PATH=/var/lib/notebook/prod/notebook.db
  CONTENT_DIR=/opt/notebook/content
  SITE_BASEURL=https://notebook.oceanheart.ai
  SITE_TITLE=Oceanheart Notebook
  # optional for remote reload:
  # RELOAD_TOKEN=changeme-strong-token
  ```

### Feeds & SEO
- Atom feed: `GET /feed.xml` (latest posts).
- Sitemap: `GET /sitemap.xml`.
- `summary` becomes the meta description; canonical URLs and Open Graph/Twitter tags are set automatically.
