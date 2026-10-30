import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, expect, it, vi } from 'vitest'
import Diploma from './Diploma.jsx'
import DiplomaDocument from '../components/DiplomaDocument.jsx'
import useResource from '../lib/useResource.js'
import { api } from '../lib/api.js'

vi.mock('../lib/useResource.js', () => ({ default: vi.fn() }))
vi.mock('../lib/api.js', () => ({ api: { post: vi.fn() }, errorMessage: error => error.message }))
const record = { id: 'diploma-private-id', recipient_name: 'Ana María López', issued_at: '2026-10-30T12:00:00Z', course_title: 'Ruta guiada de práctica A1–C2', activity_count: 203 }
const view = () => render(<MemoryRouter><Diploma /></MemoryRouter>)

beforeEach(() => {
  vi.clearAllMocks()
  useResource.mockReturnValue({ loading: false, error: '', data: { eligible: true, completed: 203, total: 203, diploma: null } })
})

it('does not request a diploma name or show an issue button before completion', () => {
  useResource.mockReturnValue({ loading: false, error: '', data: { eligible: false, completed: 3, total: 203, diploma: null } })
  view()
  expect(screen.getByText(/Has aprobado 3 de 203/)).toBeInTheDocument()
  expect(screen.queryByLabelText('Nombre para el diploma')).not.toBeInTheDocument()
  expect(screen.queryByRole('button', { name: 'Confirmar nombre y emitir diploma' })).not.toBeInTheDocument()
})
it('asks for the chosen name at the end and prints only an issued document', async () => {
  api.post.mockResolvedValue({ data: { data: record } })
  const print = vi.spyOn(window, 'print').mockImplementation(() => {})
  view()
  expect(screen.getByLabelText('Nombre para el diploma')).toHaveValue('')
  fireEvent.change(screen.getByLabelText('Nombre para el diploma'), { target: { value: 'Ana María López' } })
  fireEvent.click(screen.getByRole('button', { name: 'Confirmar nombre y emitir diploma' }))
  await waitFor(() => expect(screen.getByText('Ana María López')).toBeInTheDocument())
  expect(api.post).toHaveBeenCalledWith('/diploma', { name: 'Ana María López' })
  fireEvent.click(screen.getByRole('button', { name: 'Imprimir o guardar PDF' }))
  expect(print).toHaveBeenCalledTimes(1)
  print.mockRestore()
})
it('does not fabricate a diploma when issuance fails', async () => {
  api.post.mockRejectedValue(new Error('Primero completa las actividades.'))
  view()
  fireEvent.change(screen.getByLabelText('Nombre para el diploma'), { target: { value: 'Ana María López' } })
  fireEvent.click(screen.getByRole('button', { name: 'Confirmar nombre y emitir diploma' }))
  await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('Primero completa las actividades.'))
  expect(screen.queryByRole('article', { name: 'Diploma simbólico de finalización' })).not.toBeInTheDocument()
})
it('reopens the immutable issued diploma instead of asking for a new name', () => {
  useResource.mockReturnValue({ loading: false, error: '', data: { eligible: true, completed: 203, total: 203, diploma: record } })
  view()
  expect(screen.getByText('Ana María López')).toBeInTheDocument()
  expect(screen.queryByLabelText('Nombre para el diploma')).not.toBeInTheDocument()
  expect(screen.getByText(/No es un título oficial/)).toBeInTheDocument()
})
it('renders recipient text without interpreting markup', () => {
  const { container } = render(<DiplomaDocument diploma={{ ...record, recipient_name: '<img src=x onerror=alert(1)>' }} />)
  expect(screen.getByText('<img src=x onerror=alert(1)>')).toBeInTheDocument()
  expect(container.querySelector('img')).toBeNull()
})
