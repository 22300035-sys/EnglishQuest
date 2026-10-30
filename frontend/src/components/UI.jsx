import { Link } from 'react-router-dom'

export function Notice({ children }) {
  return (
    <p className="notice" role="alert">
      {children}
    </p>
  )
}
export function Loading() {
  return (
    <div className="loading" role="status">
      <span className="loading-dot" />
      Preparando tu aventura…
    </div>
  )
}
export function ResourceState({ resource, children }) {
  if (resource.loading) return <Loading />
  if (resource.error)
    return (
      <section className="card">
        <Notice>{resource.error}</Notice>
        <button className="btn secondary" onClick={resource.reload}>
          Volver a intentar
        </button>
      </section>
    )
  return children
}
export function PageTitle({ eyebrow, title, children }) {
  return (
    <header className="page-title">
      <span className="eyebrow">{eyebrow}</span>
      <h1>{title}</h1>
      {children && <p>{children}</p>}
    </header>
  )
}
export function Progress({ value, label = 'Progreso' }) {
  return (
    <div
      className="progress"
      role="progressbar"
      aria-label={label}
      aria-valuenow={value}
      aria-valuemin={0}
      aria-valuemax={100}
    >
      <span style={{ width: `${value}%` }} />
    </div>
  )
}
export function Empty({ children }) {
  return (
    <div className="card empty">
      {children}
      <Link to="/mapa" className="text-link">
        Explorar lecciones →
      </Link>
    </div>
  )
}
