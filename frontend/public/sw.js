/* Public assets only. NEVER cache API calls, sessions, forms or private routes. */
const CACHE = 'englishquest-public-v1'
const PUBLIC_FILES = ['/offline.html', '/favicon.svg', '/manifest.webmanifest', '/pwa/icon-192.png', '/pwa/icon-512.png']

self.addEventListener('install', event => {
  event.waitUntil(caches.open(CACHE).then(cache => cache.addAll(PUBLIC_FILES.map(url => new Request(url, { credentials: 'omit' })))).then(() => self.skipWaiting()))
})
self.addEventListener('activate', event => {
  event.waitUntil(caches.keys().then(keys => Promise.all(keys.filter(key => key.startsWith('englishquest-public-') && key !== CACHE).map(key => caches.delete(key)))).then(() => self.clients.claim()))
})
self.addEventListener('fetch', event => {
  const request = event.request
  const url = new URL(request.url)
  if (request.method !== 'GET' || url.origin !== self.location.origin || url.pathname.startsWith('/api/')) return
  if (request.mode === 'navigate') {
    event.respondWith(fetch(request).catch(async () => (await caches.match('/offline.html')) || Response.error()))
    return
  }
  const isPublic = PUBLIC_FILES.includes(url.pathname) || ['/assets/', '/images/', '/audio/', '/pwa/'].some(prefix => url.pathname.startsWith(prefix))
  if (!isPublic) return
  const publicRequest = new Request(request, { credentials: 'omit' })
  event.respondWith(fetch(publicRequest).then(response => {
    if (response.ok && response.type === 'basic' && !response.headers.get('Cache-Control')?.includes('no-store')) {
      const copy = response.clone()
      event.waitUntil(caches.open(CACHE).then(cache => cache.put(publicRequest, copy)))
    }
    return response
  }).catch(async () => (await caches.match(request)) || Response.error()))
})
