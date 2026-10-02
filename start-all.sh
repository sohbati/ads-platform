#!/usr/bin/env bash
# Open each numbered server launcher in its own Cursor terminal tab.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CURSOR_BIN="/Applications/Cursor.app/Contents/Resources/app/bin/cursor"

if [[ ! -d "/Applications/Cursor.app" ]]; then
  echo "Error: Cursor.app not found" >&2
  exit 1
fi

if ! pgrep -x Cursor >/dev/null 2>&1; then
  if [[ -x "$CURSOR_BIN" ]]; then
    "$CURSOR_BIN" -r "$ROOT_DIR" >/dev/null 2>&1 &
  else
    open -a Cursor "$ROOT_DIR"
  fi
  sleep 3
fi

open -a Cursor

open_cursor_tab() {
  local command="$1"

  if ! osascript - "$command" <<'APPLESCRIPT'
on run argv
  set theCommand to item 1 of argv
  tell application "Cursor" to activate
  delay 0.35
  tell application "System Events"
    tell process "Cursor"
      set frontmost to true
      try
        click menu item "New Terminal" of menu "Terminal" of menu bar 1
      on error
        keystroke "`" using {control down, shift down}
      end try
    end tell
  end tell
  delay 0.55
  tell application "System Events"
    tell process "Cursor"
      keystroke theCommand
      key code 36
    end tell
  end tell
end run
APPLESCRIPT
  then
    echo "Error: could not drive Cursor terminals." >&2
    echo "Allow Cursor under System Settings → Privacy & Security → Accessibility, then retry." >&2
    exit 1
  fi
}

started=0
# sort -V keeps 10 after 9 instead of after 1.
for script in $(printf '%s\n' "$ROOT_DIR"/[0-9]*.sh | sort -V); do
  name="$(basename "$script")"
  if [[ "$name" == 11-* ]]; then
    continue
  fi
  if [[ ! -x "$script" ]]; then
    chmod +x "$script"
  fi
  echo "  $name"
  open_cursor_tab "$script"
  started=$((started + 1))
  sleep 0.35
done

if [[ "$started" -eq 0 ]]; then
  echo "Error: no numbered server start scripts found in $ROOT_DIR" >&2
  exit 1
fi

echo "Opened $started Cursor terminal tabs."
