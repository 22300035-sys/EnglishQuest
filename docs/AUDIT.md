# Repository audit and implementation plan

## Starting state (2026-10-26)

- Git: `main`, no commits, no remote configured. All existing frontend files
  were untracked. No historical baseline exists to produce a conventional diff.
- Frontend: React 19.3.1, Vite 8.3.4, Tailwind 4.3.3, Router 7.18.4,
  Axios 1.20.0, Lucide 1.60.0 (installed versions/ranges retained).
- Welcome: functional React rendering with dark/emerald branding, but two
  buttons had no actions. No authentication, navigation or API integration.
- `App.css`, SVGs, icons and `hero.png` came from the starter template.
  `hero.png` depicts a technical layered shape, not an educational landscape.
  Existing assets and unused stylesheet were preserved.
- Backend and docs initially contained empty directory scaffolding only.
  No Go module, database configuration, schema, migrations or tests existed.
- HTML still had title `frontend`. No root secret/backup ignore rules existed.
- Environment: Go 1.27.3, Node 26.12.1, npm 12.2.1, PostgreSQL binaries 18.6.
  Existing PostgreSQL on port 5432 was not accessible using the current OS role.
  Its configuration, roles and data were not changed.
- Baseline frontend build and ESLint passed. No baseline automated tests existed.

## Phased plan and gates

| Phase | Scope and completion gate | Status after first increment |
| --- | --- | --- |
| 1 | Audit, architecture, env validation, reproducible checks | COMPLETADO |
| 2 | Go API, pool, explicit SQL migrations, PostgreSQL integration | EN PROGRESO: initial users/session schema tested; educational schema pending |
| 3 | Registration, secure sessions, server authorization, protected UI; test full session lifecycle | PENDIENTE |
| 4 | CEFR catalogue, A1 micro-lessons, initial A2; validated content in PostgreSQL | PENDIENTE |
| 5 | Accessible responsive dashboard using real per-user API metrics | PENDIENTE |
| 6 | Original learning path with backend unlock requirements | PENDIENTE |
| 7 | Word Hunter, Sentence Builder, Memory Match; server grading, persisted attempts | PENDIENTE |
| 8 | Atomic progress, idempotent XP, missions, achievements, timezone-aware streaks; no pets | PENDIENTE |
| 9 | Role-authorized content/user administration | PENDIENTE |
| 10 | Cross-cutting security, tests, documentation, Render preparation | EN PROGRESO |

Every increment must build, pass its executable tests, distinguish unit fixtures
from actual persistence, and document unavailable features honestly. The current
UI is an infrastructure increment, not the requested final educational demo.

## Architecture decisions

- Keep the existing monorepo: `frontend/`, `backend/`, `docs/`.
- Go layers: config → routes/middleware → handlers → services (next phase) →
  repositories → PostgreSQL. Domain models are independent of HTTP.
- GORM for parameterized repository operations, versioned SQL for schema;
  never production AutoMigrate.
- `/health` means process liveness; `/ready` means live PostgreSQL connectivity,
  not content availability or migration/schema completeness.
- No automatic migrations on server startup. Explicit forward migration command
  only; no destructive rollback command is provided.
- Session schema is preparation only. Incomplete auth service drafts created
  during this session were removed before the stable foundation was delivered.
- Subsequent authentication must validate passwords, rate-limit requests, avoid
  browser token persistence, authorize roles from the database, implement atomic
  rotation/revocation and CSRF protections as appropriate to cookie deployment.
- Future result submission must lock attempts and write score/progress/XP in
  one transaction, with unique reward sources and server-only answer keys.
