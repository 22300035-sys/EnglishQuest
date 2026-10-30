import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getServiceStatus } from '../lib/api.js'

export default function ServiceStatus() {
  const [state, setState] = useState('loading')
  const [attempt, setAttempt] = useState(0)

  useEffect(() => {
    const controller = new AbortController()
    getServiceStatus(controller.signal)
      .then(() => {
        if (!controller.signal.aborted) setState('ready')
      })
      .catch(() => {
        if (!controller.signal.aborted) setState('error')
      })
    return () => controller.abort()
  }, [attempt])

  return (
    <main className="min-h-screen bg-slate-950 text-white px-6 py-16">
      <div className="max-w-2xl mx-auto">
        <Link to="/" className="text-emerald-300 underline">
          ← EnglishQuest
        </Link>
        <h1 className="text-3xl font-bold mt-8">Estado de los servicios</h1>
        <p className="text-slate-300 mt-4">
          Comprobación real de la API y su conexión con PostgreSQL. No
          representa la disponibilidad de los cursos.
        </p>
        <section
          className="border border-slate-700 rounded-2xl p-6 my-8"
          aria-live="polite"
          aria-busy={state === 'loading'}
        >
          {state === 'loading' && <p role="status">Comprobando conexión…</p>}
          {state === 'ready' && (
            <p role="status" className="text-emerald-300">
              API disponible · PostgreSQL conectado
            </p>
          )}
          {state === 'error' && (
            <p role="alert" className="text-amber-200">
              No se pudo verificar la conexión. Comprueba que el backend esté
              iniciado y que la URL de la API y CORS estén configurados.
            </p>
          )}
          <button
            disabled={state === 'loading'}
            onClick={() => {
              setState('loading')
              setAttempt((n) => n + 1)
            }}
            className="mt-6 bg-emerald-700 rounded-xl px-5 py-3 disabled:opacity-50"
          >
            Volver a comprobar
          </button>
        </section>
        <h2 className="text-xl font-semibold">Disponibilidad del prototipo</h2>
        <p className="text-slate-300 mt-3">
          Puedes crear una cuenta, recorrer práctica guiada A1–C2 y guardar tus
          resultados. La PWA ofrece una página explicativa sin conexión, no práctica
          offline. La recuperación por correo y la acreditación de nivel no están disponibles.
        </p>
        <Link
          to="/registro"
          className="text-emerald-300 underline mt-4 inline-block"
        >
          Crear una cuenta →
        </Link>
      </div>
    </main>
  )
}
