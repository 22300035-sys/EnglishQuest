import {
  ArrowRight,
  BookOpen,
  Flame,
  Gamepad2,
  Sparkles,
  Star,
  Target,
  Trophy,
} from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../auth/context.js'
import useResource from '../lib/useResource.js'
import { PageTitle, Progress, ResourceState } from '../components/UI.jsx'
import { gameCatalog } from '../lib/gameCatalog.js'

export default function Dashboard() {
  const { user } = useAuth()
  const resource = useResource('/dashboard')
  const data = resource.data
  const completed = data?.lessons.filter((l) => l.completed).length || 0
  const progress = data?.lessons.length
    ? Math.round((completed * 100) / data.lessons.length)
    : 0
  const next =
    data?.lessons.find(
      (l) => !l.completed && !l.locked && l.level_code === user.cefr_level,
    ) || data?.lessons.find((l) => !l.completed && !l.locked)
  return (
    <>
      <PageTitle
        eyebrow="CADA PASO CUENTA"
        title={`¡Hola, ${user.first_name}!`}
      >
        Hoy es un buen día para descubrir algo nuevo.
      </PageTitle>
      <ResourceState resource={resource}>
        {data && (
          <>
            <section className="dashboard-grid">
              <article className="journey-banner">
                <div>
                  <span className="pill light">TU PRÓXIMO DESCUBRIMIENTO</span>
                  <h2>
                    {next ? next.title : '¡Completaste la muestra inicial!'}
                  </h2>
                  <p>
                    {next
                      ? next.objective
                      : 'Puedes repasar las actividades y mejorar tus resultados. Pronto habrá más contenido.'}
                  </p>
                  <Link
                    to={next ? `/lecciones/${next.id}` : '/mapa'}
                    className="btn gold"
                  >
                    {next ? 'Continuar aprendiendo' : 'Repasar mi ruta'}{' '}
                    <ArrowRight size={18} />
                  </Link>
                </div>
                <div className="banner-compass" aria-hidden="true">
                  <BookOpen size={64} />
                  <Star size={26} fill="currentColor" />
                </div>
              </article>
              <article className="card streak-card">
                <span className="fire">
                  <Flame size={42} />
                </span>
                <h2>
                  {data.streak.current_days}{' '}
                  {data.streak.current_days === 1 ? 'día' : 'días'}
                </h2>
                <p>Tu racha de aprendizaje</p>
                <div className="week-calendar">
                  {Array.from({ length: 7 }, (_, i) => {
                    const date = new Date(`${data.today}T12:00:00Z`)
                    date.setUTCDate(date.getUTCDate() - 6 + i)
                    const day = date.toISOString().slice(0, 10)
                    return (
                      <span
                        key={day}
                        className={
                          data.practice_days.includes(day) ? 'practiced' : ''
                        }
                        title={day}
                      >
                        {date.toLocaleDateString('es', {
                          weekday: 'narrow',
                          timeZone: 'UTC',
                        })}
                        <i>{data.practice_days.includes(day) ? '✓' : '·'}</i>
                      </span>
                    )
                  })}
                </div>
                <small>
                  Récord: {data.streak.best_days} días · {user.timezone}
                </small>
              </article>
            </section>
            <section className="stats-grid">
              <article className="stat-card">
                <span className="feature-icon gold">
                  <Sparkles />
                </span>
                <div>
                  <strong>{data.xp}</strong>
                  <span>XP acumulada</span>
                </div>
              </article>
              <article className="stat-card">
                <span className="feature-icon blue">
                  <BookOpen />
                </span>
                <div>
                  <strong>
                    {completed} / {data.lessons.length}
                  </strong>
                  <span>Actividades aprobadas</span>
                </div>
              </article>
              <article className="stat-card">
                <span className="feature-icon green">
                  <Trophy />
                </span>
                <div>
                  <strong>
                    {data.achievements.filter((a) => a.unlocked).length}
                  </strong>
                  <span>Logros desbloqueados</span>
                </div>
              </article>
            </section>
            <section className="two-columns">
              <article className="card">
                <div className="section-heading">
                  <h2>Tu camino, paso a paso</h2>
                  <Link to="/mapa">Ver mapa →</Link>
                </div>
                <p>{progress}% de la muestra introductoria completada</p>
                <Progress value={progress} />
                <p className="muted">
                  No representa dominio ni certificación de un nivel CEFR.
                </p>
              </article>
              <article className="card">
                <div className="section-heading">
                  <h2>
                    <Target size={20} /> Misión del día
                  </h2>
                  <Link to="/misiones">Ver misión →</Link>
                </div>
                {data.missions.filter(m => m.code === 'three_activities').map((m) => (
                  <div key={m.code}>
                    <p>{m.title}</p>
                    <Progress
                      value={Math.min(100, (m.progress * 100) / m.target)}
                      label={m.title}
                    />
                    <small>
                      {Math.min(m.progress, m.target)} / {m.target} · +
                      {m.reward} XP {m.rewarded && '· Recompensa obtenida'}
                    </small>
                  </div>
                ))}
              </article>
            </section>
            <div className="section-heading">
              <h2>Aprende jugando</h2>
              <Gamepad2 className="muted" />
            </div>
            <section className="game-grid">
              {gameCatalog.map(([kind, title, text], index) => {
                const color = ['blue', 'gold', 'green'][index % 3]
                const lesson = data.lessons.find(
                  (l) => l.kind === kind && !l.locked,
                )
                return (
                  <article className="card game-card" key={kind}>
                    <span className={`feature-icon ${color}`}>
                      <Gamepad2 />
                    </span>
                    <h3>{title}</h3>
                    <p>{text}</p>
                    {lesson ? (
                      <Link
                        className="text-link"
                        to={`/lecciones/${lesson.id}`}
                      >
                        Practicar →
                      </Link>
                    ) : (
                      <p className="muted">
                        Completa las lecciones previas para desbloquearlo.
                      </p>
                    )}
                  </article>
                )
              })}
            </section>
          </>
        )}
      </ResourceState>
    </>
  )
}
