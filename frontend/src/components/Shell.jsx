import { useState } from 'react'
import {
  BookOpen,
  Compass,
  Flame,
  LayoutDashboard,
  LogOut,
  Medal,
  Shield,
  Target,
  TrendingUp,
  UserRound,
  GraduationCap,
} from 'lucide-react'
import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../auth/context.js'
import { errorMessage } from '../lib/api.js'
import { Notice } from './UI.jsx'
import PWAControls from './PWAControls.jsx'

export default function Shell() {
  const { user, logout } = useAuth()
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const links = [
    ['/dashboard', LayoutDashboard, 'Mi aventura'],
    ['/mapa', Compass, 'Aprender'],
    ['/misiones', Target, 'Misiones'],
    ['/logros', Medal, 'Logros'],
    ['/estadisticas', TrendingUp, 'Progreso'],
    ['/diploma', GraduationCap, 'Diploma'],
    ['/perfil', UserRound, 'Mi perfil'],
  ]
  if (user.role === 'ADMIN') links.push(['/admin', Shield, 'Administración'])
  async function signOut() {
    setBusy(true)
    setError('')
    try {
      await logout()
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  return (
    <div className="app-layout">
      <aside className="sidebar">
        <NavLink to="/dashboard" className="brand">
          <BookOpen aria-hidden="true" />
          <span>
            English<span>Quest</span>
          </span>
        </NavLink>
        <p className="sidebar-caption">UN POCO CADA DÍA. UN GRAN VIAJE.</p>
        <nav aria-label="Navegación principal">
          {links.map(([path, Icon, text]) => (
            <NavLink
              key={path}
              to={path}
              className={({ isActive }) =>
                isActive ? 'nav-item active' : 'nav-item'
              }
            >
              <Icon size={20} aria-hidden="true" />
              <span>{text}</span>
            </NavLink>
          ))}
        </nav>
        <div className="sidebar-bottom">
          <Flame size={24} />
          <p>
            Tu constancia importa.
            <br />
            <small>Sin prisa, sin comparaciones.</small>
          </p>
        </div>
        <button className="nav-item logout" onClick={signOut} disabled={busy}>
          <LogOut size={19} />
          {busy ? 'Cerrando…' : 'Cerrar sesión'}
        </button>
      </aside>
      <div className="workspace">
        <header className="topbar">
          <span className="pill">Prototipo educativo · práctica guiada</span>
          <PWAControls />
          <NavLink to="/perfil" className="user-chip">
            <span className="avatar">
              {user.first_name.slice(0, 1).toUpperCase()}
            </span>
            {user.first_name}
          </NavLink>
        </header>
        <main className="page-content">
          {error && <Notice>{error}</Notice>}
          <Outlet />
        </main>
        <footer className="app-footer">
          EnglishQuest · Aprende. Practica. Avanza.{' '}
          <NavLink to="/estado">Estado del servicio</NavLink>
        </footer>
      </div>
    </div>
  )
}
