#!/usr/bin/env bash
# Media CDN: nginx in front of MinIO (GET /ads-media/*).
set -euo pipefail

printf '\033]0;%s\007' 'ads-media-cdn'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

if [[ ! -f .env && -f config.example ]]; then
  cp config.example .env
  echo "Created .env from config.example"
fi

if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

MEDIA_CDN_PORT="${MEDIA_CDN_PORT:-8098}"
MINIO_ORIGIN="${MINIO_ORIGIN:-http://127.0.0.1:9000}"
MINIO_CONTAINER="${MINIO_CONTAINER:-minio}"
MINIO_ACCESS_KEY="${MINIO_ACCESS_KEY:-minioadmin}"
MINIO_SECRET_KEY="${MINIO_SECRET_KEY:-minioadmin123}"
MINIO_BUCKET="${MINIO_BUCKET:-ads-media}"

if ! command -v nginx >/dev/null 2>&1; then
  echo "Error: nginx is not installed or not on PATH" >&2
  echo "Install with: brew install nginx" >&2
  exit 1
fi

CONF="$ROOT_DIR/nginx.conf"
if [[ ! -f "$CONF" ]]; then
  echo "Error: $CONF not found" >&2
  exit 1
fi

mkdir -p "$ROOT_DIR/logs/cache" "$ROOT_DIR/logs/tmp/client_body" "$ROOT_DIR/logs/tmp/proxy"

if ! grep -q "listen ${MEDIA_CDN_PORT};" "$CONF"; then
  echo "Warning: nginx.conf listen port does not match MEDIA_CDN_PORT=${MEDIA_CDN_PORT}" >&2
fi

if ! curl -sf -o /dev/null --max-time 2 "${MINIO_ORIGIN}/minio/health/live"; then
  echo "Warning: MinIO origin ${MINIO_ORIGIN} is not reachable. Start MinIO before serving images." >&2
fi

if command -v docker >/dev/null 2>&1 && docker ps --format '{{.Names}}' 2>/dev/null | grep -qx "${MINIO_CONTAINER}"; then
  docker exec "${MINIO_CONTAINER}" mc alias set local http://127.0.0.1:9000 "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}" >/dev/null
  docker exec "${MINIO_CONTAINER}" mc anonymous set download "local/${MINIO_BUCKET}" >/dev/null
  echo "Anonymous GET enabled on MinIO bucket ${MINIO_BUCKET}"
else
  echo "Warning: MinIO container '${MINIO_CONTAINER}' not found. Bucket must allow anonymous GetObject for the CDN." >&2
fi

if lsof -nP -iTCP:"${MEDIA_CDN_PORT}" -sTCP:LISTEN >/dev/null 2>&1; then
  echo "Error: port ${MEDIA_CDN_PORT} is already in use" >&2
  exit 1
fi

echo "Starting ads-media-cdn (nginx) on :${MEDIA_CDN_PORT} -> ${MINIO_ORIGIN} ..."
exec nginx -p "$ROOT_DIR/" -c "$CONF"
