#!/usr/bin/env bash
set -euo pipefail

# Simple publish script: sync content and binary to remote host and restart service.
# Requirements:
# - Environment variables:
#   NOTEBOOK_HOST   (e.g., example.com) [required]
#   NOTEBOOK_USER   (default: $USER)
#   NOTEBOOK_REMOTE_DIR (default: /opt/notebook)
#   NOTEBOOK_SERVICE (default: notebook)
# - SSH access with sudo rights to restart the service

NOTEBOOK_HOST=${NOTEBOOK_HOST:-}
NOTEBOOK_USER=${NOTEBOOK_USER:-"${USER}"}
NOTEBOOK_REMOTE_DIR=${NOTEBOOK_REMOTE_DIR:-"/opt/notebook"}
NOTEBOOK_SERVICE=${NOTEBOOK_SERVICE:-"notebook"}

if [[ -z "${NOTEBOOK_HOST}" ]]; then
  echo "ERROR: NOTEBOOK_HOST is required (e.g., export NOTEBOOK_HOST=your.server)" >&2
  exit 1
fi

echo "Building binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o notebook ./cmd/notebook

echo "Syncing content and binary to ${NOTEBOOK_USER}@${NOTEBOOK_HOST}:${NOTEBOOK_REMOTE_DIR}"
rsync -avz --delete content/ "${NOTEBOOK_USER}@${NOTEBOOK_HOST}:${NOTEBOOK_REMOTE_DIR}/content/"
rsync -avz notebook "${NOTEBOOK_USER}@${NOTEBOOK_HOST}:${NOTEBOOK_REMOTE_DIR}/notebook"

echo "Restarting remote service ${NOTEBOOK_SERVICE}..."
ssh "${NOTEBOOK_USER}@${NOTEBOOK_HOST}" \
  "sudo systemctl restart ${NOTEBOOK_SERVICE} && sudo systemctl --no-pager --full status ${NOTEBOOK_SERVICE} | sed -n '1,20p'"

echo "Publish complete."

