#!/usr/bin/env python3
"""Run PostgreSQL clients with URL credentials in environment, not process args.

PGDATABASE alone does not expand a connection URI for these client tools.
Only reviewed libpq URL parameters are supported; unknown ones fail closed.
"""
import os
import subprocess
import sys
from urllib.parse import parse_qs, unquote, urlsplit


def main():
    if len(sys.argv) < 3 or sys.argv[2] not in {"pg_dump", "psql"}:
        raise ValueError("Invalid PostgreSQL tool invocation")
    url = urlsplit(os.environ[sys.argv[1]])
    if url.scheme not in {"postgres", "postgresql"} or not url.hostname or not url.path.strip("/"):
        raise ValueError("Invalid PostgreSQL URL")
    # Do not inherit accidental PGHOST/PGSERVICE/PGOPTIONS overrides.
    env = {k: v for k, v in os.environ.items() if not k.startswith("PG")}
    for key in ("PGPASSFILE",):
        if key in os.environ:
            env[key] = os.environ[key]
    env.update({
        "PGHOST": url.hostname,
        "PGPORT": str(url.port or 5432),
        "PGDATABASE": unquote(url.path[1:]),
        "PGCONNECT_TIMEOUT": "10",
    })
    if url.username is not None:
        env["PGUSER"] = unquote(url.username)
    if url.password is not None:
        env["PGPASSWORD"] = unquote(url.password)
    supported = {
        "sslmode": "PGSSLMODE", "sslrootcert": "PGSSLROOTCERT",
        "sslcert": "PGSSLCERT", "sslkey": "PGSSLKEY",
        "connect_timeout": "PGCONNECT_TIMEOUT",
        "channel_binding": "PGCHANNELBINDING",
        "application_name": "PGAPPNAME",
    }
    for key, values in parse_qs(url.query, keep_blank_values=True).items():
        if key not in supported or len(values) != 1:
            raise ValueError("Unsupported or duplicated connection parameter")
        env[supported[key]] = values[0]
    return subprocess.run(sys.argv[2:], env=env, check=False).returncode


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (ValueError, KeyError, OSError):
        print("Invalid connection configuration or PostgreSQL client unavailable.", file=sys.stderr)
        sys.exit(2)
