import { test, expect } from '@playwright/test'

const API = 'http://127.0.0.1:8080/api/v1'
const origin = { Origin: 'http://127.0.0.1:5173' }
const wordMeanings = { apple: 'Manzana', banana: 'Plátano', orange: 'Naranja', book: 'Libro', pen: 'Bolígrafo', school: 'Escuela', breakfast: 'Desayuno', night: 'Noche' }
const legacySentences = {
  'Me gustan las manzanas.': ['I', 'like', 'apples.'],
  'Esto es un plátano.': ['This', 'is', 'a', 'banana.'],
  'La naranja es dulce.': ['The', 'orange', 'is', 'sweet.'],
  'Esto es un libro.': ['This', 'is', 'a', 'book.'],
  'Tengo un bolígrafo.': ['I', 'have', 'a', 'pen.'],
  'Voy a la escuela.': ['I', 'go', 'to', 'school.'],
}

test('seven additional games, real audio, server timer, booster and password lifecycle', async ({ page, context }, testInfo) => {
  test.setTimeout(90000)
  await page.emulateMedia({ reducedMotion: 'reduce' })
  const failures = []
  page.on('response', response => { if (response.url().startsWith(API) && response.status() >= 500) failures.push(response.status()) })
  const email = `extended-${Date.now()}@example.org`
  const password = 'Extended-learning-test-2026!'
  const registration = await context.request.post(`${API}/auth/register`, { headers: origin, data: { name: 'Extended Explorer', email, password, timezone: 'UTC' } })
  expect(registration.status()).toBe(201)
  expect((await context.request.post(`${API}/auth/login`, { headers: origin, data: { email, password } })).status()).toBe(200)
  // Setup completes prerequisite activities through real, authorized API calls.
  // Correct keys are NOT fetched from the database or an admin endpoint.
  const lessons = (await (await context.request.get(`${API}/lessons`)).json()).data
  for (const lesson of lessons.filter(l => l.level_code === 'A1').slice(0, 6)) {
    const attempt = (await (await context.request.post(`${API}/activities/${lesson.activity_id}/start`, { headers: origin })).json()).data
    const answers = {}
    for (const question of attempt.activity.questions) {
      const option = label => question.options.find(o => o.label === label).id
      if (lesson.kind === 'word_hunter') answers[question.id] = option(wordMeanings[question.prompt.replace('Find: ', '')])
      if (lesson.kind === 'sentence_builder') answers[question.id] = legacySentences[question.prompt].map(option)
      if (lesson.kind === 'memory_match') answers[question.id] = Object.entries(wordMeanings).filter(([en, es]) => question.options.some(o => o.label === en) && question.options.some(o => o.label === es)).map(([en, es]) => [option(en), option(es)])
    }
    const result = await context.request.post(`${API}/attempts/${attempt.id}/submit`, { headers: origin, data: { answers } })
    expect(result.status()).toBe(200)
    expect((await result.json()).data.score).toBe(100)
  }
  async function open(title, level = 'A1') {
    await page.goto(`/mapa?level=${level}`)
    await page.getByRole('link', { name: `Abrir ${title}`, exact: true }).first().click()
    await page.getByRole('button', { name: 'Empezar desafío', exact: true }).click()
    await expect(page.locator('.question-card').first()).toBeVisible()
  }
  async function submit() {
    await page.getByRole('button', { name: 'Comprobar respuestas', exact: true }).click()
    await expect(page.locator('.result-score')).toHaveText('100/100')
  }
  await open('Sound Detective')
  const listening = ['Blue', 'Two', 'On the desk']
  for (let i = 0; i < listening.length; i++) {
    const question = page.locator('.question-card').nth(i)
    const audio = question.locator('audio')
    await audio.evaluate(element => element.play())
    await expect.poll(() => audio.evaluate(element => element.ended), { timeout: 10000 }).toBe(true)
    await question.getByRole('button', { name: listening[i], exact: true }).click()
  }
  await submit()
  await open('Grammar Rescue')
  for (const [i, answer] of ['She goes to school.', 'I have a book.', 'They are happy.'].entries()) await page.locator('.question-card').nth(i).getByRole('textbox').fill(answer)
  await submit()
  await open('Spell Master')
  for (const [i, word] of ['apple', 'school', 'banana'].entries()) {
    const bank = page.locator('.question-card').nth(i).locator('.word-bank')
    for (const letter of word) await bank.locator('button:not([disabled])').filter({ hasText: new RegExp(`^${letter}$`) }).first().click()
  }
  await submit()
  await open('Conversation Quest')
  for (const [i, answer] of ['I am fine, thank you.', 'My name is Sam.', 'You are welcome.'].entries()) await page.locator('.question-card').nth(i).getByRole('button', { name: answer, exact: true }).click()
  await submit()
  await open('Story Explorer')
  await expect(page.locator('.reading-passage').first()).toContainText('Milo is a student.')
  for (const [i, answer] of ['Blue', 'An apple', 'Ten minutes'].entries()) await page.locator('.question-card').nth(i).getByRole('button', { name: answer, exact: true }).click()
  await submit()
  await open('Word Rain')
  const attemptPath = new URL(page.url()).pathname.split('/').pop()
  const before = (await (await context.request.get(`${API}/attempts/${attemptPath}`)).json()).data.expires_at
  await page.reload()
  await expect(page.getByRole('timer')).toBeVisible()
  const after = (await (await context.request.get(`${API}/attempts/${attemptPath}`)).json()).data.expires_at
  expect(after).toBe(before)
  for (const word of ['apple', 'banana', 'orange']) await page.locator('.rain-field').getByRole('button', { name: word, exact: true }).click()
  await submit()
  await open('English Rush')
  await page.locator('.question-card').nth(0).getByRole('button', { name: 'children', exact: true }).click()
  for (const word of ['I', 'am', 'reading.']) await page.locator('.question-card').nth(1).locator('.word-bank').getByRole('button', { name: word, exact: true }).click()
  await page.locator('.question-card').nth(2).getByRole('textbox').fill('goes')
  await submit()
  await page.getByRole('link', { name: 'Logros', exact: true }).click()
  await expect(page.locator('.achievement:not(.earned)').filter({ hasText: 'Ruta de nivel completada' })).toBeVisible()
  await page.getByRole('link', { name: 'Misiones', exact: true }).click()
  await expect(page.getByRole('button', { name: 'Activar doble XP', exact: true })).toBeVisible()
  await page.getByRole('button', { name: 'Activar doble XP', exact: true }).click()
  await expect(page.getByText(/Activo ·/)).toBeVisible()
  await page.reload()
  await expect(page.getByText(/Activo ·/)).toBeVisible()
  await page.screenshot({ path: testInfo.outputPath('missions-and-booster.png'), fullPage: true, animations: 'disabled' })
  await open('Discover words', 'A2')
  for (const [word, meaning] of [['breakfast', 'Desayuno'], ['school', 'Escuela'], ['night', 'Noche']]) {
    await page.locator('.question-card').filter({ has: page.getByRole('heading', { name: `Find: ${word}`, exact: true }) }).getByRole('button', { name: meaning, exact: true }).click()
  }
  await submit()
  await expect(page.getByText(/Potenciador aplicado/)).toBeVisible()
  await expect(page.locator('.xp-result')).toContainText('+50 XP')
  await page.getByRole('link', { name: 'Mi perfil', exact: true }).click()
  await page.getByLabel('Contraseña actual', { exact: true }).fill(password)
  await page.getByLabel('Nueva contraseña', { exact: true }).fill('Updated-learning-test-2026!')
  await page.getByLabel('Confirmar nueva contraseña', { exact: true }).fill('Updated-learning-test-2026!')
  await page.getByRole('button', { name: 'Actualizar contraseña y cerrar sesiones' }).click()
  await expect(page.getByRole('heading', { name: '¡Qué bueno verte!' })).toBeVisible()
  expect((await context.request.post(`${API}/auth/login`, { headers: origin, data: { email, password } })).status()).toBe(401)
  expect((await context.request.post(`${API}/auth/login`, { headers: origin, data: { email, password: 'Updated-learning-test-2026!' } })).status()).toBe(200)
  expect(failures).toEqual([])
})
