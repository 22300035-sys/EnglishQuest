import {
  ArrowRight,
  BookOpen,
  Check,
  Compass,
  Sparkles,
  Star,
} from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../auth/context.js'
import PWAControls from '../components/PWAControls.jsx'

export default function Welcome() {
  const { user } = useAuth()
  return (
    <main className="welcome">
      <header className="welcome-nav">
        <Link to="/" className="brand">
          <BookOpen />
          <span>
            English<span>Quest</span>
          </span>
        </Link>
        <Link
          className="btn secondary small"
          to={user ? '/dashboard' : '/login'}
        >
          {user ? 'Mi aventura' : 'Iniciar sesión'}
        </Link>
      </header>
      <section className="welcome-hero">
        <div className="hero-copy">
          <span className="eyebrow">
            <Sparkles size={16} /> PEQUEÑAS LECCIONES, GRANDES DESCUBRIMIENTOS
          </span>
          <h1>
            Tu inglés.
            <br />
            Tu ritmo.
            <br />
            <em>Tu aventura.</em>
          </h1>
          <p>
            Descubre palabras, construye ideas y supera pequeños desafíos.
            Aprender inglés empieza con un primer paso.
          </p>
          <div className="hero-actions">
            <Link to={user ? '/dashboard' : '/registro'} className="btn">
              Comenzar aventura <ArrowRight size={20} />
            </Link>
            <Link to="/niveles" className="btn secondary">
              Explorar niveles
            </Link>
          </div>
          <div className="hero-promises">
            <span>
              <Check /> Sin costo
            </span>
            <span>
              <Check /> Progreso guardado
            </span>
            <span>
              <Check /> A tu ritmo
            </span>
          </div>
        </div>
        <div
          className="adventure-scene"
          role="img"
          aria-label="Paisaje original de montañas, cielo azul y un camino de aprendizaje"
        >
          <div className="sun" />
          <div className="cloud cloud-one" />
          <div className="cloud cloud-two" />
          <div className="mountain mountain-back" />
          <div className="mountain mountain-front" />
          <div className="hill hill-back" />
          <div className="hill hill-front" />
          <div className="trail" />
          <div className="scene-token token-one">
            <BookOpen />
          </div>
          <div className="scene-token token-two">
            <Star fill="currentColor" />
          </div>
          <div className="scene-token token-three">
            <Compass />
          </div>
          <div className="word-bubble bubble-one">
            Hello! <span>👋</span>
          </div>
          <div className="word-bubble bubble-two">Let’s learn!</div>
          <div className="scene-label">
            <span className="avatar">A1</span>
            <div>
              El inicio de algo grande<small>Un desafío a la vez</small>
            </div>
          </div>
        </div>
      </section>
      <section className="welcome-features">
        <article>
          <span className="feature-icon blue">
            <BookOpen />
          </span>
          <h2>Aprende con sentido</h2>
          <p>Vocabulario y frases cotidianas en microlecciones.</p>
        </article>
        <article>
          <span className="feature-icon gold">
            <Star />
          </span>
          <h2>Practica jugando</h2>
          <p>Diez desafíos para escuchar, leer, escribir y construir.</p>
        </article>
        <article>
          <span className="feature-icon green">
            <Compass />
          </span>
          <h2>Ve tu progreso</h2>
          <p>Experiencia, logros y una ruta que avanza contigo.</p>
        </article>
      </section>
      <footer className="welcome-footer">
        <PWAControls />
        Rutas de práctica guiada A1–C2. No equivale a una acreditación CEFR.{' '}
        <Link to="/estado">Estado del servicio</Link>
      </footer>
    </main>
  )
}
