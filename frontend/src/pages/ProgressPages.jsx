import { Link } from 'react-router-dom'
import { Award, Check, LockKeyhole, Target } from 'lucide-react'
import useResource from '../lib/useResource.js'
import { Empty, PageTitle, Progress, ResourceState } from '../components/UI.jsx'
import BoosterInventory from '../components/BoosterInventory.jsx'

export default function ProgressPages({ kind }) {
  const resource = useResource('/dashboard')
  const data = resource.data
  const title =
    kind === 'achievements'
      ? 'Pequeñas victorias, grandes recuerdos'
      : kind === 'missions'
        ? 'Un objetivo para hoy'
        : 'Mira cuánto has avanzado'
  const completed = data?.lessons.filter((l) => l.completed) || []
  return (
    <>
      <PageTitle
        eyebrow={
          kind === 'achievements'
            ? 'TUS INSIGNIAS'
            : kind === 'missions'
              ? 'MISIONES DIARIAS'
              : 'PROGRESO ACADÉMICO'
        }
        title={title}
      >
        {kind === 'missions'
          ? 'Las recompensas se entregan automáticamente al cumplir el objetivo, según tu zona horaria.'
          : 'Cada dato viene de tu práctica guardada. Sigue aprendiendo a tu ritmo.'}
      </PageTitle>
      <ResourceState resource={resource}>
        {data && (
          <>
            {kind === 'achievements' && (
              <div className="achievement-grid">
                {data.achievements.map((a) => (
                  <article
                    className={`card achievement ${a.unlocked ? 'earned' : ''}`}
                    key={a.code}
                  >
                    <span className="achievement-icon">
                      {a.unlocked ? (
                        <Award size={38} />
                      ) : (
                        <LockKeyhole size={30} />
                      )}
                    </span>
                    <h2>{a.title}</h2>
                    <p>{a.description}</p>
                    <span className="pill">
                      {a.unlocked ? 'Desbloqueado' : 'Por descubrir'}
                    </span>
                  </article>
                ))}
              </div>
            )}
            {kind === 'missions' && (
              <>
                {data.missions.map((m) => (
                  <article className="card mission-detail" key={m.code}>
                    <span className="feature-icon gold">
                      <Target />
                    </span>
                    <h2>{m.title}</h2>
                    <p>
                      {m.description} Repetir una actividad ya aprobada no aumenta esta misión.
                    </p>
                    <Progress
                      value={Math.min(100, (m.progress * 100) / m.target)}
                    />
                    <p>
                      {Math.min(m.progress, m.target)} / {m.target} ·
                      Recompensa: {m.reward} XP
                    </p>
                    {m.rewarded ? (
                      <p className="inline success-note">
                        <Check /> ¡Misión completada y recompensa guardada!
                      </p>
                    ) : (
                      <Link to="/mapa" className="btn">
                        Ir a practicar
                      </Link>
                    )}
                  </article>
                ))}
                <BoosterInventory boosters={data.boosters} serverTime={data.server_time} reload={resource.reload} />
                <p className="fine-print">
                  Se considera el día local de tu cuenta. La muestra tiene
                  contenido limitado: cuando lo completes, puedes seguir
                  repasando aunque no haya actividades nuevas para esta misión.
                </p>
              </>
            )}
            {kind === 'statistics' && (
              <>
                <section className="stats-grid">
                  <article className="stat-card">
                    <div>
                      <strong>{completed.length}</strong>
                      <span>Actividades distintas aprobadas</span>
                    </div>
                  </article>
                  <article className="stat-card">
                    <div>
                      <strong>{data.xp}</strong>
                      <span>XP acumulada</span>
                    </div>
                  </article>
                  <article className="stat-card">
                    <div>
                      <strong>
                        {completed.length
                          ? Math.round(
                              completed.reduce((n, l) => n + l.best_score, 0) /
                                completed.length,
                            )
                          : 0}
                        %
                      </strong>
                      <span>Promedio de mejores puntuaciones aprobadas</span>
                    </div>
                  </article>
                </section>
                <article className="card">
                  <h2>Progreso de la muestra por nivel</h2>
                  {[...new Set(data.lessons.map(l => l.level_code))].map((level) => {
                    const lessons = data.lessons.filter(
                      (l) => l.level_code === level,
                    )
                    const done = lessons.filter((l) => l.completed).length
                    return (
                      <div className="level-progress" key={level}>
                        <p>
                          {level} · {done} de {lessons.length} actividades
                        </p>
                        <Progress
                          value={
                            lessons.length
                              ? Math.round((done * 100) / lessons.length)
                              : 0
                          }
                          label={`Progreso ${level}`}
                        />
                      </div>
                    )
                  })}
                  <p className="fine-print">
                    Palabras objetivo practicadas en primeras aprobaciones: {data.words_practiced}. No son necesariamente palabras diferentes.{' '}
                    Estas métricas no certifican dominio del nivel CEFR.
                  </p>
                </article>
                <h2 className="section-title">Últimos 20 intentos</h2>
                {data.history.length ? (
                  <div className="card history-list">
                    {data.history.map((h) => (
                      <Link to={`/resultados/${h.id}`} key={h.id}>
                        <span>
                          <strong>{h.title}</strong>
                          <small>
                            {new Date(h.completed_at).toLocaleString('es')}
                          </small>
                        </span>
                        <span>{h.score}% →</span>
                      </Link>
                    ))}
                  </div>
                ) : (
                  <Empty>Todavía no hay intentos completados.</Empty>
                )}
              </>
            )}
          </>
        )}
      </ResourceState>
    </>
  )
}
