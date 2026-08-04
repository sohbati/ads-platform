#!/usr/bin/env bash
# call "ads-platform-cache-service"
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-cache-service'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CACHE_RUN="$ROOT_DIR/ads-platform-cache-service/run.sh"

if [[ ! -x "$CACHE_RUN" ]]; then
  if [[ -f "$CACHE_RUN" ]]; then
    chmod +x "$CACHE_RUN"
  else
    echo "Error: $CACHE_RUN not found" >&2
    exit 1
  fi
fi

exec "$CACHE_RUN"
