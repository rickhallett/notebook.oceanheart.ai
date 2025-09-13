# Deployment Strategy (PRD)

## Executive Summary
Define a simple, reliable way to run the blog on a remote host behind a reverse proxy with clear dev/prod separation for SQLite DBs. Local authoring runs at `http://notebook.lvh.me:3010`; remote production is `https://notebook.oceanheart.ai`. Publishing is a one-command rsync + restart.

## Problem Statement
We need a developer-friendly setup to: (1) run locally with drafts, (2) serve production with TLS, (3) manage separate SQLite DB files, and (4) publish new posts quickly without complex CI.

## Requirements
- Domains: prod `notebook.oceanheart.ai`; local `notebook.lvh.me:3010`.
- Reverse proxy: automatic TLS, small config, gzip, security headers.
- DB: SQLite with distinct files for dev/prod; safe backups.
- Ops: systemd service; health checks; easy restart/rollback.
- Workflow: write locally, preview, publish remotely via a single command.

## Implementation Phases
- Phase 1 (MVP): Caddy reverse proxy, systemd unit, rsync publishing, DB backup.
- Phase 2: Git-based deploy hook (optional), zero-downtime restart, log rotation.
- Phase 3: Admin endpoint for on-demand content reload (future), metrics.

## Implementation Notes
### Remote Host Layout (Linux)
- Binary: `/opt/notebook/notebook`
- Content dir: `/opt/notebook/content/`
- DBs: `/var/lib/notebook/prod/notebook.db`, `/var/lib/notebook/dev/notebook.db`
- Config: `/etc/notebook.env` (owned by root, 0600)

Example `/etc/notebook.env`:
```
ENV=prod
PORT=8080
DB_PATH=/var/lib/notebook/prod/notebook.db
CONTENT_DIR=/opt/notebook/content
SITE_BASEURL=https://notebook.oceanheart.ai
SITE_TITLE=Oceanheart Notebook
```

### systemd service
`/etc/systemd/system/notebook.service`:
```
[Unit]
Description=Oceanheart Notebook
After=network.target

[Service]
EnvironmentFile=/etc/notebook.env
ExecStart=/opt/notebook/notebook
WorkingDirectory=/opt/notebook
User=www-data
Group=www-data
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
Commands:
```
sudo systemctl daemon-reload
sudo systemctl enable --now notebook
sudo systemctl status notebook
```

### Reverse Proxy (Caddy)
`/etc/caddy/Caddyfile`:
```
notebook.oceanheart.ai {
  encode gzip
  reverse_proxy 127.0.0.1:8080
}
```
Notes: Automatic TLS via Let’s Encrypt; minimal config. Alternative: nginx with Certbot.

### Local Development
No proxy needed: `lvh.me` resolves to `127.0.0.1`.
```
ENV=dev PORT=3010 \
DB_PATH=./notebook.dev.db \
CONTENT_DIR=./content \
SITE_BASEURL=http://notebook.lvh.me:3010 \
go run ./cmd/notebook
```
Open: `http://notebook.lvh.me:3010` (drafts visible in dev mode).

### Publish Workflow (simple and effective)
- Build locally (optional): `go build -o notebook ./cmd/notebook`
- Sync content and binary, then restart service:
```
# From repo root
rsync -avz --delete content/ user@your-host:/opt/notebook/content/
rsync -avz notebook user@your-host:/opt/notebook/notebook
ssh user@your-host 'sudo systemctl restart notebook && systemctl --no-pager --full status notebook'
```
Make it one command by adding a local script `scripts/publish.sh` that runs the above.

### Git-based Deploy (optional Phase 2)
- Keep server’s `/opt/notebook` as a Git repo.
- Push to `main`; a post-receive hook runs `go build`, copies binary, and restarts.
Pros: provenance; Cons: more moving parts than rsync.

### Backups & Maintenance
- Backup DB daily with sqlite backup to avoid locking issues:
```
sqlite3 /var/lib/notebook/prod/notebook.db \
  ".backup '/var/backups/notebook/notebook-$(date +%F).db'"
```
- Cron: `0 3 * * * root /usr/bin/sqlite3 ...` (rotate old backups with `find -mtime`).
- Logs: use `journalctl -u notebook -f`.

## Security Considerations
- Run as non-root (`www-data`), least-privilege directories.
- Caddy terminates TLS; app runs on localhost.
- No secrets in repo; use `/etc/notebook.env`.
- Drafts hidden in prod (`ENV=prod`).

## Success Metrics
- Deploy in <10 minutes from clean host.
- Publish <30 seconds end-to-end (rsync + restart).
- Zero 5xx errors during normal operation; healthy `/healthz`.

## Future Enhancements
- Add `SIGHUP` or admin endpoint to reload content without restart.
- Blue/green or socket-activation for zero-downtime restarts.
- Observability: request logs to file, metrics endpoint, uptime alert.
