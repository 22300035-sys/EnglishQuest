import { render, screen, fireEvent } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import QuestionEditor from './QuestionEditor.jsx'

it('keeps exact question labels usable when editing saved content', () => {
  const onChange = vi.fn()
  const question = {
    id: 'saved-question', prompt: 'Find: apple', explanation: 'Apple means manzana.',
    options: [{ id: 'a', label: 'Manzana', image: 'apple' }, { id: 'b', label: 'Plátano', image: 'banana' }], solution: 'a',
  }
  render(<QuestionEditor questions={[question]} kind="word_hunter" resources={[]} onChange={onChange} />)
  fireEvent.change(screen.getByLabelText('Enunciado', { exact: true }), { target: { value: 'Find: banana' } })
  expect(onChange).toHaveBeenLastCalledWith([{ ...question, prompt: 'Find: banana' }])
  fireEvent.change(screen.getByLabelText('Explicación después de responder', { exact: true }), { target: { value: 'Banana means plátano.' } })
  expect(onChange).toHaveBeenLastCalledWith([{ ...question, explanation: 'Banana means plátano.' }])
  fireEvent.change(screen.getByLabelText('Opción correcta', { exact: true }), { target: { value: 'b' } })
  expect(onChange).toHaveBeenLastCalledWith([{ ...question, solution: 'b' }])
})
