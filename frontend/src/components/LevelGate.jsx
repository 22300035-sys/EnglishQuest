import { Link } from 'react-router-dom'
import { LockKeyhole } from 'lucide-react'

export default function LevelGate({ level }) {
  return <section className="card level-gate">
    <LockKeyhole size={32} aria-hidden="true" />
    <h2>{level.code} todavía está bloqueado</h2>
    <p>Completa todas las actividades publicadas de {level.previous_level} o aprueba el examen de acceso a {level.code}. Abrir un nivel no completa sus actividades.</p>
    <div className="hero-actions"><Link className="btn secondary" to={`/mapa?level=${level.previous_level}`}>Continuar {level.previous_level}</Link>
      {level.can_take_exam ? <Link className="btn" to={`/acceso/${level.code}`}>Examen de acceso a {level.code}</Link> : <p className="fine-print">Primero abre {level.previous_level} para poder presentar el examen siguiente.</p>}
    </div>
  </section>
}
