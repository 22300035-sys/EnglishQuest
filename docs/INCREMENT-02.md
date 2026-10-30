# Increment 02 — Functional educational prototype

## Outcome and scope

**COMPLETADO — MVP demonstration flow:** register → login → select A1/A2 → read
lesson → play → Go grading → saved result → academic progress/rewards → reload.
Both student and basic administrator workflows were exercised against real
PostgreSQL. There is no production mock mode or localStorage persistence.

**EN PROGRESO — full original product specification:** broader content, remaining
games, structural catalogue administration, timed boosters and production review
remain future work. This report does not describe the entire initial wish list
as 100% delivered. The prototype's A1/A2 content is an introductory sample.

## Implemented blocks

### Accounts and security

- Registration with normalized unique email, bcrypt password hashes, name and
  validated explicit timezone. Signup cannot supply ADMIN or arbitrary fields.
- Login, database-backed absolute-expiry sessions, HttpOnly cookies and logout.
  Raw random tokens are never stored in PostgreSQL or returned in JSON.
- Current active user/role checked on every protected API request. Disabled
  accounts have sessions revoked; reactivation does not resurrect them.
- Credentialed CORS, exact Origin validation for unsafe methods, bounded auth
  rate limiting, profile updates and owner-only CLI admin promotion.
- Production cookie flags checked in unit tests; real Render browser behavior
  still requires deployment-specific verification and authorization.

### Learning

- Versioned schema and seed content stored in PostgreSQL: six CEFR entries,
  three introductory units, nine published lessons/activities, normalized
  questions and options. B1–C2 remain unavailable.
- Original responsive adventure welcome and eight local SVG illustrations.
- Dashboard, CEFR selection, unit view, learning roadmap, micro-lessons, exercise,
  results, achievements, daily mission, statistics/history, profile and basic admin.
- Three real games: Word Hunter, Sentence Builder, Memory Match. Memory checks
  all proposed pairs on final submission, not by inventing local match feedback.
- Server-only grading, threshold-based unlocks, saved best scores and stars.
- Private immutable attempt snapshots. Random opaque option IDs and shuffled
  choices do not encode ordering or pair solutions. Solutions/explanations are
  withheld until the attempt is submitted.

### Rewards

- One-time activity/checkpoint XP, first-perfect bonuses, one daily mission,
  four achievement definitions, practice-day calendar and current/best streak.
- Rewards, result and academic progress commit in one transaction under row
  locks and unique-source constraints. Duplicate/concurrent requests do not farm XP.
- Streak dates use the account's timezone and civil-day arithmetic; no pets.
- Repeats can improve scores and maintain practice but cannot repeat activity XP.

### Content administration

- Add a lesson to an existing unit; edit micro-lesson, activity kind, instructions,
  questions/options/solutions, explanations and publication/evaluation flags.
- Server-side validation and an audited transaction for content edits.
- Question authoring is a functional JSON editor, not a finished visual builder.
- Basic real aggregate statistics and activation/deactivation of student accounts.

## Files created

- Migrations: `backend/migrations/sql/000002_learning.up.sql`,
  `000003_starter_content.up.sql`.
- Domain: `backend/internal/models/learning.go`.
- Repositories: `backend/internal/repositories/accounts.go`, `learning.go`,
  `admin.go`.
- Services/tests: `backend/internal/services/auth.go`, `auth_test.go`,
  `grading.go`, `grading_test.go`, `content.go`.
- HTTP/tests: `backend/internal/handler/app.go`, `app_test.go`,
  `backend/cmd/api/prototype_test.go`.
- Frontend auth: `frontend/src/auth/context.js`, `AuthProvider.jsx`.
- Shared UI: `frontend/src/components/UI.jsx`, `Shell.jsx`, `Games.jsx`,
  `Games.test.jsx`; `frontend/src/lib/useResource.js`, `frontend/src/quest.css`.
- Pages: `frontend/src/pages/Welcome.jsx`, `Auth.jsx`, `Levels.jsx`,
  `Dashboard.jsx`, `Learning.jsx`, `Lesson.jsx`, `Exercise.jsx`, `Results.jsx`,
  `ProgressPages.jsx`, `Profile.jsx`, `Admin.jsx`.
- Art: `frontend/public/images/{apple,banana,orange,book,pen,school,breakfast,night}.svg`.
- Browser checks: `frontend/playwright.config.js`,
  `frontend/e2e/learning.spec.js`, `frontend/e2e/admin.spec.js`.
- Local operations: `scripts/dev-db.sh`, `dev.sh`, `stop-dev.sh`.
- Guides: `docs/USER_GUIDE.md`, this report.

## Existing files modified

`backend/cmd/api/main.go`, `backend/internal/routes/router.go`,
`backend/internal/middleware/cors.go`, `backend/internal/handler/response.go`,
`response_test.go`, `backend/go.mod`, `go.sum`; frontend `package.json`, lockfile,
`vite.config.js`, `src/App.jsx`, `src/App.test.jsx`, `src/main.jsx`,
`src/lib/api.js`, `src/pages/ServiceStatus.jsx`, `public/favicon.svg`,
`frontend/README.md`; root `.gitignore`, `README.md`, `docs/TECHNICAL.md`,
`docs/DEPLOYMENT.md`.

The generic Vite favicon was replaced with an original EnglishQuest book icon;
the existing unused hero/template assets remain preserved. Welcome was upgraded
to the requested adventure identity, not replaced with a generic external image.

## Executed verification (2026-10-26)

| Check | Actual result |
| --- | --- |
| `go test -race -count=1 ./...` with dedicated PostgreSQL | PASS, including full API integration |
| Fresh database: apply all 3 migrations and repeat migration command | PASS |
| `go vet ./...`, `go build ./...` | PASS |
| Registration/login, password hashing, ownership, server role checks | PASS |
| Nine seeded activities covering all 3 game types | PASS |
| Concurrent repeated submissions, new repeat attempt and XP ledger | PASS: 360 total first-perfect sample XP, no replay farming |
| Mission, streak day, progress persisted in dashboard | PASS |
| Snapshot remains unchanged after admin edit | PASS |
| Admin creation/edit, disable/reactivate, session revocation and expiry | PASS |
| Grading, 60/80% thresholds, civil-date streaks, input validation | PASS |
| Production cookie flags, auth rate limit, request cancellation handling | PASS |
| `npm run test` | PASS: 13 component/form/navigation/game tests |
| `npm run lint`, `npm run build` | PASS |
| `npm audit --audit-level=moderate` | PASS: 0 reported vulnerabilities |
| Clean `npm ci` followed by tests/lint/build | PASS |
| Chromium Playwright against real API/DB, no HTTP mocks | PASS: 3 browser tests |
| Browser student flow: three games, 140 XP, mission, achievements, reload/profile/logout | PASS |
| Browser authorized admin: create hidden lesson, edit it, reload/read persisted content | PASS |
| Browser availability labels and mobile overflow assertions | PASS |
| Browser workflows assert no received API 5xx responses | PASS |
| Desktop/mobile screenshots | Captured and inspected locally; not a complete accessibility audit |
| Python backup helper tests | PASS: 4 tests |
| `bash -n scripts/*.sh` | PASS |
| Backup of new prototype schema/data → separate empty restore DB | PASS: six CEFR entries, lesson/attempt data, migration version 3 and dirty=false |
| Existing-backup and nonempty-restore-target refusal | PASS |
| Focused known-secret/token/artifact scan, including frontend build | PASS; not a comprehensive security audit |

Backend database testing used a newly initialized disposable PostgreSQL 18.6
cluster on 55439. Browser testing used the separate EnglishQuest development
cluster on 55440. The original instance on 5432 was neither changed nor used.
The disposable cluster on 55439 was stopped after verification; the persistent
development app/database on 5173/8080/55440 remain available for inspection.
Final HTTP checks returned 200 for the frontend and database readiness endpoint.

Administrative browser-test accounts were explicitly created/promoted only in
the local development database; random test passwords were supplied to the test
process environment, not committed or offered as default credentials. QA users,
attempts and hidden authored lessons remain as local test records; no destructive
cleanup was performed. The published seed remains nine activities.

Without `E2E_ADMIN_EMAIL`/`E2E_ADMIN_PASSWORD`, the optional admin browser test
explicitly skips. The recorded successful three-test execution supplied an
authorized local test account. Without TEST_DATABASE_URL, Go integration checks
also explicitly skip; a plain test run without infrastructure is not equivalent
to the recorded PostgreSQL verification.

## Failures found and corrected

- A Word Hunter component test exposed duplicate accessible naming from an image
  and label. Added an explicit button label and hid redundant visible text from
  the accessible name.
- Initial browser completion expected welcome after logout while protected routes
  redirected to login. Logout now relies on one consistent protected-route redirect.
- A repeated integration run incorrectly assumed a fixed pre-edit title. The test
  now compares the actual immutable snapshot before/after editing.
- StrictMode intentionally aborts obsolete requests. Cancellation was initially
  logged as 500/503; canceled requests now use 499, and deadlines use 504.
- The full Chromium download timed out. The smaller official headless-shell
  browser installed successfully, and the real browser tests then ran.
- Adjusted grid spacing and screenshot animation handling after inspecting mobile
  and desktop captures. The screens fit the tested 390px viewport without page overflow.

## Remaining work, not claimed complete

- **PENDIENTE:** full academically reviewed A1/A2 courses and validated B1–C2 content.
- **PENDIENTE:** Sound Detective, Grammar Rescue, Spell Master, Conversation Quest,
  Story Explorer, Word Rain and English Rush; listening/reading/writing coverage.
- **PENDIENTE:** safe email delivery/recovery, email ownership verification,
  password/email changes, account deletion/privacy workflows.
- **PENDIENTE:** timed/double-XP boosters, diverse daily missions, additional
  achievements and spaced repetition.
- **PENDIENTE:** structural level/unit management, media uploads, role assignment
  UI and a nontechnical visual question editor. Current admin covers existing
  units and lesson/activity authoring, not the entire original admin specification.
- **PENDIENTE:** PWA/offline architecture, tablet/Firefox/Safari tests, comprehensive
  keyboard/screen-reader and contrast review, production threat review.
- **BLOQUEADO POR AUTORIZACIÓN:** external deployment/service creation, actual Render
  cookie/proxy verification, production migrations/restoration and GitHub upload.

## How to inspect the delivered prototype

Open **http://127.0.0.1:5173**, register an account and follow the learning route.
The helper-started application and its database are left running locally for
inspection. Stop them safely with `bash scripts/stop-dev.sh`, without deleting
data. Restart with `bash scripts/dev.sh`. Private credentials and persistent
development PostgreSQL files live under Git-ignored `.local/`, not localStorage.

See README for explicit owner-only promotion of your own account, and
`docs/USER_GUIDE.md` for user instructions and prioritized extensions.

No production migrations, deployments, paid resources, commits or pushes were
performed. Candidate commit after reviewing all untracked files:
`feat: deliver persistent EnglishQuest learning prototype`.
