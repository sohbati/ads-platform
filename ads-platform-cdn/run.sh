#!/usr/bin/env bash
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-cdn'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

if [[ ! -f .env ]]; then
  if [[ -f config.example ]]; then
    cp config.example .env
    echo "Created .env from config.example"
  fi
fi

if [[ ! -f "$ROOT_DIR/go.mod" ]]; then
  echo "Error: go.mod not found in $ROOT_DIR" >&2
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Error: go is not installed or not on PATH" >&2
  exit 1
fi

echo "Starting ads-platform-cdn..."
exec go run -C "$ROOT_DIR" ./cmd/server
