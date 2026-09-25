#!/usr/bin/env bash
# call "crm-ui" — CRM and admin Next.js UI
set -euo pipefail

printf '\033]0;%s\007' 'crm-ui'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CRM_RUN="$ROOT_DIR/crm-ui/run.sh"

if [[ ! -x "$CRM_RUN" ]]; then
  if [[ -f "$CRM_RUN" ]]; then
    chmod +x "$CRM_RUN"
  else
    echo "Error: $CRM_RUN not found" >&2
    exit 1
  fi
fi

exec "$CRM_RUN"
