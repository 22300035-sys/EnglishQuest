import { useState } from 'react'
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom'
import useResource from '../lib/useResource.js'
import { api, errorMessage } from '../lib/api.js'
import { Notice, PageTitle, ResourceState } from '../components/UI.jsx'
import {
  WordHunter,
  SentenceBuilder,
  MemoryMatch,
} from '../components/Games.jsx'
import { AttemptClock, ChoiceQuestion, MixedQuestion, SoundDetective, WordRain, WrittenAnswer } from '../components/AdvancedGames.jsx'
import { answerType } from '../lib/gameCatalog.js'

const games = {
  word_hunter: WordHunter,
  sentence_builder: SentenceBuilder,
  memory_match: MemoryMatch,
  sound_detective: SoundDetective,
  grammar_rescue: WrittenAnswer,
  spell_master: SentenceBuilder,
  conversation_quest: ChoiceQuestion,
  story_explorer: ChoiceQuestion,
  word_rain: WordRain,
  english_rush: MixedQuestion,
}

export default function Exercise() {
  const { id } = useParams()
  const resource = useResource(`/attempts/${id}`)
  const attempt = resource.data
  const [answers, setAnswers] = useState({})
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  if (attempt?.result) return <Navigate to={`/resultados/${id}`} replace />
  const activity = attempt?.activity
  const Game = activity && games[activity.kind]
  const ready = activity?.questions.every(q => {
    const type = answerType(activity.kind, q)
    if (type === 'choice' || type === 'text') return typeof answers[q.id] === 'string' && answers[q.id].trim().length > 0
    if (type === 'multi') return Array.isArray(answers[q.id]) && answers[q.id].length > 0
    return answers[q.id]?.length === (type === 'pairs' ? q.options.length / 2 : q.options.length)
  })
  const canSubmit = ready || activity?.time_limit_seconds > 0
  async function submit() {
    setBusy(true)
    setError('')
    try {
      await api.post(`/attempts/${id}/submit`, { answers })
      navigate(`/resultados/${id}`, { replace: true })
    } catch (e) {
      setError(errorMessage(e))
      setBusy(false)
    }
  }
  return (
    <ResourceState resource={resource}>
      {activity && (
        <>
          <Link to={`/lecciones/${activity.lesson_id}`} className="text-link">
            ← Repasar la lección
          </Link>
          <PageTitle
            eyebrow={
              activity.assessment
                ? activity.time_limit_seconds ? 'EVALUACIÓN · PLAZO CONTROLADO POR EL SERVIDOR' : 'EVALUACIÓN · SIN LÍMITE DE TIEMPO'
                : 'PRACTICA A TU RITMO'
            }
            title={activity.title}
          >
            {activity.instructions}
          </PageTitle>
          {attempt.expires_at && <AttemptClock key={attempt.id} expiresAt={attempt.expires_at} serverTime={attempt.server_time} />}
          {Game ? (
            activity.questions.map((q, i) => (
              <section className="card question-card" key={q.id}>
                <span className="eyebrow">
                  DESAFÍO {i + 1} / {activity.questions.length}
                </span>
                <h2 lang={activity.kind === 'sentence_builder' ? 'es' : 'en'}>
                  {q.prompt}
                </h2>
                {q.passage && <article className="reading-passage" lang="en" aria-label="Texto de la actividad">{q.passage}</article>}
                <Game
                  question={q}
                  value={answers[q.id]}
                  onChange={(value) =>
                    setAnswers((previous) => ({ ...previous, [q.id]: value }))
                  }
                />
              </section>
            ))
          ) : (
            <Notice>Este tipo de actividad no está disponible.</Notice>
          )}
          {error && <Notice>{error}</Notice>}
          <div className="exercise-footer">
            <p>
              {ready
                ? '¡Todo listo! Revisa tus respuestas antes de enviar.'
                : activity?.time_limit_seconds ? 'Puedes enviar respuestas parciales; las omitidas cuentan como incorrectas.' : 'Completa todos los desafíos para enviar.'}
            </p>
            <button className="btn" disabled={!canSubmit || busy || !Game} onClick={submit}>
              {busy ? 'Comprobando y guardando…' : 'Comprobar respuestas'}
            </button>
          </div>
        </>
      )}
    </ResourceState>
  )
}
