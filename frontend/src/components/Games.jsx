import { useState } from 'react'

export function WordHunter({ question, value, onChange }) {
  return (
    <div className="picture-options">
      {question.options.map((option) => (
        <PictureChoice option={option} selected={value === option.id} choose={onChange} key={option.id} />
      ))}
    </div>
  )
}

function PictureChoice({ option, selected, choose }) {
  const [broken, setBroken] = useState(false)
  return <button
          key={option.id}
          className={`picture-option ${selected ? 'chosen' : ''}`}
          aria-label={option.label}
          aria-pressed={selected}
          disabled={broken}
          onClick={() => choose(option.id)}
        >
          {broken ? <span role="status">Imagen no disponible</span> : <img src={option.image_url || `/images/${option.image}.svg`} alt={option.label} referrerPolicy="no-referrer" onError={() => setBroken(true)} />}
          <span aria-hidden="true">{option.label}</span>
        </button>
}

export function SentenceBuilder({ question, value = [], onChange }) {
  return (
    <div>
      <div
        className="sentence-answer"
        aria-label="Tu construcción"
        aria-live="polite"
      >
        {value.length ? (
          value.map((id) => (
            <button
              key={id}
              className="word-token selected"
              onClick={() => onChange(value.filter((v) => v !== id))}
              title="Quitar palabra"
            >
              {question.options.find((o) => o.id === id)?.label}
            </button>
          ))
        ) : (
          <span className="muted">
            Selecciona las palabras o letras en el orden correcto…
          </span>
        )}
      </div>
      <div className="word-bank">
        {question.options.map((option) => (
          <button
            key={option.id}
            className="word-token"
            disabled={value.includes(option.id)}
            onClick={() => onChange([...value, option.id])}
          >
            {option.label}
          </button>
        ))}
      </div>
      <button className="text-button" onClick={() => onChange([])}>
        Reiniciar oración
      </button>
    </div>
  )
}

export function MemoryMatch({ question, value = [], onChange }) {
  const [open, setOpen] = useState(null)
  const selected = value.flat()
  function reveal(id) {
    if (open === id) {
      setOpen(null)
      return
    }
    if (open) {
      onChange([...value, [open, id]])
      setOpen(null)
    } else setOpen(id)
  }
  return (
    <div>
      <p className="muted">
        Forma {question.options.length / 2} parejas. La corrección se mostrará
        al enviar, no al girar las tarjetas.
      </p>
      <div className="memory-grid">
        {question.options.map((option, i) => {
          const paired = selected.includes(option.id)
          const revealed = paired || open === option.id
          return (
            <button
              className={`memory-card ${revealed ? 'revealed' : ''} ${paired ? 'paired' : ''}`}
              key={option.id}
              disabled={paired}
              aria-label={revealed ? option.label : `Revelar tarjeta ${i + 1}`}
              onClick={() => reveal(option.id)}
            >
              {revealed ? option.label : <span aria-hidden="true">?</span>}
            </button>
          )
        })}
      </div>
      <p>
        {value.length} / {question.options.length / 2} parejas propuestas
      </p>
      <button
        className="text-button"
        onClick={() => {
          onChange([])
          setOpen(null)
        }}
      >
        Reiniciar parejas
      </button>
    </div>
  )
}
