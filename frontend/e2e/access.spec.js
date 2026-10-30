import { test, expect } from '@playwright/test'

test('initial half open, failed exam cannot unlock or claim a diploma', async ({ page, context }) => {
  const api = 'http://127.0.0.1:8080/api/v1'
  const headers = { Origin: 'http://127.0.0.1:5173' }
  const account = { name: 'Progressive Explorer', email: `access-${Date.now()}@example.org`, password: 'Progressive-test-password-2026!', timezone: 'UTC' }
  expect((await context.request.post(`${api}/auth/register`, { headers, data: account })).status()).toBe(201)
  expect((await context.request.post(`${api}/auth/login`, { headers, data: { email: account.email, password: account.password } })).status()).toBe(200)
  await page.goto('/niveles')
  await expect(page.getByRole('link', { name: 'Explorar', exact: true })).toHaveCount(3)
  await expect(page.getByRole('link', { name: 'Ver cómo desbloquear', exact: true })).toHaveCount(3)
  await page.goto('/mapa?level=C2')
  await expect(page.getByRole('heading', { name: 'C2 todavía está bloqueado' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Examen de acceso a C2', exact: true })).toHaveCount(0)
  const lessons = (await (await context.request.get(`${api}/lessons`)).json()).data
  const blocked = lessons.find(l => l.level_code === 'B2')
  expect(blocked.level_locked).toBe(true)
  expect((await context.request.post(`${api}/activities/${blocked.activity_id}/start`, { headers })).status()).toBe(403)
  await page.goto('/mapa?level=B2')
  await page.getByRole('link', { name: 'Examen de acceso a B2', exact: true }).click()
  await page.getByRole('button', { name: 'Comenzar o retomar examen' }).click()
  await expect(page.locator('.question-card')).toHaveCount(12)
  await page.getByRole('button', { name: 'Enviar examen', exact: true }).click()
  await expect(page.getByRole('heading', { name: 'Todavía no alcanzas el umbral' })).toBeVisible()
  await page.getByRole('link', { name: 'Volver al acceso' }).click()
  await page.getByRole('button', { name: 'Comenzar o retomar examen' }).click()
  await expect(page.getByRole('alert')).toHaveText(/Espera 10 minutos/)
  await page.goto('/diploma')
  await expect(page.getByRole('heading', { name: 'Tu diploma te espera al terminar' })).toBeVisible()
  await expect(page.getByLabel('Nombre para el diploma')).toHaveCount(0)
  expect((await context.request.post(`${api}/diploma`, { headers, data: { name: 'Premature Diploma' } })).status()).toBe(409)
})
