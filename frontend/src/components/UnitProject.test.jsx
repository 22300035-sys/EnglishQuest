import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import UnitProject from './UnitProject.jsx'
import { api } from '../lib/api.js'

vi.mock('../lib/useResource.js', () => ({ default: () => ({
  loading: false, error: '', data: { prompt: 'Reformulate the argument for a new audience.', rubric: 'Check scope and evidence.', body: 'My saved draft.', reflection: 'I checked the premise.' },
}) }))
vi.mock('../lib/api.js', () => ({ api: { put: vi.fn() }, errorMessage: error => error.message }))

it('loads the private draft and only reports success after the API saves it', async () => {
  api.put.mockResolvedValue({ data: { data: { saved: true } } })
  render(<UnitProject unitID="unit-id" />)
  expect(screen.getByLabelText('Tu borrador en inglés')).toHaveValue('My saved draft.')
  fireEvent.change(screen.getByLabelText('Tu borrador en inglés'), { target: { value: 'My revised argument.' } })
  fireEvent.click(screen.getByRole('button', { name: 'Guardar proyecto' }))
  await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('No se ha asignado una calificación automática.'))
  expect(api.put).toHaveBeenCalledWith('/units/unit-id/project', { body: 'My revised argument.', reflection: 'I checked the premise.' })
})

it('does not pretend a failed project save succeeded', async () => {
  api.put.mockRejectedValue(new Error('Sin conexión al servidor.'))
  render(<UnitProject unitID="unit-id" />)
  fireEvent.click(screen.getByRole('button', { name: 'Guardar proyecto' }))
  await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('Sin conexión al servidor.'))
  expect(screen.queryByText(/Proyecto guardado/)).not.toBeInTheDocument()
  expect(screen.getByLabelText('Tu borrador en inglés')).toHaveValue('My saved draft.')
})
