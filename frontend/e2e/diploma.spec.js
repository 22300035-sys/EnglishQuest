import { test, expect } from '@playwright/test'
import { expectedSubmission } from './helpers/course-answers.js'

test('all course activities → progressive unlocks → chosen name → durable printable diploma', async ({ page, context }, testInfo) => {
  test.setTimeout(120000)
  const api = 'http://127.0.0.1:8080/api/v1'
  const headers = { Origin: 'http://127.0.0.1:5173' }
  const account = { name: 'Course finisher', email: `diploma-${Date.now()}@example.org`, password: 'Diploma-completion-test-2026!', timezone: 'UTC' }
  expect((await context.request.post(`${api}/auth/register`, { headers, data: account })).status()).toBe(201)
  expect((await context.request.post(`${api}/auth/login`, { headers, data: { email: account.email, password: account.password } })).status()).toBe(200)
  const lessons = (await (await context.request.get(`${api}/lessons`)).json()).data
  expect(lessons.length).toBeGreaterThanOrEqual(203)
  let lastAttempt
  // Real course completion, not direct DB progress insertion or a test-only API.
  for (const lesson of lessons) {
    const response = await context.request.post(`${api}/activities/${lesson.activity_id}/start`, { headers })
    expect(response.status(), `${lesson.level_code} ${lesson.unit_title}`).toBe(201)
    const attempt = (await response.json()).data
    const result = await context.request.post(`${api}/attempts/${attempt.id}/submit`, { headers, data: expectedSubmission(lesson, attempt.activity) })
    expect(result.status()).toBe(200)
    expect((await result.json()).data.passed).toBe(true)
    lastAttempt = attempt.id
  }
  const levels = (await (await context.request.get(`${api}/learning/levels`)).json()).data
  expect(levels.filter(l => l.unlocked)).toHaveLength(6)
  expect(levels.filter(l => l.source === 'progress')).toHaveLength(3)
  await page.goto(`/resultados/${lastAttempt}`)
  await page.getByRole('link', { name: 'Elegir nombre y obtener mi diploma' }).click()
  await page.getByLabel('Nombre para el diploma').fill('María Fernanda O’Neill')
  await page.getByRole('button', { name: 'Confirmar nombre y emitir diploma' }).click()
  await expect(page.getByRole('article', { name: 'Diploma simbólico de finalización' })).toBeVisible()
  await expect(page.locator('.diploma-name')).toHaveText('María Fernanda O’Neill')
  await expect(page.getByText(/No es un título oficial/)).toBeVisible()
  const folio = await page.locator('.diploma-folio').textContent()
  await page.reload()
  await expect(page.locator('.diploma-folio')).toHaveText(folio)
  await expect(page.getByLabel('Nombre para el diploma')).toHaveCount(0)
  await page.screenshot({ path: testInfo.outputPath('symbolic-diploma.png'), fullPage: true, animations: 'disabled' })
  await page.emulateMedia({ media: 'print' })
  await expect(page.locator('.sidebar')).toBeHidden()
  await expect(page.getByRole('article', { name: 'Diploma simbólico de finalización' })).toBeVisible()
  const pdf = await page.pdf({ path: testInfo.outputPath('symbolic-diploma.pdf'), preferCSSPageSize: true, printBackground: true })
  expect(pdf.toString('latin1').match(/\/Type\s*\/Page\b/g)).toHaveLength(1)
})
