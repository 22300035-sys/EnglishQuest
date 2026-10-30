import { useEffect, useState } from 'react'
import { api, errorMessage } from '../lib/api.js'
import { Notice } from './UI.jsx'

export default function BoosterInventory({ boosters, serverTime, reload }) {
  const [base] = useState(() => ({ local: performance.now(), server: new Date(serverTime).getTime() }))
  const [now, setNow] = useState(base.server)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  useEffect(() => { const timer = setInterval(() => setNow(base.server + performance.now() - base.local), 1000); return () => clearInterval(timer) }, [base])
  const active = boosters.some(b => b.active && new Date(b.expires_at).getTime() > now)
  async function activate(id) {
    setBusy(true); setError('')
    try { await api.post(`/boosters/${id}/activate`); reload() } catch (e) { setError(errorMessage(e)); setBusy(false) }
  }
  return <section className="card booster-inventory"><div className="section-heading"><h2>Potenciadores de experiencia</h2><button className="text-button" disabled={busy} onClick={reload}>Actualizar inventario</button></div><p>Completa la misión de tres actividades nuevas para obtener doble XP base durante 15 minutos. No se acumulan ni multiplican bonos perfectos o de misiones. Solo benefician primeras aprobaciones y deben estar activos al validar.</p>{error && <Notice>{error}</Notice>}
    {!boosters.length && <p className="muted">Todavía no tienes potenciadores. Se guardan y otorgan desde el servidor, sin compras.</p>}
    <div className="booster-list">{boosters.map(booster => {
      const remaining = Math.max(0, Math.ceil((new Date(booster.expires_at).getTime() - now) / 1000))
      const running = booster.active && remaining > 0
      return <article key={booster.id}><div><h3>{booster.title}</h3><p>{booster.available ? 'Disponible para activar cuando estés listo.' : running ? `Activo · ${Math.floor(remaining / 60)}:${String(remaining % 60).padStart(2, '0')} restantes` : 'Utilizado o expirado'}</p></div>{booster.available && <button className="btn secondary small" disabled={busy || active} onClick={() => activate(booster.id)}>Activar doble XP</button>}</article>
    })}</div>
  </section>
}
