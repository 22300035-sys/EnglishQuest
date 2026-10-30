import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import App from './App.jsx'
import { AuthContext } from './auth/context.js'
import { api, getServiceStatus } from './lib/api.js'

vi.mock('./lib/api.js', () => ({
  api: { get: vi.fn(), post: vi.fn(), patch: vi.fn() },
  errorMessage: e => e.message,
  getServiceStatus: vi.fn(),
}))
afterEach(() => vi.resetAllMocks())

function view(path = '/', overrides = {}) {
  const auth = { user: null, loading: false, login: vi.fn(), logout: vi.fn(), ...overrides }
  render(<AuthContext.Provider value={auth}><MemoryRouter initialEntries={[path]}><App /></MemoryRouter></AuthContext.Provider>)
  return auth
}

describe('Navigation and availability', () => {
  it('has a real registration link and an honest scope notice', async () => {
    view()
    expect(screen.getByRole('heading', { name: /Tu inglés/ })).toBeInTheDocument()
    expect(screen.getByText(/No equivale a una acreditación CEFR/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('link', { name: /Comenzar aventura/ }))
    expect(await screen.findByRole('heading', { name: 'Tu aventura empieza aquí' })).toBeInTheDocument()
  })
  it('protects student routes without a session', () => {
    view('/dashboard')
    expect(screen.getByRole('heading', { name: '¡Qué bueno verte!' })).toBeInTheDocument()
  })
  it('shows loading before API completion', () => {
    getServiceStatus.mockImplementation(() => new Promise(() => {}))
    view('/estado')
    expect(screen.getByRole('status')).toHaveTextContent('Comprobando conexión')
    expect(screen.getByRole('button')).toBeDisabled()
  })
  it('reports successful API and database verification', async () => {
    getServiceStatus.mockResolvedValue({ api: 'ok', database: 'ready' })
    view('/estado')
    expect(await screen.findByText('API disponible · PostgreSQL conectado')).toBeInTheDocument()
  })
  it('shows connection errors and supports retry', async () => {
    getServiceStatus.mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce({ api: 'ok', database: 'ready' })
    view('/estado')
    expect(await screen.findByRole('alert')).toHaveTextContent('No se pudo verificar')
    fireEvent.click(screen.getByRole('button', { name: 'Volver a comprobar' }))
    await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('PostgreSQL conectado'))
    expect(getServiceStatus).toHaveBeenCalledTimes(2)
  })
  it('handles unknown routes', () => {
    view('/missing')
    expect(screen.getByRole('heading', { name: 'Página no encontrada' })).toBeInTheDocument()
  })
  it('does not render the admin editor for a student', () => {
    view('/admin', { user: { role: 'STUDENT', first_name: 'Test' } })
    expect(screen.getByRole('heading', { name: 'Acceso restringido' })).toBeInTheDocument()
    expect(api.get).not.toHaveBeenCalled()
  })
})

describe('Authentication forms', () => {
  it('rejects short passwords before registration', async () => {
    view('/registro')
    fireEvent.change(screen.getByLabelText('Tu nombre'), { target: { value: 'Alex' } })
    fireEvent.change(screen.getByLabelText('Correo electrónico'), { target: { value: 'alex@example.org' } })
    fireEvent.change(screen.getByLabelText('Contraseña'), { target: { value: 'short' } })
    fireEvent.submit(screen.getByRole('button', { name: 'Crear mi cuenta' }).closest('form'))
    expect(await screen.findByRole('alert')).toHaveTextContent('10 y 72 bytes')
    expect(api.post).not.toHaveBeenCalled()
  })
  it('registers and navigates to login without storing a password', async () => {
    api.post.mockResolvedValue({ data: { data: { id: 'test' } } })
    view('/registro')
    fireEvent.change(screen.getByLabelText('Tu nombre'), { target: { value: 'Alex' } })
    fireEvent.change(screen.getByLabelText('Correo electrónico'), { target: { value: 'alex@example.org' } })
    fireEvent.change(screen.getByLabelText('Contraseña'), { target: { value: 'test-password-42' } })
    fireEvent.submit(screen.getByRole('button', { name: 'Crear mi cuenta' }).closest('form'))
    expect(await screen.findByText('Cuenta creada. Inicia sesión para comenzar.')).toBeInTheDocument()
    expect(api.post).toHaveBeenCalledWith('/auth/register', expect.objectContaining({ email: 'alex@example.org', timezone: expect.any(String) }))
    expect(screen.getByLabelText('Contraseña')).toHaveValue('')
  })
  it('shows server login errors', async () => {
    const auth = view('/login')
    auth.login.mockRejectedValue(new Error('Credenciales incorrectas'))
    fireEvent.change(screen.getByLabelText('Correo electrónico'), { target: { value: 'alex@example.org' } })
    fireEvent.change(screen.getByLabelText('Contraseña'), { target: { value: 'wrong-password' } })
    fireEvent.submit(screen.getByRole('button', { name: 'Entrar a mi aventura' }).closest('form'))
    expect(await screen.findByRole('alert')).toHaveTextContent('Credenciales incorrectas')
  })
})
