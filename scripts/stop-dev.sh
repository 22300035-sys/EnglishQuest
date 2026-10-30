#!/usr/bin/env bash
set -euo pipefail
root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
if [[ -f "$root/.local/app.pids" ]]; then
  read -r api_pid web_pid < "$root/.local/app.pids"
  for pid in "$api_pid" "$web_pid"; do
    [[ "$pid" =~ ^[0-9]+$ ]] || continue
    command=$(ps -p "$pid" -o args= || true)
    # Refuse to signal a reused PID belonging to a different application.
    if [[ "$command" == "$root/.local/api" || "$command" == *"$root/frontend/node_modules/.bin/vite"* ]]; then
      kill -TERM "$pid" 2>/dev/null || true
    fi
  done
fi
if [[ -f "$root/.local/pgdata/PG_VERSION" ]] && pg_ctl -D "$root/.local/pgdata" status > /dev/null 2>&1; then
  pg_ctl -D "$root/.local/pgdata" -m fast -w stop
fi
echo 'EnglishQuest local services stopped. Development data was not deleted.'
