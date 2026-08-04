#!/usr/bin/env bash
# call "ads-platform-ui"
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-ui'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UI_RUN="$ROOT_DIR/ads-platform-ui/run.sh"

if [[ ! -x "$UI_RUN" ]]; then
  if [[ -f "$UI_RUN" ]]; then
    chmod +x "$UI_RUN"
  else
    echo "Error: $UI_RUN not found" >&2
    exit 1
  fi
fi

exec "$UI_RUN"
