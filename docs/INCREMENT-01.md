# Increment 01 — Foundation verification report

## Outcome

**COMPLETADO:** repository audit, executable Go foundation, real frontend service
status navigation, initial schema integration tests and local backup/restore
verification. **EN PROGRESO:** phase 2 schema and phase 10 deployment preparation.
This is not a completed authentication or educational application.

## Created files

- Root: `.gitignore`, `README.md`, `render.yaml`.
- Backend: `go.mod`, `go.sum`, `.env.example`;
  `cmd/api/main.go`, `cmd/api/main_test.go`;
  `internal/config/config.go`, `internal/config/config_test.go`;
  `internal/db/db.go`;
  `internal/handler/handler.go`, `health.go`, `response.go`, `response_test.go`;
  `internal/middleware/cors.go`, `middleware.go`, `middleware_test.go`;
  `internal/models/user.go`, `internal/pkg/errors/errors.go`;
  `internal/repositories/user_repository.go`;
  `internal/routes/router.go`, `internal/routes/router_test.go`;
  `migrations/migrations.go`, `migrations/sql/000001_create_users.up.sql`.
- Frontend: `.env.example`, `src/lib/api.js`, `src/pages/ServiceStatus.jsx`,
  `src/App.test.jsx`, `src/test/setup.js`.
- Operations: `scripts/backup-db.sh`, `restore-db.sh`, `pg-tool.py`,
  `test_pg_tool.py`.
- Documentation: `docs/AUDIT.md`, `DEPLOYMENT.md`, `TECHNICAL.md`, this report.

## Modified existing files

`frontend/package.json`, `frontend/package-lock.json`, `frontend/vite.config.js`,
`frontend/index.html`, `frontend/src/App.jsx`, `frontend/src/main.jsx`,
`frontend/src/index.css`. Existing assets, unused template stylesheet and the
welcome layout were not deleted. Generated `dist/` is ignored.

## Working features

- Welcome → status page → retry/back navigation; unknown-route handling.
- `/api/v1/health`: process liveness. `/api/v1/ready`: actual PostgreSQL ping;
  returns 503 on failure, without exposing credentials.
- Explicit migration invocation with repeat-safe migration tracking.
- Users repository creation/lookup tested with PostgreSQL (not public registration).
- Strict env validation, production TLS requirement and exact CORS origins;
  sanitized errors/logs, body limit, timeouts and graceful shutdown.
- Backup/restore scripts with archive validation, private permissions,
  non-overwrite protection and refusal to restore into nonempty databases.

## Executed verification — 2026-10-26

| Check | Result |
| --- | --- |
| Baseline frontend ESLint + build | PASS |
| Final `go test -race -count=1 -cover ./...` with real TEST_DATABASE_URL | PASS; 9 top-level tests, including configuration subcases |
| `go vet ./...` and `go build ./...` | PASS |
| Initial migration on fresh PostgreSQL 18.6 DB and repeat application | PASS |
| Repository persistence, unique email and role constraint rejection | PASS |
| Readiness with connected and closed DB | PASS: 200 / 503 |
| Native API process: HTTP health/readiness and allowed CORS | PASS |
| Native API process: untrusted Origin | PASS: rejected with 403 |
| Native API SIGTERM shutdown | PASS: exit 0 |
| `npm run test` | PASS: 5 React component/navigation tests |
| `npm run lint` and `npm run build` | PASS |
| `npm audit --audit-level=moderate` | PASS: 0 reported vulnerabilities |
| `python3 -m unittest discover -s scripts -p 'test_*.py'` | PASS: 4 tests |
| `bash -n scripts/backup-db.sh scripts/restore-db.sh` | PASS |
| Actual pg_dump → archive → separate empty DB restoration | PASS: user restored, migration version 1, dirty=false |
| Existing backup overwrite / nonempty restore target | PASS: both safely refused (exit 2) |
| Backup permissions | PASS: 0600 |
| Clean `npm ci`, then tests/lint/build | PASS |
| Focused private-key/token and ignored-artifact check | PASS; not a comprehensive secret audit |

The actual database tests used a newly initialized cluster in a unique directory
under `/tmp/opencode`, bound only to loopback on port 55439. It did not use,
modify or reconfigure the existing PostgreSQL instance on port 5432. Local trust
authentication was limited to this isolated test cluster; it is not production
configuration. Integration tests without TEST_DATABASE_URL explicitly skip.
The temporary PostgreSQL cluster and native API test process were both stopped
after verification. Test artifacts remain outside the repository under
`/tmp/opencode`; they are not application storage or a persistent development DB.

Component tests intentionally mock the network function. They verify UI behavior,
not live browser-to-database end-to-end integration. The native API smoke test and
PostgreSQL integration checks were executed separately. No browser visual audit,
Render deployment validation or complete educational flow is claimed.

## Failures discovered and corrected

- Early backend drafts had compile errors, duplicated response helpers, invalid
  GORM close calls and incorrect migration embed/API usage. Corrected before the
  final successful checks. Unfinished authentication drafts were removed; no
  existing user implementation was removed.
- Initial migration filenames did not follow golang-migrate's `.up.sql`
  convention. Corrected and tested against a fresh database.
- Migration initially ran automatically at server startup and did not treat
  ErrNoChange as success. Replaced with explicit `--migrate` operation.
- CORS/log/recovery middleware was registered repeatedly. Reduced to one chain;
  production no longer inherits local development origins.
- Existing welcome buttons did nothing. Replaced with real links and an explicit
  notice that learning features remain unavailable.
- Initial backup verification failed because these clients do not expand a
  PostgreSQL URI from PGDATABASE alone. Added a tested Python stdlib helper that
  maps URL components to libpq environment variables without secret process args.
  Backup/restoration and refusal checks then passed. The failed partial archive
  was not presented as a valid backup.

## Pending / blocked

- **PENDIENTE:** phases 3–9, course data, secure auth/session lifecycle and full
  registration → activity → persisted progress demonstration.
- **PENDIENTE:** server score, XP, authorization and timezone streak tests; these
  modules do not exist yet, so no successful tests are claimed for them.
- **PENDIENTE:** adventure-style illustration and visual redesign, accessibility
  audit, mobile browser tests, recovery-email provider, PWA and other games.
- **PENDIENTE:** external Render blueprint validation, deployment smoke tests and
  public production security review. YAML preparation alone is not deployment.
- **BLOQUEADO:** using the existing local PostgreSQL instance without a dedicated
  authorized URL/role. Isolated PostgreSQL testing was possible without it.

## Next stable increment

Complete authentication as a separate block: validated registration with bcrypt,
secure revocable sessions, atomic refresh/revocation, login rate limiting, role
authorization from the database, React forms/protected routes, and actual
PostgreSQL tests for duplicate requests and session lifecycle. Then build the
CEFR/A1 content schema and a genuine first learning flow.

No commit, push, service creation or production migration was performed. Candidate
commit after review: `feat: establish tested API and PostgreSQL foundation`.
Production migration/restoration and any external deployment require explicit
authorization. GitHub upload requires the owner's affirmative authorization.
