#!/usr/bin/env bash
# call "ads-platform-back"
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-back'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACK_RUN="$ROOT_DIR/ads-platform-back/run.sh"

if [[ ! -x "$BACK_RUN" ]]; then
  if [[ -f "$BACK_RUN" ]]; then
    chmod +x "$BACK_RUN"
  else
    echo "Error: $BACK_RUN not found" >&2
    exit 1
  fi
fi

exec "$BACK_RUN"
