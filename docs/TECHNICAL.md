# EnglishQuest — Technical and academic specification

## Executive Summary

EnglishQuest is a working educational prototype for learning introductory English
through micro-lessons and ten game types. Accounts, grading, results, rewards and
progress persist in PostgreSQL. The catalogue contains 203 activities
across 41 units: 43 in A1, 40 in A2 and 30 each in B1, B2, C1 and C2.
The new 36 units also offer private writing projects and self-review rubrics.
This is not an academically validated complete CEFR course,
a proficiency assessment or a certified language-learning product.

## Target Audience

Spanish-speaking beginners practicing everyday English independently. Educators
can edit lessons and activities through an authorized administration workspace.
Age-specific privacy, parental consent and public-service terms require review
before broad deployment; no claim of regulatory certification is made.

## Pedagogical Framework

Implemented: short objective-led lessons, vocabulary–picture association, sentence
ordering and retrieval through word/meaning matching. Students read instruction
and explanation, receive corrective feedback and
may repeat without losing progress. The next step unlocks after passing.

Gamification supports practice: repeat practice can improve a best score or
maintain a streak, but cannot farm activity XP. The sample has not undergone
independent CEFR or educational-efficacy validation. Spaced repetition and
adaptive difficulty remain future work.

## Linguistic Scope

- A1 Fruit garden: apple, banana, orange; basic statements and identification.
- A1 My classroom: book, pen, school; everyday subject–verb–object sentences.
- Initial A2 Everyday routines: breakfast, school, night; present simple,
  third-person -s/-es and time expressions.
- Each unit has word association, sentence ordering and a matching checkpoint.
- Fourteen additional lessons cover listening (six bundled synthetic clips),
  target-sentence correction, spelling, contextual conversation, short stories,
  category recognition and mixed checkpoints. Timed games enforce server deadlines.
- Written answers check a normalized target response, not arbitrary free-form
  production. Pronunciation recognition and independent writing assessment are absent.
- Six guided routes A1–C2 are published, with six new units per level. Advanced
  tasks address argument, evidence, register, stance, irony, ambiguity and faithful
  synthesis. See `CURRICULUM.md` for all topics and the limits of the CEFR labels.
- Productive writing projects save drafts and reflections privately in PostgreSQL.
  They do not receive automatic grades, XP or CEFR credentials. Oral rehearsal
  instructions do not record or automatically assess speech.

## Architecture

- React/JSX + Vite + Tailwind 4, Router, Axios and Lucide. Original CSS adventure
  landscape and eight local SVG illustrations; no paid or external image service.
- Go/Gin API `/api/v1`, handlers, auth/grading/validation services, repositories,
domain models, safe JSON errors and restricted credentialed CORS.
- PostgreSQL/GORM repositories, parameterized SQL, explicit embedded migrations,
  database constraints, bounded connection pool, timeouts and graceful shutdown.
- Passwords: bcrypt with 10–72-byte input bounds. Sessions: 256-bit opaque random
  values, only SHA-256 hashes stored. Cookies: HttpOnly; production adds Secure,
  SameSite=None, Partitioned and the `__Host-` prefix. JWT fields in configuration
  are retained for compatibility; the application uses revocable opaque sessions,
  not signed browser JWTs. `JWT_REFRESH_TTL` sets absolute session lifetime.
- Each state-changing browser request must have an exact trusted Origin, including
  login. Authorization loads the current active user and role from PostgreSQL.
  No account/session/progress data is stored in localStorage.
- Auth rate limits are bounded and process-local; distributed rate limiting and
  deployment-specific trusted proxy configuration remain hardening tasks.

## Data Flow

Registration → PostgreSQL user/password hash → login → hashed session and HttpOnly
cookie → level/map selection → lesson → server-created attempt → interactive game
→ submitted answer IDs → Go grading → transactional result/progress/rewards
→ persistent result screen → refreshed dashboard/statistics.

Questions and answer options are relational. Each attempt stores a private,
immutable activity snapshot and opaque randomized option IDs. Neither solutions
nor explanations are returned before completion. Admin edits do not alter old
or running attempts. Only the attempt owner can submit/read a result; privileged
content access is separately ADMIN-authorized.

## Scoring and Gamification Rules

### Level access and symbolic recognition

A1/A2/B1 are initially open. Advanced access is granted by completion of the
previous published route or a private server-graded access exam. An exam samples
12 questions from 16 per advanced level, lasts 20 minutes and requires 80%
(10 correct); retakes wait 10 minutes after submission. Access exams award no
activity completion, XP or streak. Previous-level access is required to take the
next exam. Existing advanced attempts are grandfathered by migration 9.

All activity starts and unit project reads/writes enforce level access on the
server. Course completion is separate from access. After all published activities
across six available levels are passed, the user chooses a diploma name. Issuance
is owner-only, unique and durable, with an immutable snapshot of completed content.
The browser prints an A4 landscape symbolic diploma; no external PDF service or
public name lookup is used. This is not CEFR accreditation or identity verification.

### Activity rewards

- Practice passes at 60%; unit checkpoints pass at 80%. Memory currently has one
  question containing all pairs, so it passes only when every pair is correct.
- Grade is the proportion of completely correct questions. 1/2/3 stars represent
  passing/at least 80%/100%. Partial sentence/pair credit is not implemented.
- First passed activity: 20 XP; first passed checkpoint: 40 XP. A perfect first
  completion adds 10 XP. No repeat activity or request-replay XP.
- Daily mission: first-pass three different activities in the account's local
  civil day; 30 XP once. Assignments/progress/reward records are stored in SQL.
- Additional missions reward first completions by skill and perfect practice.
  The three-activity mission grants one 15-minute booster per local day. Activation
  is owner-only, single-use and serialized; only one booster may run at a time.
  Active boosters double base activity XP at submission, never mission/perfect bonuses.
- Achievements: first activity, ten different activities, first checkpoint,
  seven consecutive days, 100 target-word practices (not necessarily distinct words),
  and completing all published activities in a level (not CEFR certification).
- Streak: at most one increment per civil date after valid passed practice,
  including a fresh repeat attempt. Duplicate submission is not fresh practice.
  IANA timezone is fixed on registration. DST and year boundaries are handled
  using calendar dates, not elapsed 24-hour durations. No pets or streak avatars.
- A transaction locks the user and attempt, writes results, applies progress and
  rewards, and commits together. Unique source keys and a partial open-attempt
  index protect repeated/concurrent requests.

## User Stories and Acceptance Criteria

### US1 — Create and protect my account (implemented)
As a learner, I want a private account so that my practice persists across sessions.
- **Given** valid unique details, **When** I register, **Then** a STUDENT account
  with a bcrypt hash and validated timezone is saved, never a plaintext password.
- **Given** an expired, revoked or absent session, **When** I request protected
  data, **Then** Go denies access; a forged React role does not grant admin access.

### US2 — Complete a learning challenge (implemented)
As an A1 learner, I want feedback and saved results so that I can practice reliably.
- **Given** an unlocked activity, **When** I start, **Then** I receive randomized
  public questions without solutions, backed by a persistent attempt snapshot.
- **Given** completed answers, **When** I submit, **Then** Go grades the snapshot,
  saves feedback/progress and unlocks the next lesson only if I pass.

### US3 — Earn trustworthy rewards (implemented)
As a learner, I want fair XP and streaks so that progress reflects real practice.
- **Given** an already submitted attempt, **When** requests are replayed or sent
  concurrently, **Then** the saved result is returned without duplicate XP.
- **Given** a valid fresh attempt on the same local day, **When** I pass again,
  **Then** my streak grows no more than once for that civil date.

### US4 — Manage teaching content (implemented, basic scope)
As an administrator, I want to create/edit lessons without changing source code.
- **Given** a database-authorized ADMIN, **When** valid lesson/activity/question
  content is saved, **Then** all related edits commit together and are audited.
- **Given** a started attempt, **When** an admin edits the source activity,
  **Then** that attempt's snapshot and grading remain unchanged.

## UI/UX

Responsive adventure-themed welcome, original landscape, strong focus outlines,
moderate reduced-motion-aware CSS transitions, loading/error/retry states, real
forms and a persistent application navigation shell. Fourteen screen categories
are represented: welcome, registration, login, dashboard, map, unit, lesson,
exercise, results, achievements, missions, statistics, profile and admin.

The admin question editor provides visual authoring for all ten types with a
validated advanced JSON alternative. Level, unit and image/audio resource catalogue
editing is available; it does not upload binary media or assign roles. The map uses original nodes/illustration,
not copyrighted Candy Crush characters or assets. Locked/empty/future content is
explicitly labeled. Session checks can redirect to login after expiration.

## Testing

Unit/component tests cover grading, streak calendar rules, content validation,
configuration, route protection, forms, game interactions, errors and navigation.
Real PostgreSQL integration covers registration/login, all ten games, concurrent
submissions, XP anti-farming, daily mission, roles, owner isolation, snapshot
stability, admin edits/creation, password/session revocation, account disabling,
deadline expiry and booster ownership/concurrency/expiry.

Playwright tests target the real locally running app without network interception:
signup → login → three games → saved results → dashboard/reload → profile/logout,
plus unavailable levels and mobile-width overflow checks. Execution status and
limitations belong in the increment report; having a test file is not a pass.

## Conclusions and Extensions

The prototype supports a genuine learning flow with durable PostgreSQL data and
server-controlled grading/rewards. It is a demonstration foundation, not the full
long-term product. Highest-value extensions: academically reviewed complete courses,
safe email password recovery, independent writing/pronunciation assessment,
spaced repetition and broader accessibility testing. A production-only service worker
caches public assets and provides an offline explanation; it never caches private
routes/API responses or validates offline XP. Render deployment remains unperformed.
