# Production Deploy (Fly.io) + Turso Migration Plan

## Overview
- Vendor: Fly.io for app hosting (generous free tier: small VM + bandwidth).
- DB: Start with on‑disk SQLite; migrate to Turso (libSQL) when ready.
- Outcome: Single container, TLS, custom domain, GitHub Actions deploy, and a clear path to Turso.

## Fly.io Deployment
1) Prerequisites
- Install: `flyctl auth login`
- Repo contains: Dockerfile, `fly.toml` (internal port `8080`)

2) Dockerfile (multi‑stage)
```Dockerfile
# builder
FROM golang:1.22-bookworm AS build
RUN apt-get update && apt-get install -y build-essential sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY . .
ENV CGO_ENABLED=1
RUN go build -o /out/notebook ./cmd/notebook

# runtime
FROM debian:bookworm-slim
RUN useradd -u 10001 -m app
WORKDIR /app
COPY --from=build /out/notebook ./notebook
COPY content/ ./content/
COPY internal/view/assets/ ./internal/view/assets/
USER app
ENV PORT=8080
EXPOSE 8080
CMD ["./notebook"]
```

3) fly.toml essentials
```toml
app = "notebook-oceanheart-ai"
[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
```

4) First deploy
```bash
flyctl apps create notebook-oceanheart-ai
# Set production env (SQLite for now)
flyctl secrets set \
  ENV=prod PORT=8080 \
  DB_PATH=/data/notebook.db \
  CONTENT_DIR=/app/content \
  SITE_BASEURL=https://notebook.oceanheart.ai \
  SITE_TITLE="Oceanheart Notebook"

# Optional: attach a small volume for SQLite
flyctl volumes create data --region <iad> --size 3

flyctl deploy
flyctl certs add notebook.oceanheart.ai   # follow DNS instructions
```

5) CI (optional)
- Add a GitHub Actions workflow to run `flyctl deploy --remote-only` on `main`.

## Turso Migration (libSQL)
Target: Replace on‑disk SQLite with Turso’s hosted libSQL. Minimal code changes; keep `database/sql`.

1) Create DB and get credentials
```bash
# Install turso CLI first
# Create database and auth token
turso db create notebook-prod
turso db tokens create notebook-prod
# Get URL
turso db show notebook-prod --url
```
Store as env vars (secrets):
- `DB_URL=libsql://<DATABASE>.turso.io`
- `DB_AUTH_TOKEN=<token>`

2) Choose client library
- Embedded replicas (CGO, best DX + offline): `go get github.com/tursodatabase/go-libsql`
- Remote‑only (no CGO): `go get github.com/tursodatabase/libsql-client-go/libsql`

3) Code outline (remote‑only)
```go
import (
  "database/sql"
  _ "github.com/tursodatabase/libsql-client-go/libsql"
)

func openTurso() (*sql.DB, error) {
  url := os.Getenv("DB_URL") + "?authToken=" + os.Getenv("DB_AUTH_TOKEN")
  return sql.Open("libsql", url)
}
```
- Update your store open path to use `sql.Open("libsql", dsn)` when `DB_URL` is present; otherwise fall back to SQLite path (`DB_PATH`).
- Reuse existing schema migrations by executing the same SQL against Turso (once) using a small bootstrap routine or the CLI.

4) Code outline (embedded replicas)
```go
import (
  "database/sql"
  "path/filepath"
  "github.com/tursodatabase/go-libsql"
)

func openReplica(dbDir string) (*sql.DB, *libsql.EmbeddedReplicaConnector, error) {
  dbPath := filepath.Join(dbDir, "replica.db")
  primary := os.Getenv("DB_URL")
  token := os.Getenv("DB_AUTH_TOKEN")
  c, err := libsql.NewEmbeddedReplicaConnector(dbPath, primary, libsql.WithAuthToken(token))
  if err != nil { return nil, nil, err }
  db := sql.OpenDB(c)
  return db, c, nil
}
// In dev you can call c.Sync() periodically, or use WithSyncInterval.
```

5) Fly secrets for Turso
```bash
flyctl secrets set \
  DB_URL=libsql://<DATABASE>.turso.io \
  DB_AUTH_TOKEN=<token>
# Optional: remove DB_PATH and volume after cutover
```

6) Migration & Cutover
- Pre‑flight: run the app locally against Turso; verify reads/writes and migrations.
- Deploy with Turso secrets set; keep SQLite vars in place as rollback.
- Verify healthz, content pages, sitemap/feed; then remove SQLite volume/env.
- Rollback: unset Turso env, restore `DB_PATH` + volume, deploy.

## Backups & Ops
- Fly: keep content in image; DB lives in Turso.
- Turso: use CLI or dashboard to manage backups/snapshots and access tokens.
- Logs: `flyctl logs -a notebook-oceanheart-ai`.

## Security
- Keep tokens in Fly secrets.
- Restrict Turso token scope; rotate regularly.
- Do not commit credentials to the repo.

## Local Dev Options
- Continue with on‑disk SQLite (`DB_PATH=./notebook.dev.db`).
- Or use Turso embedded replica for closer parity and offline sync.

