import { render, screen, fireEvent } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import { MemoryMatch, SentenceBuilder, WordHunter } from './Games.jsx'

const question = {
  options: [
    { id: 'opaque-a', label: 'Apple', image: 'apple' },
    { id: 'opaque-b', label: 'Manzana', image: 'banana' },
  ],
}

it('Word Hunter submits only the selected option ID, not a score', () => {
  const change = vi.fn()
  render(<WordHunter question={question} onChange={change} />)
  fireEvent.click(screen.getByRole('button', { name: 'Apple' }))
  expect(change).toHaveBeenCalledWith('opaque-a')
})
it('Sentence Builder appends tokens in learner-selected order', () => {
  const change = vi.fn()
  render(
    <SentenceBuilder
      question={question}
      value={['opaque-a']}
      onChange={change}
    />,
  )
  fireEvent.click(screen.getByRole('button', { name: 'Manzana' }))
  expect(change).toHaveBeenCalledWith(['opaque-a', 'opaque-b'])
})
it('Memory Match forms pairs without claiming they are correct', () => {
  const change = vi.fn()
  render(<MemoryMatch question={question} onChange={change} />)
  fireEvent.click(screen.getByRole('button', { name: 'Revelar tarjeta 1' }))
  expect(screen.getByRole('button', { name: 'Apple' })).toBeInTheDocument()
  fireEvent.click(screen.getByRole('button', { name: 'Revelar tarjeta 2' }))
  expect(change).toHaveBeenCalledWith([['opaque-a', 'opaque-b']])
  expect(screen.queryByText('Correcto')).not.toBeInTheDocument()
})
