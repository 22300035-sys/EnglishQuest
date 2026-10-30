# EnglishQuest frontend

React, Vite, JavaScript/JSX, Tailwind CSS 4, React Router, Axios and Lucide.
See the root [README](../README.md) for the real Go/PostgreSQL development setup.

```bash
npm ci
npm run dev
npm run test
npm run lint
npm run build
```

`VITE_API_URL` is public and points to `/api/v1`. Never put database URLs or
private credentials in Vite variables. Authentication uses HttpOnly cookies,
not localStorage. Browser tests require the actual local API/database:

```bash
npx playwright install chromium --only-shell
npm run test:e2e
```

The optional admin browser test requires an explicitly authorized local
`E2E_ADMIN_EMAIL` and `E2E_ADMIN_PASSWORD` in the test process environment.
Without them it explicitly skips; no admin account is silently created/promoted.
