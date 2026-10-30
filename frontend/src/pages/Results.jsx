import { Link, useParams } from 'react-router-dom'
import { Sparkles, Star, Trophy } from 'lucide-react'
import useResource from '../lib/useResource.js'
import { PageTitle, ResourceState } from '../components/UI.jsx'

function solutionText(question, solution) {
  const label = (id) => question.options.find((o) => o.id === id)?.label || '—'
  if (typeof solution === 'string') return question.options.some(o => o.id === solution) ? label(solution) : solution
  if (Array.isArray(solution?.[0]))
    return solution.map((pair) => pair.map(label).join(' ↔ ')).join(' · ')
  return Array.isArray(solution) ? solution.map(label).join(' ') : '—'
}

export default function Results() {
  const { id } = useParams()
  const resource = useResource(`/attempts/${id}`)
  const diploma = useResource('/diploma')
  const { result, activity } = resource.data || {}
  return (
    <ResourceState resource={resource}>
      {result ? (
        <>
          <section className={`result-hero ${result.passed ? 'passed' : ''}`}>
            <span className="result-medal">
              <Trophy size={50} />
            </span>
            <PageTitle
              eyebrow="RESULTADO GUARDADO"
              title={
                result.passed
                  ? '¡Un paso más en tu aventura!'
                  : 'Cada intento te enseña algo'
              }
            >
            {result.passed
                ? 'Tu práctica ya forma parte de tu progreso.'
                : 'Repasa la explicación y vuelve a intentarlo. No hay penalización.'}
            </PageTitle>
            {result.expired && <p role="status" className="notice">El servidor recibió este intento fuera de plazo. La puntuación es 0 y no se entregaron recompensas. Puedes repasar e iniciar un intento nuevo.</p>}
            <div className="result-score">
              {result.score}
              <span>/100</span>
            </div>
            <div className="stars large">
              {[1, 2, 3].map((n) => (
                <Star
                  key={n}
                  fill={result.stars >= n ? 'currentColor' : 'none'}
                />
              ))}
            </div>
            <p className="xp-result">
              <Sparkles size={20} /> +{result.xp} XP en este intento
            </p>
            {result.passed && !result.new_completion && (
              <small>
                Actividad ya completada: no se repite la recompensa.
              </small>
            )}
            {result.multiplier === 2 && <p className="success-note">Potenciador aplicado: ×2 a la XP base. Los bonos no se multiplican.</p>}
            <div className="hero-actions">
              <Link className="btn" to="/mapa">
                Continuar mi ruta
              </Link>
              <Link
                className="btn secondary"
                to={`/lecciones/${activity.lesson_id}`}
              >
                Repasar y practicar
              </Link>
            </div>
          </section>
          {result.passed && (diploma.data?.eligible || diploma.data?.diploma) && <section className="card"><h2>¡Terminaste las actividades del curso!</h2><p>Tu reconocimiento es simbólico, no una certificación CEFR.</p><Link className="btn" to="/diploma">{diploma.data.diploma ? 'Ver mi diploma' : 'Elegir nombre y obtener mi diploma'}</Link></section>}
          <section className="feedback-list">
            <h2>Lo que aprendiste</h2>
            {result.feedback.map((feedback) => {
              const question = activity.questions.find(
                (q) => q.id === feedback.question_id,
              )
              return (
                <article
                  className={`card feedback ${feedback.correct ? 'correct' : 'incorrect'}`}
                  key={feedback.question_id}
                >
                  <span className="eyebrow">
                    {feedback.correct ? '✓ CORRECTO' : '↻ PARA REPASAR'}
                  </span>
                  <h3>{question?.prompt}</h3>
                  <p lang="en">{feedback.explanation}</p>
                  {question && (
                    <p className="solution">
                      Respuesta: {solutionText(question, feedback.solution)}
                    </p>
                  )}
                </article>
              )
            })}
          </section>
        </>
      ) : (
        <div className="card">
          <p>Este intento todavía no tiene un resultado.</p>
          <Link className="btn" to={`/ejercicios/${id}`}>
            Continuar actividad
          </Link>
        </div>
      )}
    </ResourceState>
  )
}
