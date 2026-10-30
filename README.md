# EnglishQuest

A real educational prototype built with React/JSX, Vite, Tailwind CSS 4, React
Router, Axios, Lucide, Go, Gin, GORM and PostgreSQL. No Docker or paid service.

## Implemented prototype

- Registration, login, revocable HttpOnly sessions, logout and protected routes.
- Student/admin authorization enforced by Go, strict Origin checks and auth rate limits.
- Profile updates, password changes with session revocation and immutable registration timezone.
- PostgreSQL-backed catalogue with guided A1, A2, B1, B2, C1 and C2 practice routes.
- A1/A2/B1 initially open; B2/C1/C2 unlock through previous-level completion or a
  server-graded access exam (12 questions, 20 minutes, 80%, 10-minute retake wait).
- Owner-only symbolic completion diploma after all six routes' published activities,
  a chosen recipient name, durable issuance snapshot and browser print/PDF layout.
- 203 seeded activities in 41 units: 23 original activities plus 180 new guided
  lessons in 36 units, with progressive unlocks and preserved historical progress.
- Ten game types, including listening, guided writing, spelling, reading and timed
  challenges, graded by the server. 42 bundled synthetic English recordings.
- 36 private writing projects with self-review rubrics, server-side persistence
  and optional oral rehearsal instructions; no fabricated automatic writing/oral grades.
- Attempt snapshots, saved results, repeat-safe rewards, XP, perfect bonuses,
  varied daily missions, timed double-base-XP boosters, achievements and timezone-aware streaks without pets.
- Dashboard, academic statistics, history and per-level sample progress.
- Visual admin lesson/activity/question/option editor, levels/units/resource catalogue, publication control, account
  activation/deactivation, aggregate statistics and an audit trail.
- Explicit migrations, backup/restoration scripts and Render blueprint preparation.
- PWA manifest, browser-offered installation and honest offline fallback; no private
  data or API responses cached and no offline grading/rewards.

These are **guided practice routes, not academically validated complete CEFR courses
or proficiency certification**. See [curriculum scope](docs/CURRICULUM.md).
Teacher moderation, live oral assessment, email recovery, spaced repetition,
a media upload service and production deployment remain pending.
See the [presentation guide](docs/PRESENTACION.md),
the [latest verification report](docs/INCREMENT-05.md),
[technical scope](docs/TECHNICAL.md) and [user guide](docs/USER_GUIDE.md).

## Quick local start (Linux, PostgreSQL binaries installed)

Prerequisites: Go >=1.26, Node/npm compatible with Vite, PostgreSQL >=13 binaries
(`initdb`, `pg_ctl`, `psql`, `createdb`), OpenSSL. Tested with Go 1.27.3, Node
26.12.1 and PostgreSQL 18.6. Run as a normal OS user, never root.

```bash
cd frontend
npm ci
cd ..
bash scripts/dev.sh
```

Open **http://127.0.0.1:5173** and register your own account. There are no default
accounts or public admin passwords. This script creates an EnglishQuest-owned
development cluster on **127.0.0.1:55440**, with SCRAM password authentication,
private generated credentials in `.local/dev.env`, and persistent data in
`.local/pgdata` (both Git-ignored). It never changes an existing database on 5432.
It applies reviewed development migrations explicitly before starting the app.

Ctrl+C stops the API/frontend; PostgreSQL remains running. Stop only this cluster:

```bash
pg_ctl -D "$PWD/.local/pgdata" -m fast -w stop
```

To stop both the helper-started app and its PostgreSQL without deleting data:
`bash scripts/stop-dev.sh`.

Do not delete `.local` if you want to keep your development records. Use the backup
tools before any local cleanup. This directory is never production storage.

## Existing authorized development database

Instead of the helper, copy `backend/.env.example` to `backend/.env`, edit it
privately, export the trusted file, then run `go run ./cmd/api --migrate` and
`go run ./cmd/api` inside `backend`. The server never auto-migrates. Configure
`frontend/.env` from its example. Keep frontend/API on the same hostname in local
development (both `localhost` or both `127.0.0.1`) so SameSite cookies work.

## Administrative account

Register normally first. From the project root, the local owner can explicitly
promote that existing account:

```bash
set -a
source .local/dev.env
set +a
cd backend
go run ./cmd/api --promote-admin your-email@example.org
```

This is an owner-only database operation, not a public endpoint. Production role
changes require separate explicit owner authorization. The frontend never grants
permissions by itself.

## Tests

```bash
cd backend
go test -race ./...
go vet ./...
go build ./...
# Disposable development DB only; applies forward migrations and adds test data:
TEST_DATABASE_URL='postgres://USER:PASS@localhost:5432/englishquest_test?sslmode=disable' go test -race -count=1 ./...
```

Integration tests explicitly skip without TEST_DATABASE_URL. Never use a live
production URL or your demonstration database: tests edit test content and users.

```bash
cd frontend
npm run test
npm run lint
npm run build
# Real browser tests, with scripts/dev.sh already running:
npx playwright install chromium --only-shell
npm run test:e2e
```

Browser tests use real accounts/API/PostgreSQL, without request interception.
Screenshots/reports are ignored by Git. Python backup helper checks:
`python3 -m unittest discover -s scripts -p 'test_*.py'`.

See [Render/backup instructions](docs/DEPLOYMENT.md), [initial audit](docs/AUDIT.md)
and [the latest increment report](docs/INCREMENT-05.md) for execution status.
No external services, production changes, commits or pushes are authorized by
this README or performed automatically.
