#!/usr/bin/env bash
# call "ads-platform-crm" — CRM and admin for the ads platform
set -euo pipefail

printf '\033]0;%s\007' 'ads-platform-crm'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CRM_RUN="$ROOT_DIR/ads-platform-crm/run.sh"

if [[ ! -x "$CRM_RUN" ]]; then
  if [[ -f "$CRM_RUN" ]]; then
    chmod +x "$CRM_RUN"
  else
    echo "Error: $CRM_RUN not found" >&2
    exit 1
  fi
fi

exec "$CRM_RUN"
