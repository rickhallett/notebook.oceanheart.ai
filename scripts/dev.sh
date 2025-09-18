#!/usr/bin/env bash
set -euo pipefail

# Developer hot-reload script
# - Restarts server on Go changes
# - Triggers content reload on Markdown changes
#
# Requirements: one of watchexec, reflex, or entr

PORT=${PORT:-8003}
DB_PATH=${DB_PATH:-"./notebook.dev.db"}
CONTENT_DIR=${CONTENT_DIR:-"./content"}
SITE_BASEURL=${SITE_BASEURL:-"http://notebook.lvh.me:${PORT}"}

has_cmd() { command -v "$1" >/dev/null 2>&1; }

echo "Dev server: http://localhost:${PORT}"

cleanup() {
  jobs -p | xargs -r kill 2>/dev/null || true
}
trap cleanup EXIT INT TERM

if has_cmd watchexec; then
  echo "Using watchexec for watching."
  watchexec -w "${CONTENT_DIR}" -e md -- \
    sh -c "curl -fsS http://localhost:${PORT}/admin/reload >/dev/null || true" &

  exec watchexec -r -e go -- \
    env ENV=dev PORT=${PORT} DB_PATH="${DB_PATH}" CONTENT_DIR="${CONTENT_DIR}" SITE_BASEURL="${SITE_BASEURL}" \
    go run ./cmd/notebook

elif has_cmd reflex; then
  echo "Using reflex for watching."
  (
    cd "${CONTENT_DIR}" && \
    reflex -r '\.md$' -- sh -c "curl -fsS http://localhost:${PORT}/admin/reload >/dev/null || true"
  ) &

  exec reflex -r '\.(go)$' -- \
    sh -c "ENV=dev PORT=${PORT} DB_PATH='${DB_PATH}' CONTENT_DIR='${CONTENT_DIR}' SITE_BASEURL='${SITE_BASEURL}' go run ./cmd/notebook"

elif has_cmd entr; then
  echo "Using entr for watching."
  # Content reload watcher
  find "${CONTENT_DIR}" -type f -name '*.md' | entr -p sh -c "curl -fsS http://localhost:${PORT}/admin/reload >/dev/null || true" &

  # Server restart on Go changes
  find . -type f -name '*.go' | \
    entr -r sh -c "ENV=dev PORT=${PORT} DB_PATH='${DB_PATH}' CONTENT_DIR='${CONTENT_DIR}' SITE_BASEURL='${SITE_BASEURL}' go run ./cmd/notebook"

else
  echo "ERROR: Please install one of: watchexec, reflex, or entr." >&2
  echo "Examples:" >&2
  echo "  brew install watchexec || cargo install watchexec-cli" >&2
  echo "  brew install reflex || go install github.com/cespare/reflex@latest" >&2
  echo "  brew install entr || sudo apt install entr" >&2
  exit 1
fi
