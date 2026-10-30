import { Link } from 'react-router-dom'
import { LockKeyhole, ArrowRight } from 'lucide-react'
import useResource from '../lib/useResource.js'
import { PageTitle, ResourceState } from '../components/UI.jsx'
import { useAuth } from '../auth/context.js'

export default function Levels() {
  const { user } = useAuth()
  const resource = useResource(user ? '/learning/levels' : '/levels')
  return (
    <div className="public-page">
      <Link className="text-link" to={user ? '/dashboard' : '/'}>
        ← Volver
      </Link>
      <PageTitle
        eyebrow="TU RUTA DE APRENDIZAJE"
        title="Cada gran viaje tiene un comienzo"
      >
        Explora las rutas de práctica guiada desde A1 hasta C2. Cada ruta incluye
        lectura, gramática contextual, interacción escrita y escucha. Completar
        actividades no equivale a una acreditación CEFR.
        {' '}A1, A2 y B1 empiezan abiertos. Los siguientes se abren por progreso o examen.
      </PageTitle>
      <ResourceState resource={resource}>
        <div className="level-grid">
          {resource.data?.map((level) => (
            <article
              className={`card level-card ${!level.available ? 'unavailable' : ''}`}
              key={level.code}
            >
              <span className="level-code">{level.code}</span>
              <h2>{level.title}</h2>
              {level.available && <p className="fine-print">{level.unit_count} unidades · {level.lesson_count} actividades publicadas</p>}
              <p>
                {level.available
                  ? level.description || 'Práctica introductoria disponible'
                  : 'Contenido pendiente de desarrollo y validación'}
              </p>
              {level.available ? (
                <Link
                  className="btn"
                  to={user ? `/mapa?level=${level.code}` : '/login'}
                >
                  {level.unlocked ? 'Explorar' : 'Ver cómo desbloquear'} <ArrowRight size={18} />
                </Link>
              ) : (
                <span className="muted inline">
                  <LockKeyhole size={16} /> Próximamente
                </span>
              )}
              {level.available && !level.unlocked && <p className="muted inline"><LockKeyhole size={16} /> Completa el nivel anterior o aprueba su examen de acceso.</p>}
              {level.source === 'legacy' && <p className="fine-print">Conservas el acceso porque ya habías iniciado este nivel antes de incorporar los bloqueos.</p>}
            </article>
          ))}
        </div>
      </ResourceState>
    </div>
  )
}
