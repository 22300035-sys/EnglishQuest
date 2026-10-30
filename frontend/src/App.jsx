import {
  Navigate,
  Outlet,
  Route,
  Routes,
  useLocation,
  Link,
} from 'react-router-dom'
import { useAuth } from './auth/context.js'
import { Loading } from './components/UI.jsx'
import Shell from './components/Shell.jsx'
import Welcome from './pages/Welcome.jsx'
import Auth from './pages/Auth.jsx'
import Levels from './pages/Levels.jsx'
import Dashboard from './pages/Dashboard.jsx'
import Learning from './pages/Learning.jsx'
import Lesson from './pages/Lesson.jsx'
import Exercise from './pages/Exercise.jsx'
import Results from './pages/Results.jsx'
import ProgressPages from './pages/ProgressPages.jsx'
import Profile from './pages/Profile.jsx'
import Admin from './pages/Admin.jsx'
import ServiceStatus from './pages/ServiceStatus.jsx'
import AccessExam, { AccessExamStart } from './pages/AccessExam.jsx'
import Diploma from './pages/Diploma.jsx'

function RequireAuth() {
  const { user, loading } = useAuth()
  const location = useLocation()
  if (loading) return <Loading />
  return user ? (
    <Outlet />
  ) : (
    <Navigate
      to="/login"
      replace
      state={{ from: location.pathname + location.search }}
    />
  )
}

export default function App() {
  const location = useLocation()
  return (
    <Routes>
      <Route path="/" element={<Welcome />} />
      <Route path="/login" element={<Auth key="login" />} />
      <Route path="/registro" element={<Auth key="register" register />} />
      <Route path="/niveles" element={<Levels />} />
      <Route path="/estado" element={<ServiceStatus />} />
      <Route element={<RequireAuth />}>
        <Route element={<Shell />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/mapa" element={<Learning />} />
          <Route path="/unidades/:id" element={<Learning unit />} />
          <Route path="/lecciones/:id" element={<Lesson />} />
          <Route path="/ejercicios/:id" element={<Exercise />} />
          <Route path="/resultados/:id" element={<Results />} />
          <Route
            path="/logros"
            element={<ProgressPages kind="achievements" />}
          />
          <Route path="/misiones" element={<ProgressPages kind="missions" />} />
          <Route
            path="/estadisticas"
            element={<ProgressPages kind="statistics" />}
          />
          <Route path="/perfil" element={<Profile />} />
          <Route path="/acceso/:code" element={<AccessExamStart />} />
          <Route path="/examen-acceso/:id" element={<AccessExam key={location.pathname} />} />
          <Route path="/diploma" element={<Diploma />} />
          <Route path="/admin" element={<Admin />} />
        </Route>
      </Route>
      <Route
        path="*"
        element={
          <main className="public-page">
            <h1>Página no encontrada</h1>
            <Link className="btn" to="/">
              Volver al inicio
            </Link>
          </main>
        }
      />
    </Routes>
  )
}
