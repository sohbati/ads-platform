#!/usr/bin/env bash
# call "message-broker"
set -euo pipefail

printf '\033]0;%s\007' 'message-broker'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BROKER_RUN="$ROOT_DIR/message-broker/run.sh"

if [[ ! -x "$BROKER_RUN" ]]; then
  if [[ -f "$BROKER_RUN" ]]; then
    chmod +x "$BROKER_RUN"
  else
    echo "Error: $BROKER_RUN not found" >&2
    exit 1
  fi
fi

exec "$BROKER_RUN"
