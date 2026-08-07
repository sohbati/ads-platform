#!/usr/bin/env bash
# Run integration tests with testcontainers
set -euo pipefail

printf '\033]0;%s\007' 'integration-test'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR/integration-test"

if ! command -v docker >/dev/null 2>&1; then
  echo "Error: docker is not installed or not on PATH" >&2
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  echo "Error: docker daemon is not running" >&2
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Error: go is not installed or not on PATH" >&2
  exit 1
fi

echo "Running integration tests..."
make tidy
make test-integration
