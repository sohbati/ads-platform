#!/usr/bin/env bash
set -euo pipefail

printf '\033]0;%s\007' 'nats-message-broker'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

if [[ ! -f .env ]]; then
  if [[ -f config.example ]]; then
    cp config.example .env
    echo "Created .env from config.example"
  else
    echo "Warning: no .env or config.example found; using built-in defaults" >&2
  fi
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Error: go is not installed or not on PATH" >&2
  exit 1
fi

echo "Starting nats-message-broker..."
exec go run ./cmd/server
