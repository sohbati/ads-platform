#!/usr/bin/env bash
# Standalone git-bug runner. Issues live in this git repo (objects, not files).
set -euo pipefail

printf '\033]0;%s\007' 'git-bug'

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$PROJECT_DIR/.." && pwd)"
BIN_DIR="$PROJECT_DIR/bin"
PORT_DEFAULT=8100

cd "$REPO_ROOT"

if [[ ! -d .git ]]; then
  echo "Error: $REPO_ROOT is not a git repository" >&2
  exit 1
fi

if [[ ! -f "$PROJECT_DIR/.env" && -f "$PROJECT_DIR/config.example" ]]; then
  cp "$PROJECT_DIR/config.example" "$PROJECT_DIR/.env"
  echo "Created git-bug/.env from config.example"
fi

if [[ -f "$PROJECT_DIR/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$PROJECT_DIR/.env"
  set +a
fi

PORT="${PORT:-$PORT_DEFAULT}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
esac

asset="git-bug_${os}_${arch}"
if [[ "$os" == "darwin" ]]; then
  asset="git-bug_darwin_${arch}"
elif [[ "$os" == "linux" ]]; then
  asset="git-bug_linux_${arch}"
fi

find_git_bug() {
  if command -v git-bug >/dev/null 2>&1; then
    command -v git-bug
    return 0
  fi
  if [[ -x "$BIN_DIR/git-bug" ]]; then
    echo "$BIN_DIR/git-bug"
    return 0
  fi
  if command -v brew >/dev/null 2>&1; then
    local prefix
    prefix="$(brew --prefix git-bug 2>/dev/null || true)"
    if [[ -n "$prefix" && -x "$prefix/bin/git-bug" ]]; then
      echo "$prefix/bin/git-bug"
      return 0
    fi
  fi
  return 1
}

install_git_bug() {
  mkdir -p "$BIN_DIR"
  local url="https://github.com/git-bug/git-bug/releases/latest/download/${asset}"
  echo "Downloading git-bug ($asset)..."
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$BIN_DIR/git-bug" "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -q -O "$BIN_DIR/git-bug" "$url"
  else
    echo "Error: install git-bug (brew install git-bug) or install curl" >&2
    exit 1
  fi
  chmod +x "$BIN_DIR/git-bug"
}

GIT_BUG_BIN="$(find_git_bug || true)"
if [[ -z "${GIT_BUG_BIN}" ]]; then
  install_git_bug
  GIT_BUG_BIN="$BIN_DIR/git-bug"
fi

export PATH="$(dirname "$GIT_BUG_BIN"):$PATH"

if ! git bug version >/dev/null 2>&1; then
  echo "Error: git-bug is not usable ($GIT_BUG_BIN)" >&2
  exit 1
fi

ensure_identity() {
  local listed
  listed="$(git bug user 2>/dev/null || true)"
  if [[ -n "$listed" ]]; then
    return 0
  fi
  local name email
  name="$(git config --get user.name || true)"
  email="$(git config --get user.email || true)"
  if [[ -z "$name" || -z "$email" ]]; then
    echo "Error: no git-bug identity. Set git user.name and user.email, then rerun." >&2
    echo "  Or: git bug user new -n 'Your Name' -e you@example.com --non-interactive" >&2
    exit 1
  fi
  echo "Creating git-bug identity from git config ($name <$email>)..."
  git bug user new -n "$name" -e "$email" --non-interactive
}

ensure_identity

echo "Starting git-bug web UI on http://127.0.0.1:${PORT}"
echo "Issues are stored as git objects in $REPO_ROOT (not as working-tree files)."
exec git bug webui --host 127.0.0.1 --port "$PORT" --no-open
