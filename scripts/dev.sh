#!/usr/bin/env bash
set -euo pipefail
root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
bash "$root/scripts/dev-db.sh"
set -a
source "$root/.local/dev.env"
set +a
if [[ "$APP_ENV" != development || "$DATABASE_URL" != postgres://englishquest:*@127.0.0.1:55440/englishquest\?sslmode=disable ]]; then
  echo 'Refusing migrations: dev.sh only accepts its own local EnglishQuest database on 55440.' >&2
  exit 2
fi
export CORS_ALLOWED_ORIGINS="$CORS_ALLOWED_ORIGINS,http://127.0.0.1:4173,http://localhost:4173"
(cd "$root/backend" && go build -o "$root/.local/api" ./cmd/api)
"$root/.local/api" --migrate
"$root/.local/api" &
api_pid=$!
(cd "$root/frontend" && unset DATABASE_URL JWT_SECRET EQ_DB_PASSWORD && exec ./node_modules/.bin/vite --host 127.0.0.1 --strictPort) &
web_pid=$!
printf '%s %s\n' "$api_pid" "$web_pid" > "$root/.local/app.pids"
trap 'kill "$api_pid" "$web_pid" 2>/dev/null || true' EXIT INT TERM
echo 'Open http://127.0.0.1:5173 — register your own account.'
echo 'Ctrl+C stops the app. PostgreSQL remains running to preserve your data.'
wait -n "$api_pid" "$web_pid"
