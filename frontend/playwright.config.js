import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 45000,
  use: { baseURL: 'http://127.0.0.1:5173', headless: true, screenshot: 'only-on-failure' },
  // Start the real local API, PostgreSQL and Vite with scripts/dev.sh first.
  // No HTTP request is intercepted or mocked in these browser tests.
  workers: 1,
})
