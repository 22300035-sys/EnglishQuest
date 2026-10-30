import { test, expect } from '@playwright/test'
import { env } from 'node:process'

test('production PWA provides honest offline fallback without private caches', async ({ page, context }, testInfo) => {
  test.skip(!env.E2E_PWA_URL, 'Run a production preview and supply E2E_PWA_URL to exercise the service worker.')
  await page.goto(env.E2E_PWA_URL)
  const manifest = await page.locator('link[rel="manifest"]').getAttribute('href')
  expect(manifest).toBe('/manifest.webmanifest')
  await page.evaluate(async () => { await navigator.serviceWorker.ready })
  await page.reload()
  await expect.poll(() => page.evaluate(() => !!navigator.serviceWorker.controller)).toBe(true)
  const icons = await page.evaluate(async () => {
    const data = await (await fetch('/manifest.webmanifest')).json()
    return Promise.all(data.icons.map(async icon => ({ status: (await fetch(icon.src)).status, size: icon.sizes })))
  })
  expect(icons).toEqual([{ status: 200, size: '192x192' }, { status: 200, size: '512x512' }])
  await context.setOffline(true)
  await page.goto(`${env.E2E_PWA_URL}/dashboard`)
  await expect(page.getByRole('heading', { name: 'Estás sin conexión' })).toBeVisible()
  await expect(page.getByText(/no simula calificaciones ni XP/)).toBeVisible()
  const keys = await page.evaluate(async () => {
    const names = await caches.keys()
    const groups = await Promise.all(names.map(async name => (await (await caches.open(name)).keys()).map(request => new URL(request.url).pathname)))
    return groups.flat()
  })
  expect(keys.every(path => !path.startsWith('/api/') && !['/dashboard', '/perfil', '/login'].includes(path))).toBe(true)
  await page.screenshot({ path: testInfo.outputPath('offline-fallback.png'), fullPage: true })
  await context.setOffline(false)
})
