import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { ArrowRight, BookOpen, LockKeyhole } from 'lucide-react'
import useResource from '../lib/useResource.js'
import { api, errorMessage } from '../lib/api.js'
import { Notice, PageTitle, ResourceState } from '../components/UI.jsx'

export default function Lesson() {
  const { id } = useParams()
  const resource = useResource('/lessons')
  const lesson = resource.data?.find((l) => l.id === id)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  async function start() {
    setBusy(true)
    setError('')
    try {
      const { data } = await api.post(`/activities/${lesson.activity_id}/start`)
      navigate(`/ejercicios/${data.data.id}`)
    } catch (e) {
      setError(errorMessage(e))
      setBusy(false)
    }
  }
  return (
    <>
      <Link className="text-link" to="/mapa">
        ← Volver al mapa
      </Link>
      <ResourceState resource={resource}>
        {lesson ? (
          <>
            <PageTitle
              eyebrow={`${lesson.level_code} · ${lesson.unit_title}`}
              title={lesson.title}
            >
              {lesson.objective}
            </PageTitle>
            <article className="card lesson-card">
              <span className="feature-icon blue">
                <BookOpen />
              </span>
              <h2>Antes de jugar</h2>
              <p className="lesson-content" lang="en">
                {lesson.content}
              </p>
              <aside className="learning-tip">
                Consejo: lee los ejemplos en voz alta y piensa en su
                significado. No necesitas responder rápido.
              </aside>
              {lesson.assessment && (
                <p>
                  Evaluación de unidad: necesitas al menos 80%. Puedes volver a
                  practicar sin penalización.
                </p>
              )}
              {error && <Notice>{error}</Notice>}
              {lesson.locked ? (
                <p className="inline">
                  <LockKeyhole /> Completa la lección anterior para desbloquear
                  esta actividad.
                </p>
              ) : (
                <button onClick={start} className="btn" disabled={busy}>
                  {busy
                    ? 'Preparando actividad…'
                    : lesson.completed
                      ? 'Volver a practicar'
                      : 'Empezar desafío'}{' '}
                  <ArrowRight size={18} />
                </button>
              )}
              <p className="fine-print">
                XP una sola vez por actividad aprobada. Repetir sirve para
                aprender y mejorar tu puntuación.
              </p>
            </article>
          </>
        ) : (
          <div className="card">Lección no encontrada.</div>
        )}
      </ResourceState>
    </>
  )
}
