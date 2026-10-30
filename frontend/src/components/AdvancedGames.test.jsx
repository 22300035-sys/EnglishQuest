import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import { ChoiceQuestion, SoundDetective, WrittenAnswer } from './AdvancedGames.jsx'
import BoosterInventory from './BoosterInventory.jsx'
import PWAControls from './PWAControls.jsx'
import { api } from '../lib/api.js'

vi.mock('../lib/api.js', () => ({ api: { post: vi.fn() }, errorMessage: error => error.message }))
const question = { options: [{ id: 'opaque-one', label: 'Blue' }, { id: 'opaque-two', label: 'Red' }], audio_url: '/audio/eq-listen-a1-1.wav' }

it('choice questions send opaque IDs, not scores', () => {
  const change = vi.fn()
  render(<ChoiceQuestion question={question} onChange={change} />)
  fireEvent.click(screen.getByRole('button', { name: 'Blue' }))
  expect(change).toHaveBeenCalledWith('opaque-one')
})
it('written practice sends the text for server validation', () => {
  const change = vi.fn()
  render(<WrittenAnswer onChange={change} />)
  fireEvent.change(screen.getByRole('textbox'), { target: { value: 'She goes to school.' } })
  expect(change).toHaveBeenCalledWith('She goes to school.')
})
it('unplayable audio is reported honestly and does not unlock choices', () => {
  render(<SoundDetective question={question} onChange={vi.fn()} />)
  expect(screen.getByRole('button', { name: 'Blue' })).toBeDisabled()
  fireEvent.error(screen.getByLabelText('Audio del desafío'))
  expect(screen.getByRole('alert')).toHaveTextContent('No se simulará la escucha')
  expect(screen.getByRole('button', { name: 'Blue' })).toBeDisabled()
})
it('booster activation calls the API rather than changing local XP', async () => {
  api.post.mockResolvedValue({ data: { data: { multiplier: 2 } } })
  const reload = vi.fn()
  render(<BoosterInventory boosters={[{ id: 'booster-id', title: 'Doble XP', available: true }]} serverTime="2026-10-27T12:00:00Z" reload={reload} />)
  fireEvent.click(screen.getByRole('button', { name: 'Activar doble XP' }))
  await waitFor(() => expect(reload).toHaveBeenCalled())
  expect(api.post).toHaveBeenCalledWith('/boosters/booster-id/activate')
})
it('installation controls are absent when the browser has not offered installation', () => {
  render(<PWAControls />)
  expect(screen.queryByRole('button', { name: 'Instalar EnglishQuest' })).not.toBeInTheDocument()
})
it('an actual browser install offer is consumed once', async () => {
  render(<PWAControls />)
  const event = new Event('beforeinstallprompt', { cancelable: true })
  event.prompt = vi.fn().mockResolvedValue()
  event.userChoice = Promise.resolve({ outcome: 'dismissed' })
  fireEvent(window, event)
  fireEvent.click(await screen.findByRole('button', { name: 'Instalar EnglishQuest' }))
  await waitFor(() => expect(screen.queryByRole('button', { name: 'Instalar EnglishQuest' })).not.toBeInTheDocument())
  expect(event.prompt).toHaveBeenCalledTimes(1)
})
