## Local Dev Hot Reloader (PRD)

### Executive Summary
Improve developer experience with a fast, predictable hot‑reload loop. In dev, templates already reparse on each request; this PRD adds file watching to reload Markdown content automatically and restarts the server on Go code changes. No production impact or runtime dependencies.

### Problem Statement
Current workflow requires manual restarts or hitting the reload endpoint to pick up content changes; Go code changes require manual restart. We want sub‑second feedback for templates/CSS and near‑instant reloads for content and code with minimal tooling.

### Requirements
- Templates: live on each request (already in place via view.Manager in dev).
- Assets: serve CSS with no-cache in dev so changes reflect on refresh.
- Content: auto‑call `/admin/reload` when files under `content/**/*.md` change.
- Code: auto‑restart `go run ./cmd/notebook` on `**/*.go` changes.
- DX: single command to run watchers + app; zero config.
- Safety: only active when `ENV=dev`; no effect in prod builds.

### Implementation Phases
- Phase 1: Watch content and trigger reload
  - Use a lightweight watcher (e.g., `watchexec`/`reflex`/`entr`) to `curl http://localhost:$PORT/admin/reload` on `.md` changes.
- Phase 2: Watch Go sources and restart app
  - Use the same tool to restart `go run ./cmd/notebook` on `.go` changes.
- Phase 3: One‑shot DX wrapper
  - Add `scripts/dev.sh` or a `make dev` target that runs both watchers concurrently.
- Phase 4 (optional): In‑app fsnotify
  - Integrate `fsnotify` to watch `content/` and call the in‑process reload logic. Keep behind `ENV=dev`.

### Implementation Notes
- Example with watchexec (recommended):
  ```bash
  # Terminal 1: restart on Go changes
  watchexec -r -e go -- ENV=dev PORT=3010 \
    DB_PATH=./notebook.dev.db CONTENT_DIR=./content \
    SITE_BASEURL=http://notebook.lvh.me:8003 \
    go run ./cmd/notebook

  # Terminal 2: reload on Markdown changes
  watchexec -w content -e md -- \
    curl -fsS http://localhost:3010/admin/reload >/dev/null || true
  ```
- Example with reflex:
  ```bash
  reflex -r '\\.(go)$' -- sh -c 'ENV=dev PORT=3010 DB_PATH=./notebook.dev.db CONTENT_DIR=./content SITE_BASEURL=http://notebook.lvh.me:8003 go run ./cmd/notebook'
  reflex -r '\\.(md)$'  -- sh -c 'curl -fsS http://localhost:3010/admin/reload >/dev/null || true'
  ```
- CSS & Templates:
  - Templates are already reparsed per request in dev; no extra work.
  - Ensure `Cache-Control: no-store` for `/static/*` in dev to avoid stale CSS (small middleware tweak if needed).

### Security Considerations
- `/admin/reload` requires no token in dev; in prod, ensure `RELOAD_TOKEN` is set and required (already enforced).
- Watchers run only on developer machines; do not add them to production images.

### Success Metrics
- Template/CSS changes reflect on refresh with no manual steps.
- Markdown edits visible within 1s after file save.
- Go code changes restart and serve within ~2–3s.
- One command starts the full dev loop reliably.

### Future Enhancements
- Unified `scripts/dev.sh` that spawns both watchers and handles app lifecycle.
- File‑level selective reload to avoid reloading all posts.
- Browser live‑reload (websocket snippet) for auto refresh after rebuild.
