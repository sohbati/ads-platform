#!/usr/bin/env bash
# call "ads-platform-stats" — NATS worker that rolls up ad view/contact events
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-stats'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STATS_RUN="$ROOT_DIR/ads-platform-stats/run.sh"

if [[ ! -x "$STATS_RUN" ]]; then
  if [[ -f "$STATS_RUN" ]]; then
    chmod +x "$STATS_RUN"
  else
    echo "Error: $STATS_RUN not found" >&2
    exit 1
  fi
fi

exec "$STATS_RUN"
