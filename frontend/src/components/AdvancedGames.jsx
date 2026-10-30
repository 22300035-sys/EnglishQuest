import { useEffect, useRef, useState } from 'react'
import { SentenceBuilder, MemoryMatch } from './Games.jsx'
import { answerType } from '../lib/gameCatalog.js'

export function ChoiceQuestion({ question, value, onChange, disabled = false }) {
  return <div className="choice-options">{question.options.map(option => <button key={option.id} className={`choice-option ${value === option.id ? 'chosen' : ''}`} aria-pressed={value === option.id} disabled={disabled} onClick={() => onChange(option.id)}><span>{option.label}</span>{value === option.id && <span aria-hidden="true">✓</span>}</button>)}</div>
}
export function WrittenAnswer({ value = '', onChange }) {
  return <label className="written-answer">Tu respuesta en inglés<textarea lang="en" autoComplete="off" autoCapitalize="sentences" spellCheck={false} value={value} maxLength={1000} onChange={e => onChange(e.target.value)} rows={3} placeholder="Escribe la corrección o la palabra solicitada…" /><small>Se comprueba la corrección indicada; no todas las paráfrasis posibles. Mayúsculas, espacios repetidos y puntuación final no afectan.</small></label>
}
export function SoundDetective(props) {
  const audio = useRef(null)
  const [played, setPlayed] = useState(false)
  const [error, setError] = useState(false)
  return <div><audio ref={audio} src={props.question.audio_url} controls preload="metadata" crossOrigin="anonymous" onPlay={() => setPlayed(true)} onError={() => setError(true)} aria-label="Audio del desafío" />
    <p className="fine-print">Narración sintética pregrabada. Puedes escucharla varias veces; la transcripción se muestra al terminar.</p>
    {error ? <p role="alert" className="notice">No se pudo reproducir el audio. No se simulará la escucha.<button className="text-button" onClick={() => { setError(false); audio.current?.load() }}>Reintentar audio</button></p> : !played && <p className="muted">Reproduce el audio para habilitar las respuestas.</p>}
    <ChoiceQuestion {...props} disabled={!played || error} />
  </div>
}
export function WordRain({ question, value = [], onChange }) {
  function toggle(id) { onChange(value.includes(id) ? value.filter(v => v !== id) : [...value, id]) }
  return <div><p className="muted">Las palabras vuelven a aparecer. Selecciona solo la categoría indicada; puedes quitar una selección antes de enviar.</p><div className="rain-field">{question.options.map((option, i) => <button key={option.id} className={`rain-word ${value.includes(option.id) ? 'caught' : ''}`} style={{ '--lane': `${(i % 3) * 31 + 2}%`, '--delay': `${i * -1.2}s`, '--duration': `${8 + (i % 3)}s` }} aria-pressed={value.includes(option.id)} onClick={() => toggle(option.id)}>{option.label}</button>)}</div><div className="caught-words" aria-live="polite"><strong>Tus palabras:</strong>{value.length ? value.map(id => <button className="word-token" key={id} onClick={() => toggle(id)}>{question.options.find(o => o.id === id)?.label} ×</button>) : <span className="muted">Todavía no has seleccionado ninguna.</span>}</div></div>
}
export function MixedQuestion(props) {
  switch (answerType('english_rush', props.question)) {
    case 'order': case 'letters': return <SentenceBuilder {...props} />
    case 'pairs': return <MemoryMatch {...props} />
    case 'text': return <WrittenAnswer {...props} />
    default: return props.question.audio_url ? <SoundDetective {...props} /> : <ChoiceQuestion {...props} />
  }
}
export function AttemptClock({ expiresAt, serverTime }) {
  const [base] = useState(() => ({ at: performance.now(), left: Math.max(0, (new Date(expiresAt) - new Date(serverTime)) / 1000) }))
  const [left, setLeft] = useState(Math.ceil(base.left))
  useEffect(() => { const timer = setInterval(() => setLeft(Math.max(0, Math.ceil(base.left - (performance.now() - base.at) / 1000))), 250); return () => clearInterval(timer) }, [base])
  return <aside className={`attempt-clock ${left === 0 ? 'expired' : ''}`}><span role="timer" aria-label="Tiempo restante" aria-live="off">{Math.floor(left / 60)}:{String(left % 60).padStart(2, '0')}</span><p>{left ? 'El servidor controla el plazo total, incluso al recargar. Puedes enviar respuestas parciales.' : 'El plazo terminó. Envía para ver el resultado; los intentos vencidos no dan XP.'}</p></aside>
}
