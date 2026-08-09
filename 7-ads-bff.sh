#!/usr/bin/env bash
# call "ads-bff"
set -euo pipefail

printf '\033]0;%s\007' 'ads-bff'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BFF_RUN="$ROOT_DIR/ads-bff/run.sh"

if [[ ! -x "$BFF_RUN" ]]; then
  if [[ -f "$BFF_RUN" ]]; then
    chmod +x "$BFF_RUN"
  else
    echo "Error: $BFF_RUN not found" >&2
    exit 1
  fi
fi

exec "$BFF_RUN"
