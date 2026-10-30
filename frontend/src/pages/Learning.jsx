import { Link, useParams, useSearchParams } from 'react-router-dom'
import { BookOpen, Check, LockKeyhole, Star } from 'lucide-react'
import { useAuth } from '../auth/context.js'
import useResource from '../lib/useResource.js'
import { PageTitle, ResourceState } from '../components/UI.jsx'
import UnitProject from '../components/UnitProject.jsx'
import LevelGate from '../components/LevelGate.jsx'

export default function Learning({ unit = false }) {
  const { user } = useAuth()
  const { id } = useParams()
  const [search, setSearch] = useSearchParams()
  const level = search.get('level') || user.cefr_level
  const resource = useResource('/lessons')
  const levels = useResource('/learning/levels')
  const lessons =
    resource.data?.filter((l) =>
      unit ? l.unit_id === id : l.level_code === level,
    ) || []
  const units = [...new Set(lessons.map((l) => l.unit_id))]
  const currentLevel = levels.data?.find(l => l.code === (unit ? lessons[0]?.level_code : level))
  const levelLocked = currentLevel?.available && !currentLevel.unlocked
  return (
    <>
      <PageTitle
        eyebrow="TU CAMINO DE DESCUBRIMIENTOS"
        title={
          unit ? lessons[0]?.unit_title || 'Unidad' : 'Un desafío a la vez'
        }
      >
        Lee cada microlección, practica y aprueba para abrir el siguiente paso.
      </PageTitle>
      {!unit && (
        <div className="level-tabs" aria-label="Nivel de inglés">
          {levels.data?.filter(l => l.available).map(({code,title,unlocked}) => (
            <button
              key={code}
              className={level === code ? 'selected' : ''}
              onClick={() => setSearch({ level: code })}
            >
              {!unlocked && <LockKeyhole size={14} aria-hidden="true" />} {code} <span>{title}</span>
            </button>
          ))}
          <Link to="/niveles">Todos los niveles →</Link>
        </div>
      )}
      {levels.error && <p className="notice" role="alert">{levels.error}<button className="text-button" onClick={levels.reload}>Reintentar niveles</button></p>}
      <ResourceState resource={resource}>
        {levelLocked && <LevelGate level={currentLevel} />}
        {!lessons.length && (
          <div className="card">
            No hay lecciones disponibles para esta selección.
          </div>
        )}
        {!levelLocked && units.map((unitID) => (
          <section key={unitID} className="map-unit">
            <div className="unit-heading">
              <span className="eyebrow">
                {lessons.find((l) => l.unit_id === unitID).level_code} ·
                UNIDAD DE PRÁCTICA
              </span>
              <h2>{lessons.find((l) => l.unit_id === unitID).unit_title}</h2>
              <p>{lessons.find(l => l.unit_id === unitID).unit_description}</p>
              {!unit && (
                <Link to={`/unidades/${unitID}`} className="text-link">
                  Ver unidad →
                </Link>
              )}
            </div>
            <ol className="learning-path">
              {lessons
                .filter((l) => l.unit_id === unitID)
                .map((lesson, i) => (
                  <li
                    key={lesson.id}
                    className={`path-step step-${i % 3} ${lesson.completed ? 'done' : ''} ${lesson.locked ? 'locked' : ''}`}
                  >
                    <div className="path-connection" />
                    {lesson.locked ? (
                      <span className="path-node" aria-label="Bloqueada">
                        <LockKeyhole />
                      </span>
                    ) : (
                      <Link
                        to={`/lecciones/${lesson.id}`}
                        className="path-node"
                        aria-label={`Abrir ${lesson.title}`}
                      >
                        {lesson.completed ? (
                          <Check size={30} />
                        ) : (
                          <BookOpen size={26} />
                        )}
                      </Link>
                    )}
                    <div>
                      <span className="step-caption">
                        {lesson.assessment
                          ? 'EVALUACIÓN DE UNIDAD'
                          : `PASO ${i + 1}`}
                      </span>
                      <h3>{lesson.title}</h3>
                      <p>
                        {lesson.locked
                          ? 'Aprueba el paso anterior para continuar.'
                          : lesson.objective}
                      </p>
                      {lesson.completed && (
                        <span
                          className="stars"
                          aria-label={`Mejor puntuación: ${lesson.best_score}%`}
                        >
                          {[1, 2, 3].map((n) => (
                            <Star
                              key={n}
                              size={16}
                              fill={
                                lesson.best_score >=
                                (n === 3 ? 100 : n === 2 ? 80 : 60)
                                  ? 'currentColor'
                                  : 'none'
                              }
                            />
                          ))}{' '}
                          {lesson.best_score}%
                        </span>
                      )}
                    </div>
                  </li>
                ))}
            </ol>
          </section>
        ))}
        {unit && !levelLocked && !lessons[0]?.level_locked && lessons[0]?.has_project && <UnitProject unitID={id} key={id} />}
      </ResourceState>
    </>
  )
}
