#!/usr/bin/env bash
set -euo pipefail
umask 077
: "${RESTORE_DATABASE_URL:?Set RESTORE_DATABASE_URL to a reviewed EMPTY target database}"
if [[ $# != 2 || "$2" != '--confirm-empty-target' || ! -f "$1" ]]; then
  echo 'Usage: bash scripts/restore-db.sh backup.dump --confirm-empty-target' >&2
  echo 'Production restoration requires owner approval. Never target a live database.' >&2
  exit 2
fi
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
tables=$(python3 "$script_dir/pg-tool.py" RESTORE_DATABASE_URL psql --no-password --no-psqlrc -Atc "SELECT count(*) FROM information_schema.tables WHERE table_schema NOT IN ('pg_catalog', 'information_schema')")
if [[ "$tables" != '0' ]]; then
  echo 'Refusing to restore: target database is not empty.' >&2
  exit 2
fi
pg_restore --list "$1" > /dev/null
# No --clean, DROP, owner reassignment or role creation. All SQL is committed
# together; SQL errors roll back the restoration instead of leaving partial data.
pg_restore --no-owner --no-acl --file=- "$1" |
  python3 "$script_dir/pg-tool.py" RESTORE_DATABASE_URL psql --no-password --no-psqlrc --set=ON_ERROR_STOP=1 --single-transaction > /dev/null
echo 'Restoration committed. Verify row counts and application behavior.'
