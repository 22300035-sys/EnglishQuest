import { test, expect } from '@playwright/test'
import { env } from 'node:process'

test('authorized administrator authors and edits persisted content', async ({
  page,
}, testInfo) => {
  test.skip(
    !env.E2E_ADMIN_EMAIL || !env.E2E_ADMIN_PASSWORD,
    'Supply an explicitly authorized local admin account; never a production account.',
  )
  const serverFailures = []
  page.on('response', (response) => {
    if (response.url().includes(':8080/api/v1/') && response.status() >= 500)
      serverFailures.push(response.status())
  })
  const title = `Browser authored lesson ${Date.now()}`
  await page.goto('/login')
  await page.getByLabel('Correo electrónico').fill(env.E2E_ADMIN_EMAIL)
  await page
    .getByLabel('Contraseña', { exact: true })
    .fill(env.E2E_ADMIN_PASSWORD)
  await page.getByRole('button', { name: 'Entrar a mi aventura' }).click()
  await page.getByRole('link', { name: 'Administración', exact: true }).click()
  await expect(
    page.getByRole('heading', { name: 'Panel de administración' }),
  ).toBeVisible()
  await page.getByRole('button', { name: 'Nueva lección' }).click()
  await page.getByLabel('Título de la lección').fill(title)
  await page
    .getByLabel('Objetivo de aprendizaje')
    .fill('Recognize everyday fruit words.')
  await page
    .getByLabel('Microlección (en inglés)')
    .fill('An apple is a fruit. Read the word and look at the picture.')
  await page.getByLabel('Publicada', { exact: true }).uncheck()
  await page
    .getByRole('button', { name: 'Guardar contenido', exact: true })
    .click()
  await expect(
    page.getByText(
      'Contenido guardado. Los intentos en curso conservan su versión original.',
    ),
  ).toBeVisible()
  const row = page.locator('.admin-list > div').filter({ hasText: title })
  await expect(row).toBeVisible()
  await row.getByRole('button', { name: 'Editar', exact: true }).click()
  await expect(page.getByLabel('Microlección (en inglés)')).toHaveValue(
    'An apple is a fruit. Read the word and look at the picture.',
  )
  await page
    .getByLabel('Microlección (en inglés)')
    .fill(
      'An apple is a fruit. Apples can be red or green. Practice the word apple.',
    )
  await page.getByLabel('Enunciado', { exact: true }).fill('Find: banana')
  await page.getByLabel('Explicación después de responder', { exact: true }).fill('Banana means plátano.')
  await page.getByLabel('Opción correcta', { exact: true }).selectOption({ label: 'Plátano' })
  await page
    .getByRole('button', { name: 'Guardar contenido', exact: true })
    .click()
  await expect(
    page.getByText(
      'Contenido guardado. Los intentos en curso conservan su versión original.',
    ),
  ).toBeVisible()
  await page.reload()
  await page
    .locator('.admin-list > div')
    .filter({ hasText: title })
    .getByRole('button', { name: 'Editar', exact: true })
    .click()
  await expect(page.getByLabel('Microlección (en inglés)')).toHaveValue(
    'An apple is a fruit. Apples can be red or green. Practice the word apple.',
  )
  await expect(page.getByLabel('Enunciado', { exact: true })).toHaveValue('Find: banana')
  await expect(page.getByLabel('Opción correcta', { exact: true })).toHaveValue('b')
  await page.screenshot({
    path: testInfo.outputPath('admin-editor.png'),
    fullPage: true,
    animations: 'disabled',
  })
  await page.locator('.admin-editor').getByRole('button', { name: 'Cancelar', exact: true }).click()
  const unitTitle = `Browser unit ${Date.now()}`
  await page.getByRole('button', { name: 'Nueva unidad', exact: true }).click()
  await page.getByLabel('Título del catálogo', { exact: true }).fill(unitTitle)
  await page.getByLabel('Descripción', { exact: true }).fill('A private draft unit for browser verification.')
  await page.getByRole('button', { name: 'Guardar catálogo', exact: true }).click()
  await expect(page.locator('.admin-list > div').filter({ hasText: unitTitle })).toBeVisible()
  const key = `browser_image_${Date.now()}`
  await page.getByRole('button', { name: 'Nuevo recurso', exact: true }).click()
  await page.getByLabel('Clave del recurso').fill(key)
  await page.getByLabel('Título del catálogo', { exact: true }).fill('Browser authorized illustration')
  await page.getByLabel('URL o ruta del archivo').fill('/images/apple.svg')
  await page.getByLabel('Autoría, licencia o atribución').fill('Original EnglishQuest illustration')
  await page.getByRole('button', { name: 'Guardar catálogo', exact: true }).click()
  await expect(page.locator('.admin-list > div').filter({ hasText: key })).toBeVisible()
  await page.getByRole('button', { name: 'Editar nivel B1', exact: true }).click()
  await page.getByLabel('Habilitar nivel', { exact: true }).check()
  await page.getByRole('button', { name: 'Guardar catálogo', exact: true }).click()
  await expect(page.getByText('Catálogo guardado. Los cambios se aplican a las nuevas sesiones de aprendizaje.')).toBeVisible()
  expect(serverFailures).toEqual([])
})
