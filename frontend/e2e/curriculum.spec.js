import { test, expect } from '@playwright/test'
import { unlockByExam } from './helpers/access.js'

test('B1 to C2: real first-unit reading, grammar, interaction, audio and review', async ({ page, context }, testInfo) => {
  test.setTimeout(200000)
  const api = 'http://127.0.0.1:8080/api/v1'
  const headers = { Origin: 'http://127.0.0.1:5173' }
  const account = { name: 'Advanced learner', email: `curriculum-${Date.now()}@example.org`, password: 'Curriculum-browser-test-2026!', timezone: 'UTC' }
  expect((await context.request.post(`${api}/auth/register`, { headers, data: account })).status()).toBe(201)
  expect((await context.request.post(`${api}/auth/login`, { headers, data: { email: account.email, password: account.password } })).status()).toBe(200)
  const units = [
    { level: 'B1', reading: ['The neighbour locked herself out.', 'Calling the daughter with a spare key'], grammar: ['was walking', 'had left'], interaction: ['I am glad we sorted it out; I can let my family know.'], audio: ['Because of the fire alarm', 'It returned to the depot.'] },
    { level: 'B2', reading: ['Use recordings as a supplement while monitoring effects.', 'All long-term effects'], grammar: ['Although participation varied, attendance remained stable.', 'whereas'], interaction: ['The pilot found no overall decline, but that does not settle every long-term question.'], audio: ['Qualified support', 'Testing different departmental approaches'] },
    { level: 'C1', reading: ['Causation and generalisability', 'They equate a limited causal claim with rejection of flexible work.'], grammar: ['Rarely have we seen such consistent findings.', 'with'], interaction: ['The revision changes a possible contribution into an unsupported guarantee.'], audio: ['Accepting a pattern while limiting generalisation', 'Replication in less predictable settings'] },
    { level: 'C2', reading: ['That failure is already established', 'The failure claim still needs an argued criterion.'], grammar: ['What we need is an explicit criterion.', 'whether'], interaction: ['Let us specify the criteria: some targets were missed, others met, and the main outcome is unmeasured.'], audio: ['A label substituting for justification', 'A different story about deservingness'] },
  ]
  for (const unit of units) {
    if (unit.level !== 'B1') await unlockByExam(page, unit.level)
    await page.goto(`/mapa?level=${unit.level}`)
    await expect(page.locator('.map-unit')).toHaveCount(6)
    const firstUnit = page.locator('.map-unit').first()
    await expect(firstUnit.locator('.path-step')).toHaveCount(5)
    const titles = ['Read and discover', 'Grammar in context', 'Respond with purpose', 'Listen for meaning', 'Unit review']
    for (const [step, title] of titles.entries()) {
      if (step) await page.goto(`/mapa?level=${unit.level}`)
      await page.locator('.map-unit').first().getByRole('link', { name: `Abrir ${title}`, exact: true }).click()
      await page.getByRole('button', { name: 'Empezar desafío', exact: true }).click()
      await expect(page.locator('.question-card').first()).toBeVisible()
      const answers = step === 0 ? unit.reading : step === 1 ? unit.grammar : step === 2 ? unit.interaction : step === 3 ? unit.audio : [...unit.reading, ...unit.interaction]
      for (const [index, answer] of answers.entries()) {
        const question = page.locator('.question-card').nth(index)
        if (step === 1) await question.getByRole('textbox').fill(answer)
        else {
          if (step === 3) {
            const audio = question.locator('audio')
            await audio.evaluate(element => element.play())
            await expect.poll(() => audio.evaluate(element => element.ended), { timeout: 25000 }).toBe(true)
          }
          await question.getByRole('button', { name: answer, exact: true }).click()
        }
      }
      await page.getByRole('button', { name: 'Comprobar respuestas', exact: true }).click()
      await expect(page.locator('.result-score')).toHaveText('100/100')
    }
    await page.goto(`/mapa?level=${unit.level}`)
    await expect(page.locator('.map-unit').first().locator('.path-step.done')).toHaveCount(5)
    await expect(page.locator('.map-unit').nth(1).getByRole('link', { name: 'Abrir Read and discover', exact: true })).toBeVisible()
  }
  await page.reload()
  await expect(page.locator('.map-unit').first().locator('.path-step.done')).toHaveCount(5)
  await page.locator('.map-unit').first().getByRole('link', { name: 'Ver unidad' }).click()
  await page.getByLabel('Tu borrador en inglés', { exact: true }).fill('The editorial treats failure as established, although its own evidence is mixed. An adequate evaluation requires explicit criteria rather than a loaded premise.')
  await page.getByLabel('Reflexión: qué revisaste y qué dudas tienes', { exact: true }).fill('I distinguished a presupposition from an evidence-based conclusion.')
  await page.getByRole('button', { name: 'Guardar proyecto', exact: true }).click()
  await expect(page.getByText('Proyecto guardado. No se ha asignado una calificación automática.')).toBeVisible()
  await page.reload()
  await expect(page.getByLabel('Tu borrador en inglés', { exact: true })).toHaveValue(/The editorial treats failure/)
  await page.screenshot({ path: testInfo.outputPath('c2-guided-route.png'), fullPage: true, animations: 'disabled' })
})
