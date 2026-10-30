#!/usr/bin/env bash
# Starts ONLY an EnglishQuest-owned local cluster; never modifies port 5432.
set -euo pipefail
umask 077
root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
local_dir="$root/.local"
mkdir -p "$local_dir/socket"
chmod 700 "$local_dir" "$local_dir/socket"
if [[ ! -f "$local_dir/dev.env" ]]; then
  password=$(openssl rand -hex 24)
  secret=$(openssl rand -hex 32)
  set -o noclobber
  {
    printf 'APP_ENV=development\nPORT=8080\n'
    printf 'EQ_DB_PASSWORD=%s\n' "$password"
    printf 'DATABASE_URL=postgres://englishquest:%s@127.0.0.1:55440/englishquest?sslmode=disable\n' "$password"
    printf 'JWT_SECRET=%s\nJWT_REFRESH_TTL=72h\n' "$secret"
    printf 'CORS_ALLOWED_ORIGINS=http://127.0.0.1:5173,http://localhost:5173\n'
    printf 'VITE_API_URL=http://127.0.0.1:8080/api/v1\n'
  } > "$local_dir/dev.env"
  set +o noclobber
fi
set -a
source "$local_dir/dev.env"
set +a
if [[ ! -f "$local_dir/pgdata/PG_VERSION" ]]; then
  initdb -D "$local_dir/pgdata" --auth-local=peer --auth-host=scram-sha-256 --no-locale --encoding=UTF8 > "$local_dir/init.log"
fi
if ! pg_ctl -D "$local_dir/pgdata" status > /dev/null 2>&1; then
  pg_ctl -D "$local_dir/pgdata" -l "$local_dir/postgres.log" -o "-h 127.0.0.1 -p 55440 -k '$local_dir/socket'" -w start
fi
psql_local=(psql --no-password --no-psqlrc -h "$local_dir/socket" -p 55440 -U "$(id -un)" -d postgres -v ON_ERROR_STOP=1)
if [[ "$("${psql_local[@]}" -Atc "SELECT count(*) FROM pg_roles WHERE rolname='englishquest'")" == 0 ]]; then
  # Password flows through stdin, never command-line arguments.
  printf "CREATE ROLE englishquest LOGIN PASSWORD '%s';\n" "$EQ_DB_PASSWORD" | "${psql_local[@]}" > /dev/null
fi
if [[ "$("${psql_local[@]}" -Atc "SELECT count(*) FROM pg_database WHERE datname='englishquest'")" == 0 ]]; then
  createdb -h "$local_dir/socket" -p 55440 -U "$(id -un)" -O englishquest englishquest
fi
echo 'EnglishQuest development PostgreSQL ready on 127.0.0.1:55440.'
echo 'Private configuration: .local/dev.env (ignored by Git).'
