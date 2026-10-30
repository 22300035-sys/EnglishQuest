import { useState } from 'react'
import { BookOpen, ArrowLeft, Eye, EyeOff } from 'lucide-react'
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/context.js'
import { api, errorMessage } from '../lib/api.js'
import { Notice } from '../components/UI.jsx'

export default function Auth({ register = false }) {
  const { user, login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [form, setForm] = useState({
    name: '',
    email: location.state?.email || '',
    password: '',
  })
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  if (user) return <Navigate to="/dashboard" replace />
  async function submit(event) {
    event.preventDefault()
    setError('')
    if (
      register &&
      (new TextEncoder().encode(form.password).length < 10 ||
        new TextEncoder().encode(form.password).length > 72)
    ) {
      setError('La contraseña debe tener entre 10 y 72 bytes.')
      return
    }
    setBusy(true)
    try {
      if (register) {
        await api.post('/auth/register', {
          ...form,
          timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
        })
        navigate('/login', {
          replace: true,
          state: { registered: true, email: form.email },
        })
      } else {
        await login(form.email, form.password)
        navigate(location.state?.from || '/dashboard', { replace: true })
      }
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  const update = (event) =>
    setForm({ ...form, [event.target.name]: event.target.value })
  return (
    <main className="auth-page">
      <Link to="/" className="back-link">
        <ArrowLeft size={18} /> Volver al inicio
      </Link>
      <section className="auth-card">
        <span className="auth-logo">
          <BookOpen />
        </span>
        <span className="eyebrow">ENGLISHQUEST</span>
        <h1>{register ? 'Tu aventura empieza aquí' : '¡Qué bueno verte!'}</h1>
        <p>
          {register
            ? 'Crea tu cuenta y da el primer paso.'
            : 'Continúa donde lo dejaste.'}
        </p>
        {location.state?.registered && !register && (
          <p className="success-note" role="status">
            Cuenta creada. Inicia sesión para comenzar.
          </p>
        )}
        {error && <Notice>{error}</Notice>}
        <form onSubmit={submit}>
          {register && (
            <label>
              Tu nombre
              <input
                name="name"
                value={form.name}
                onChange={update}
                autoComplete="given-name"
                required
                maxLength={80}
              />
            </label>
          )}
          <label>
            Correo electrónico
            <input
              name="email"
              type="email"
              value={form.email}
              onChange={update}
              autoComplete="email"
              required
              maxLength={254}
            />
          </label>
          <label>
            Contraseña
            <div className="password-field">
              <input
                name="password"
                type={showPassword ? 'text' : 'password'}
                value={form.password}
                onChange={update}
                autoComplete={register ? 'new-password' : 'current-password'}
                required
                minLength={register ? 10 : 1}
                maxLength={72}
              />
              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                aria-label={
                  showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'
                }
              >
                {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
              </button>
            </div>
          </label>
          {register && (
            <small>Entre 10 y 72 bytes. Usa una contraseña única.</small>
          )}
          <button className="btn full" disabled={busy}>
            {busy
              ? 'Un momento…'
              : register
                ? 'Crear mi cuenta'
                : 'Entrar a mi aventura'}
          </button>
        </form>
        <p className="auth-switch">
          {register ? '¿Ya tienes una cuenta?' : '¿Primera vez aquí?'}{' '}
          <Link to={register ? '/login' : '/registro'}>
            {register ? 'Inicia sesión' : 'Crea tu cuenta'}
          </Link>
        </p>
        {!register && (
          <p className="fine-print">
            Recuperación de contraseña por correo: aún no disponible. No
            compartas tus credenciales.
          </p>
        )}
      </section>
      <p className="auth-footnote">
        Tu progreso se guarda en PostgreSQL, no en este navegador.
      </p>
    </main>
  )
}
