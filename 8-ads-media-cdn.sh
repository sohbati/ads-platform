#!/usr/bin/env bash
# call "ads-media-cdn" — nginx in front of MinIO
set -euo pipefail

printf '\033]0;%s\007' 'ads-media-cdn'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CDN_RUN="$ROOT_DIR/ads-media-cdn/run.sh"

if [[ ! -x "$CDN_RUN" ]]; then
  if [[ -f "$CDN_RUN" ]]; then
    chmod +x "$CDN_RUN"
  else
    echo "Error: $CDN_RUN not found" >&2
    exit 1
  fi
fi

exec "$CDN_RUN"
