#!/usr/bin/env bash
# call "ads-platform-cdn"
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-cdn'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CDN_RUN="$ROOT_DIR/ads-platform-cdn/run.sh"

if [[ ! -x "$CDN_RUN" ]]; then
  if [[ -f "$CDN_RUN" ]]; then
    chmod +x "$CDN_RUN"
  else
    echo "Error: $CDN_RUN not found" >&2
    exit 1
  fi
fi

exec "$CDN_RUN"
