import axios from 'axios'

// Environment variables exposed by Vite are public: never place secrets here.
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  timeout: 8000,
  withCredentials: true,
  headers: { Accept: 'application/json' },
})

export const errorMessage = (error) =>
  error.response?.data?.error ||
  (error.code === 'ECONNABORTED'
    ? 'La conexión tardó demasiado. Inténtalo otra vez.'
    : 'No se pudo conectar con el servidor. Comprueba tu conexión.')

export async function getServiceStatus(signal) {
  const [health, ready] = await Promise.all([
    api.get('/health', { signal }),
    api.get('/ready', { signal }),
  ])
  if (
    health.data?.data?.status !== 'ok' ||
    ready.data?.data?.status !== 'ready'
  ) {
    throw new Error('Unexpected service response')
  }
  return { api: 'ok', database: 'ready' }
}
