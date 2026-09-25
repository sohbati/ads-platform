#!/usr/bin/env bash
set -euo pipefail

printf '\033]0;%s\007' 'crm-ui'

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

if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "Error: npm is not installed or not on PATH" >&2
  exit 1
fi

if [[ ! -d node_modules ]]; then
  echo "Installing crm-ui dependencies..."
  npm install
fi

PORT="${PORT:-8100}"
echo "Starting crm-ui on port ${PORT}..."
exec npm run dev -- --port "${PORT}"
