#!/usr/bin/env bash
set -euo pipefail
umask 077
: "${DATABASE_URL:?Set DATABASE_URL privately before running this script}"
if [[ $# != 1 ]]; then
  echo 'Usage: bash scripts/backup-db.sh /secure/path/backup.dump' >&2
  exit 2
fi
target=$1
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
if [[ -e "$target" ]]; then
  echo 'Refusing to overwrite an existing backup.' >&2
  exit 2
fi
# noclobber also protects against a concurrent process creating the same file.
set -o noclobber
python3 "$script_dir/pg-tool.py" DATABASE_URL pg_dump --no-password --format=custom --no-owner --no-acl > "$target"
pg_restore --list "$target" > /dev/null
echo 'Backup created and archive structure verified. Test restoration separately.'
