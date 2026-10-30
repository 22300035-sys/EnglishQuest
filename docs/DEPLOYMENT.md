# Render deployment and data safety

## Authorization boundary

`render.yaml` is a preparation blueprint only. No external resources have been
created. Obtain owner authorization before deploying or creating resources, and
separate explicit approval before production migrations or restoration. Do not
select paid plans. Verify current free-plan eligibility before any creation.

## Services

- Static Site: `frontend`, `npm ci && npm run build`, publish `dist`.
  SPA rewrite sends application routes to `/index.html`.
- Go Web Service: `backend`, `go build -o bin/api ./cmd/api`, `./bin/api`.
  The server uses Render's `PORT`, timeouts, graceful shutdown and a five-connection
  database pool. Probe `/api/v1/ready` checks PostgreSQL, not course availability.
- PostgreSQL: free-plan configuration in blueprint. No permanent data may be
  stored on the Web Service's ephemeral filesystem.

## Environment variables

Backend: `APP_ENV=production`, a private `DATABASE_URL` with `sslmode=require`
(or certificate verification), a random `JWT_SECRET` of at least 32 characters,
and `CORS_ALLOWED_ORIGINS=https://YOUR-STATIC-SITE.onrender.com` (exact origin,
no trailing slash). Existing JWT configuration names are retained, but auth uses
opaque revocable sessions: `JWT_REFRESH_TTL` sets the absolute session lifetime
(default 72h). No raw tokens or password hashes are returned to the frontend.
Only explicit configured origins are accepted in production. Credentialed CORS
and Axios cookies are enabled; unsafe methods require a trusted Origin header.

Production session cookies are `__Host-eq_session`, Secure, HttpOnly,
SameSite=None and Partitioned (CHIPS), with Path=/ and no Domain attribute.
Separate onrender.com hosts can be cross-site: test signup/login, reload and
logout on the actual URLs in target browsers before accepting deployment.
Partitioned-cookie support is browser-dependent; same-site HTTPS custom domains
can improve compatibility but no domain purchase is assumed or authorized.
The client confirms `/auth/me` after login instead of assuming the browser stored
the cookie. Local HTTP uses an unpartitioned SameSite=Lax development cookie;
keep frontend/API on the same hostname. Do not weaken CSRF/CORS to hide a cookie
configuration error.

Frontend: `VITE_API_URL=https://YOUR-API.onrender.com/api/v1` before build.
Vite variables are public: never store credentials or private URLs in them.
Changing Vite env requires rebuilding. Cold starts on free services can exceed
the UI's 8-second timeout; the status screen supports retry.

## Production migration procedure (not executed)

1. Back up any existing database; verify a restoration in a separate empty DB.
2. Review all unapplied files in `backend/migrations/sql/` with the owner:
   - 000001: users, constrained role/CEFR, timezone, hashed sessions and indexes.
   - 000002: CEFR catalogue, units/lessons/activities/questions/options,
     attempts, progress, XP ledger, achievements, missions, streaks and admin audit.
   - 000003: original introductory content (three units/nine lessons in A1/A2).
   They do not drop existing tables; each migration uses a transaction.
3. Obtain explicit authorization for the target production database and changes.
4. From an authorized workstation with production variables exported privately,
   run `go run ./cmd/api --migrate` inside `backend/`. Free services need not
   provide a remote shell/predeploy hook; no paid feature is assumed.
5. Verify version and dirty flag in `schema_migrations`. Re-running an already
   applied version is safe. Never force a dirty migration version without
   diagnosing the failure and obtaining authorization.
6. Start/restart the service and verify readiness, signup/login, each game,
   persisted results, CORS/CSRF and logout using the actual frontend origin.
7. Register the owner's account normally. Only after explicit owner approval,
   run `go run ./cmd/api --promote-admin OWNER_EMAIL` from a trusted workstation
   with production variables. There are no default admin passwords or public
   role-promotion endpoints. Admin edits and student activation changes mutate
   production data and must remain limited to authorized personnel.

The process-local authentication limiter does not trust arbitrary forwarded IPs.
On a proxy deployment, requests may share the proxy IP/rate budget; review Render's
documented trusted proxy topology before configuring it, or introduce a shared
free self-hosted rate-limit store. Do not trust all forwarded IP headers.

## Free PostgreSQL expires after 30 days

Render's free PostgreSQL database has a limited 30-day lifespan. Treat expiration
as data loss risk, not a permanent persistence solution. Check the dashboard for
the precise expiration date, set reminders well before it, and export regularly.
Do not assume automatic backups or a grace period. A sustainable replacement
free-tier destination must be evaluated separately; do not silently upgrade to
a paid plan or automatically create replacement services.

## Backup

Use Python 3 and a PostgreSQL client at least as new as the server. The helper
splits the connection URL into libpq environment variables so credentials do not
appear in command-line arguments. Local processes with permission to inspect
your environment can still access them; use a trusted workstation.
Export DATABASE_URL through
a trusted secret mechanism, not command history. Disable shell tracing. Use a
private destination outside the repository, preferably encrypted storage.

```bash
bash scripts/backup-db.sh /secure/location/englishquest-YYYY-MM-DD.dump
```

The script uses custom pg_dump format, restrictive file permissions, refuses
overwrite, and checks archive structure. It does not prove recoverability until
a separate restoration succeeds. The archive contains personal data: encrypt
at rest, limit access/retention and never commit it. No backup is stored on a
Render Web Service.

## Restore (explicit, empty target only)

Obtain authorization first for any production target. Provision an empty
authorized database without deleting/replacing an existing one. Export its
connection URL as `RESTORE_DATABASE_URL` privately, then:

```bash
bash scripts/restore-db.sh /secure/location/englishquest-YYYY-MM-DD.dump --confirm-empty-target
```

The script refuses targets with user tables, never uses `--clean`, and restores
in a single transaction with SQL error checking and without owner/ACL reassignment.
Verify row counts, constraints and API behavior before switching DATABASE_URL.
Restore only trusted archives: dumps contain executable SQL. Restoration does
not merge databases, preserve role grants or automatically switch the application.
