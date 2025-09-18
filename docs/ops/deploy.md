# Deploy Guide

This guide sets up production at `https://notebook.oceanheart.ai` with Caddy, systemd, and a simple publish workflow. Local dev runs at `http://notebook.lvh.me:8003`.

## 1) Prepare the host
```bash
# Create directories
sudo mkdir -p /opt/notebook/content
sudo mkdir -p /var/lib/notebook/{prod,dev}

# Owner/group (adjust as needed)
sudo chown -R www-data:www-data /opt/notebook /var/lib/notebook
```

## 2) Environment file
Create `/etc/notebook.env`:
```bash
ENV=prod
PORT=8003
DB_PATH=/var/lib/notebook/prod/notebook.db
CONTENT_DIR=/opt/notebook/content
SITE_BASEURL=https://notebook.oceanheart.ai
SITE_TITLE=Oceanheart Notebook
# Optional: set to enable reload endpoint in prod
# RELOAD_TOKEN=changeme-strong-token
```

## 3) systemd service
Copy the sample unit and enable:
```bash
sudo cp docs/ops/notebook.service /etc/systemd/system/notebook.service
sudo systemctl daemon-reload
sudo systemctl enable --now notebook
sudo systemctl status notebook
```

## 4) Reverse proxy (Caddy)
Install Caddy (from official repos), then:
```bash
sudo cp docs/ops/Caddyfile /etc/caddy/Caddyfile
sudo systemctl reload caddy || sudo systemctl restart caddy
```

Caddyfile (provided in docs/ops/Caddyfile) minimally proxies to `127.0.0.1:8003` with gzip and automatic TLS.

## 5) Local development
```bash
ENV=dev PORT=3010 \
DB_PATH=./notebook.dev.db \
CONTENT_DIR=./content \
SITE_BASEURL=http://notebook.lvh.me:8003 \
go run ./cmd/notebook
```
Visit: http://notebook.lvh.me:8003 (drafts visible).

## 6) Publish workflow
Use the helper script to sync content and binary, then restart the service:
```bash
export NOTEBOOK_HOST=your.server
export NOTEBOOK_USER=youruser           # optional
export NOTEBOOK_REMOTE_DIR=/opt/notebook # optional
./scripts/publish.sh
```

## 7) Optional: reload content without restart
- Dev: call `GET /admin/reload`
- Prod: set `RELOAD_TOKEN` in `/etc/notebook.env`, then:
```bash
curl -H "X-Reload-Token: $RELOAD_TOKEN" https://notebook.oceanheart.ai/admin/reload
```

## 8) Backups
```bash
sudo mkdir -p /var/backups/notebook
sudo sqlite3 /var/lib/notebook/prod/notebook.db \
  ".backup '/var/backups/notebook/notebook-$(date +%F).db'"
```
Add a cronjob and prune old backups with `find -mtime`.

## 9) Logs
```bash
journalctl -u notebook -f
```
