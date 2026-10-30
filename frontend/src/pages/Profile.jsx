import { useState } from 'react'
import { useAuth } from '../auth/context.js'
import { api, errorMessage } from '../lib/api.js'
import { Notice, PageTitle } from '../components/UI.jsx'
import useResource from '../lib/useResource.js'

export default function Profile() {
  const { user, setUser } = useAuth()
  const levels = useResource('/levels')
  const [name, setName] = useState(user.first_name)
  const [level, setLevel] = useState(user.cefr_level)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [passwords, setPasswords] = useState({ current_password: '', new_password: '', confirm: '' })
  const [passwordError, setPasswordError] = useState('')
  async function changePassword(event) {
    event.preventDefault(); setPasswordError('')
    if (passwords.new_password !== passwords.confirm) { setPasswordError('Las contraseñas nuevas no coinciden.'); return }
    setBusy(true)
    try { await api.post('/profile/password', { current_password: passwords.current_password, new_password: passwords.new_password }); setUser(null) } catch (e) { setPasswordError(errorMessage(e)) } finally { setBusy(false) }
  }
  async function save(event) {
    event.preventDefault()
    setBusy(true)
    setError('')
    setSuccess(false)
    try {
      const { data } = await api.patch('/profile', { name, cefr_level: level })
      setUser(data.data)
      setSuccess(true)
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  return (
    <>
      <PageTitle eyebrow="TU ESPACIO" title="Mi perfil">
        Personaliza tu nombre y tu nivel de práctica preferido.
      </PageTitle>
      <section className="card profile-card">
        <div className="profile-heading">
          <span className="avatar large">
            {user.first_name.slice(0, 1).toUpperCase()}
          </span>
          <div>
            <h2>{user.first_name}</h2>
            <p>{user.email}</p>
          </div>
        </div>
        {error && <Notice>{error}</Notice>}
        {success && (
          <p className="success-note" role="status">
            Cambios guardados.
          </p>
        )}
        <form onSubmit={save}>
          <label>
            Nombre
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              maxLength={80}
            />
          </label>
          <label>
            Nivel de práctica
            <select value={level} onChange={(e) => setLevel(e.target.value)}>
              <option value="" disabled>Selecciona un nivel disponible</option>
              {levels.data?.filter(l => l.available).map(l => <option key={l.code} value={l.code}>{l.code} — {l.title}</option>)}
            </select>
          </label>
          {levels.error && <Notice>{levels.error}<button type="button" className="text-button" onClick={levels.reload}>Reintentar niveles</button></Notice>}
          <button className="btn" disabled={busy || levels.loading}>
            {busy ? 'Guardando…' : 'Guardar cambios'}
          </button>
        </form>
        <dl className="profile-meta">
          <div>
            <dt>Zona horaria de aprendizaje</dt>
            <dd>{user.timezone}</dd>
          </div>
          <div>
            <dt>Tipo de cuenta</dt>
            <dd>{user.role === 'ADMIN' ? 'Administración' : 'Estudiante'}</dd>
          </div>
        </dl>
        <section className="password-change"><h2>Cambiar contraseña</h2><p className="fine-print">Por seguridad, se cerrarán todas tus sesiones al guardar.</p>{passwordError && <Notice>{passwordError}</Notice>}<form onSubmit={changePassword}>
          <label>Contraseña actual<input type="password" autoComplete="current-password" required maxLength={72} value={passwords.current_password} onChange={e => setPasswords({ ...passwords, current_password: e.target.value })} /></label>
          <label>Nueva contraseña<input type="password" autoComplete="new-password" required minLength={10} maxLength={72} value={passwords.new_password} onChange={e => setPasswords({ ...passwords, new_password: e.target.value })} /></label>
          <label>Confirmar nueva contraseña<input type="password" autoComplete="new-password" required minLength={10} maxLength={72} value={passwords.confirm} onChange={e => setPasswords({ ...passwords, confirm: e.target.value })} /></label>
          <button className="btn secondary" disabled={busy}>Actualizar contraseña y cerrar sesiones</button>
        </form></section>
        <p className="fine-print">
          La zona horaria se establece al registrarte para mantener coherentes
          las rachas y misiones. Cambio de correo y recuperación por
          correo: pendientes.
        </p>
      </section>
    </>
  )
}
