import { useId } from 'react'
import { answerType, newQuestion } from '../lib/gameCatalog.js'

export default function QuestionEditor({ questions, onChange, kind, resources }) {
  const editorID = useId()
  const images = resources.filter(r => r.kind === 'image' && r.active)
  const audio = resources.filter(r => r.kind === 'audio' && r.active)
  function update(index, patch) { onChange(questions.map((q, i) => i === index ? { ...q, ...patch } : q)) }
  function reset(index, q, response = q.response) {
    const source = ({ order: 'sentence_builder', pairs: 'memory_match', text: 'grammar_rescue', letters: 'spell_master', multi: 'word_rain' })[response] || kind
    const fresh = newQuestion(source)
    update(index, { options: fresh.options, solution: fresh.solution, response: kind === 'english_rush' ? response : '' })
  }
  return <div className="question-editor">
    <p className="fine-print">Las soluciones solo se usan en Go. No las escribas en el enunciado ni en el texto público. Cambiar el formato reinicia opciones y solución; el resto del contenido se conserva.</p>
    {questions.map((q, index) => {
      const mode = answerType(kind, q)
      const sequence = Array.isArray(q.solution) ? q.solution : []
      const pairs = sequence.filter(p => Array.isArray(p) && p.length === 2)
      return <fieldset key={q.id || index} className="question-authoring"><legend>Pregunta {index + 1}</legend>
        {kind === 'english_rush' && <label>Formato de respuesta<select value={q.response || 'choice'} onChange={e => reset(index, q, e.target.value)}><option value="choice">Elegir opción</option><option value="order">Ordenar palabras</option><option value="text">Escribir respuesta</option><option value="letters">Ordenar letras</option></select></label>}
        <div><label htmlFor={`${editorID}-${index}-prompt`}>Enunciado</label><textarea id={`${editorID}-${index}-prompt`} required value={q.prompt} onChange={e => update(index, { prompt: e.target.value })} maxLength={500} /></div>
        <label>Texto público de apoyo / historia / conversación<textarea value={q.passage || ''} maxLength={5000} required={kind === 'story_explorer' || kind === 'conversation_quest'} onChange={e => update(index, { passage: e.target.value })} rows={5} /></label>
        {mode === 'choice' && <label>Audio {kind === 'sound_detective' ? '(obligatorio)' : '(opcional)'}<select required={kind === 'sound_detective'} value={q.audio_key || ''} onChange={e => update(index, { audio_key: e.target.value })}><option value="">Sin audio</option>{audio.map(r => <option value={r.key} key={r.key}>{r.title}</option>)}</select></label>}
        {mode !== 'choice' && q.audio_key && <button type="button" className="text-button" onClick={() => update(index, { audio_key: '' })}>Quitar audio incompatible</button>}
        <div><label htmlFor={`${editorID}-${index}-explanation`}>Explicación después de responder</label><textarea id={`${editorID}-${index}-explanation`} required value={q.explanation} onChange={e => update(index, { explanation: e.target.value })} maxLength={1500} /></div>
        {mode !== 'text' && <><h3>Opciones {mode === 'letters' ? '/ letras' : mode === 'order' ? '/ palabras' : ''}</h3>{q.options.map((option, j) => <div className="option-authoring" key={option.id}>
          <label>Texto {j + 1}<input value={option.label} maxLength={mode === 'letters' ? 1 : 120} required onChange={e => update(index, { options: q.options.map(o => o.id === option.id ? { ...o, label: e.target.value } : o) })} /></label>
          {kind === 'word_hunter' && <label>Imagen {j + 1}<select required value={option.image || ''} onChange={e => update(index, { options: q.options.map(o => o.id === option.id ? { ...o, image: e.target.value } : o) })}><option value="">Elige una imagen</option>{images.map(r => <option key={r.key} value={r.key}>{r.title}</option>)}</select></label>}
          <button className="text-button" type="button" disabled={q.options.length <= 2} onClick={() => update(index, { options: q.options.filter(o => o.id !== option.id) })}>Quitar opción {j + 1}</button>
        </div>)}<button className="btn secondary small" type="button" disabled={q.options.length >= 12} onClick={() => update(index, { options: [...q.options, { id: `o-${crypto.randomUUID().replaceAll('-', '').slice(0, 24)}`, label: '', image: '' }] })}>Añadir opción</button></>}
        {mode === 'text' && q.options.length > 0 && <p className="notice">Este formato no usa opciones. <button type="button" className="text-button" onClick={() => update(index, { options: [] })}>Quitar las opciones para usar respuesta escrita</button></p>}
        <h3 className="section-title">Solución (solo administración)</h3>
        {mode === 'choice' && <div><label htmlFor={`${editorID}-${index}-solution`}>Opción correcta</label><select id={`${editorID}-${index}-solution`} value={typeof q.solution === 'string' ? q.solution : ''} onChange={e => update(index, { solution: e.target.value })}><option value="">Selecciona la respuesta</option>{q.options.map(o => <option value={o.id} key={o.id}>{o.label || 'Sin texto'}</option>)}</select></div>}
        {(mode === 'text' || mode === 'letters') && <label>{mode === 'letters' ? 'Palabra correcta' : 'Respuesta escrita esperada'}<input required maxLength={500} value={typeof q.solution === 'string' ? q.solution : ''} onChange={e => update(index, { solution: e.target.value })} /></label>}
        {mode === 'multi' && <div>{q.options.map(o => <label className="checkbox-label" key={o.id}><input type="checkbox" checked={sequence.includes(o.id)} onChange={e => update(index, { solution: e.target.checked ? [...sequence, o.id] : sequence.filter(id => id !== o.id) })} />{o.label || 'Sin texto'} · correcta</label>)}<small>Marca al menos una respuesta correcta y conserva algún distractor.</small></div>}
        {mode === 'order' && <div><p>Elige las palabras en el orden correcto.</p><div className="word-bank">{q.options.map(o => <button className="word-token" type="button" key={o.id} disabled={sequence.includes(o.id)} onClick={() => update(index, { solution: [...sequence, o.id] })}>{o.label || 'Sin texto'}</button>)}</div><p className="solution">{sequence.map(value => q.options.find(o => o.id === value)?.label || 'Opción eliminada').join(' ')}</p><button className="text-button" type="button" onClick={() => update(index, { solution: [] })}>Rehacer orden correcto</button></div>}
        {mode === 'pairs' && <div>{pairs.map((pair, pairIndex) => <div className="form-row" key={pairIndex}>{[0, 1].map(side => <label key={side}>Pareja {pairIndex + 1} · lado {side + 1}<select value={pair[side] || ''} onChange={e => update(index, { solution: pairs.map((p, i) => i === pairIndex ? p.map((v, s) => s === side ? e.target.value : v) : p) })}><option value="">Selecciona</option>{q.options.map(o => <option key={o.id} value={o.id}>{o.label || 'Sin texto'}</option>)}</select></label>)}<button className="text-button" type="button" onClick={() => update(index, { solution: pairs.filter((_, i) => i !== pairIndex) })}>Quitar pareja {pairIndex + 1}</button></div>)}<button className="btn secondary small" type="button" onClick={() => update(index, { solution: [...pairs, ['', '']] })}>Añadir pareja</button></div>}
        <button className="text-button" type="button" onClick={() => reset(index, q, kind === 'english_rush' ? q.response || 'choice' : '')}>Reiniciar opciones y solución de esta pregunta</button>
        <button className="text-button remove-question" type="button" disabled={questions.length <= 1} onClick={() => onChange(questions.filter((_, i) => i !== index))}>Quitar pregunta {index + 1}</button>
      </fieldset>
    })}
    <button className="btn secondary" type="button" disabled={questions.length >= 12} onClick={() => onChange([...questions, newQuestion(kind)])}>Añadir pregunta</button>
  </div>
}
