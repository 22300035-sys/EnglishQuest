import { test, expect } from '@playwright/test'

test('real account → three games → saved rewards → reload → logout', async ({
  page,
}, testInfo) => {
  const serverFailures = []
  page.on('response', (response) => {
    if (response.url().includes(':8080/api/v1/') && response.status() >= 500)
      serverFailures.push(response.status())
  })
  const email = `browser-${Date.now()}@example.org`
  const password = 'Browser-learning-test-2026!'
  await page.goto('/')
  await expect(page.getByRole('heading', { name: /Tu inglés/ })).toBeVisible()
  await page.getByRole('link', { name: 'Comenzar aventura' }).click()
  await page.getByLabel('Tu nombre').fill('Browser Explorer')
  await page.getByLabel('Correo electrónico').fill(email)
  await page.getByLabel('Contraseña', { exact: true }).fill(password)
  await page.getByRole('button', { name: 'Crear mi cuenta' }).click()
  await expect(
    page.getByText('Cuenta creada. Inicia sesión para comenzar.'),
  ).toBeVisible()
  await page.getByLabel('Contraseña', { exact: true }).fill(password)
  await page.getByRole('button', { name: 'Entrar a mi aventura' }).click()
  await expect(
    page.getByRole('heading', { name: '¡Hola, Browser Explorer!' }),
  ).toBeVisible()
  await page.getByRole('link', { name: 'Continuar aprendiendo' }).click()
  await page.getByRole('button', { name: 'Empezar desafío' }).click()
  for (const [word, meaning] of [
    ['apple', 'Manzana'],
    ['banana', 'Plátano'],
    ['orange', 'Naranja'],
  ]) {
    const question = page
      .locator('.question-card')
      .filter({
        has: page.getByRole('heading', { name: `Find: ${word}`, exact: true }),
      })
    await question
      .locator('.picture-option')
      .filter({ hasText: meaning })
      .click()
  }
  await page.getByRole('button', { name: 'Comprobar respuestas' }).click()
  await expect(page.getByText('RESULTADO GUARDADO')).toBeVisible()
  await expect(page.locator('.result-score')).toHaveText('100/100')
  await page.reload()
  await expect(page.locator('.result-score')).toHaveText('100/100')

  await page.getByRole('link', { name: 'Continuar mi ruta' }).click()
  await page
    .getByRole('link', { name: 'Abrir Build a sentence', exact: true })
    .first()
    .click()
  await page.getByRole('button', { name: 'Empezar desafío' }).click()
  const sentences = [
    ['Me gustan las manzanas.', ['I', 'like', 'apples.']],
    ['Esto es un plátano.', ['This', 'is', 'a', 'banana.']],
    ['La naranja es dulce.', ['The', 'orange', 'is', 'sweet.']],
  ]
  for (const [prompt, words] of sentences) {
    const question = page
      .locator('.question-card')
      .filter({ has: page.getByRole('heading', { name: prompt, exact: true }) })
    for (const word of words)
      await question
        .locator('.word-bank')
        .getByRole('button', { name: word, exact: true })
        .click()
  }
  await page.getByRole('button', { name: 'Comprobar respuestas' }).click()
  await expect(page.locator('.result-score')).toHaveText('100/100')
  await page.getByRole('link', { name: 'Continuar mi ruta' }).click()
  await page
    .getByRole('link', { name: 'Abrir Unit checkpoint', exact: true })
    .first()
    .click()
  await page.getByRole('button', { name: 'Empezar desafío' }).click()
  // Learn the six card locations through the actual UI, then reset and pair.
  const cards = page.locator('.memory-card')
  await expect(cards).toHaveCount(6)
  const labels = []
  for (let i = 0; i < 6; i++) {
    await cards.nth(i).click()
    labels.push(await cards.nth(i).innerText())
  }
  await page.getByRole('button', { name: 'Reiniciar parejas' }).click()
  for (const pair of [
    ['apple', 'Manzana'],
    ['banana', 'Plátano'],
    ['orange', 'Naranja'],
  ]) {
    for (const label of pair) await cards.nth(labels.indexOf(label)).click()
  }
  await page.getByRole('button', { name: 'Comprobar respuestas' }).click()
  await expect(page.locator('.result-score')).toHaveText('100/100')
  await expect(page.locator('.xp-result')).toContainText('+100 XP')

  await page.getByRole('link', { name: 'Mi aventura', exact: true }).click()
  await expect(
    page.locator('.stat-card').filter({ hasText: 'XP acumulada' }),
  ).toContainText('170')
  await expect(page.locator('.streak-card h2')).toHaveText('1 día')
  await page.reload()
  await expect(
    page.locator('.stat-card').filter({ hasText: 'XP acumulada' }),
  ).toContainText('170')
  await page.screenshot({
    path: testInfo.outputPath('dashboard-desktop.png'),
    fullPage: true,
    animations: 'disabled',
  })
  await page.setViewportSize({ width: 390, height: 844 })
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= window.innerWidth,
    ),
  ).toBe(true)
  await page.screenshot({
    path: testInfo.outputPath('dashboard-mobile.png'),
    fullPage: true,
    animations: 'disabled',
  })
  await page.getByRole('link', { name: 'Misiones', exact: true }).click()
  await expect(
    page.getByText('¡Misión completada y recompensa guardada!').first(),
  ).toBeVisible()
  await page.getByRole('link', { name: 'Logros', exact: true }).click()
  await expect(page.locator('.achievement.earned')).toHaveCount(2)
  await page.getByRole('link', { name: 'Mi perfil', exact: true }).click()
  await page.getByLabel('Nombre', { exact: true }).fill('Updated Explorer')
  await page.getByRole('button', { name: 'Guardar cambios' }).click()
  await expect(page.getByText('Cambios guardados.')).toBeVisible()
  await page.reload()
  await expect(page.getByLabel('Nombre', { exact: true })).toHaveValue(
    'Updated Explorer',
  )
  await page.getByRole('button', { name: 'Cerrar sesión' }).click()
  await expect(
    page.getByRole('heading', { name: '¡Qué bueno verte!' }),
  ).toBeVisible()
  await page.goto('/dashboard')
  await expect(
    page.getByRole('heading', { name: '¡Qué bueno verte!' }),
  ).toBeVisible()
  expect(serverFailures).toEqual([])
})

test('all six guided levels are available without accreditation claims', async ({
  page,
}, testInfo) => {
  await page.goto('/niveles')
  await expect(page.locator('.level-card')).toHaveCount(6)
  await expect(page.locator('.level-card.unavailable')).toHaveCount(0)
  await expect(page.locator('.level-card').getByRole('link', { name: 'Explorar' })).toHaveCount(3)
  await expect(page.locator('.level-card').getByRole('link', { name: 'Ver cómo desbloquear' })).toHaveCount(3)
  await expect(page.getByText(/no equivale a una acreditación CEFR/)).toBeVisible()
  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/')
  await expect(page.getByRole('heading', { name: /Tu inglés/ })).toBeVisible()
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= window.innerWidth,
    ),
  ).toBe(true)
  await page.screenshot({
    path: testInfo.outputPath('welcome-mobile.png'),
    fullPage: true,
    animations: 'disabled',
  })
})
