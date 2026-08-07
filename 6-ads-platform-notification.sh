#!/usr/bin/env bash
# call "ads-platform-notification"
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-notification'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NOTIFICATION_RUN="$ROOT_DIR/ads-platform-notification/run.sh"

if [[ ! -x "$NOTIFICATION_RUN" ]]; then
  if [[ -f "$NOTIFICATION_RUN" ]]; then
    chmod +x "$NOTIFICATION_RUN"
  else
    echo "Error: $NOTIFICATION_RUN not found" >&2
    exit 1
  fi
fi

exec "$NOTIFICATION_RUN"
