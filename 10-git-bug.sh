#!/usr/bin/env bash
# call "git-bug" — local issue tracker web UI for this repo
set -euo pipefail

printf '\033]0;%s\007' 'git-bug'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUG_RUN="$ROOT_DIR/git-bug/run.sh"

if [[ ! -x "$BUG_RUN" ]]; then
  if [[ -f "$BUG_RUN" ]]; then
    chmod +x "$BUG_RUN"
  else
    echo "Error: $BUG_RUN not found" >&2
    exit 1
  fi
fi

exec "$BUG_RUN"
