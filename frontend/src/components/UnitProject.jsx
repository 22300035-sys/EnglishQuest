import { useState } from 'react'
import useResource from '../lib/useResource.js'
import { api, errorMessage } from '../lib/api.js'
import { Notice, ResourceState } from './UI.jsx'

export default function UnitProject({ unitID }) {
  const resource = useResource(`/units/${unitID}/project`)
  return <ResourceState resource={resource}>{resource.data && <ProjectForm project={resource.data} unitID={unitID} />}</ResourceState>
}

function ProjectForm({ project, unitID }) {
  const [body, setBody] = useState(project.body)
  const [reflection, setReflection] = useState(project.reflection)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [saved, setSaved] = useState(false)
  async function save(event) {
    event.preventDefault(); setBusy(true); setError(''); setSaved(false)
    try { await api.put(`/units/${unitID}/project`, { body, reflection }); setSaved(true) }
    catch (e) { setError(errorMessage(e)) }
    finally { setBusy(false) }
  }
  const words = body.trim() ? body.trim().split(/\s+/).length : 0
  return <section className="card unit-project">
    <h2>Proyecto de producción y reflexión</h2>
    <p>Aplica lo aprendido con tus propias palabras. Se guarda de forma privada en PostgreSQL; no se califica automáticamente ni otorga XP. Puedes revisar el borrador con un docente. No se guarda mientras escribes: pulsa Guardar antes de salir.</p>
    <p className="reading-passage">{project.prompt}</p>
    <details><summary>Guía de autoevaluación</summary><p className="reading-passage">{project.rubric}</p></details>
    {error && <Notice>{error}</Notice>}
    <form onSubmit={save}>
      <label htmlFor="unit-project-body">Tu borrador en inglés</label>
      <textarea id="unit-project-body" rows={12} maxLength={12000} required value={body} onChange={e => { setBody(e.target.value); setSaved(false) }} />
      <p className="fine-print">{words} palabras · {body.length}/12000 caracteres · extensión sugerida, no nota de nivel</p>
      <label htmlFor="unit-project-reflection">Reflexión: qué revisaste y qué dudas tienes</label>
      <textarea id="unit-project-reflection" rows={4} maxLength={2000} value={reflection} onChange={e => { setReflection(e.target.value); setSaved(false) }} />
      <button className="btn" disabled={busy || !body.trim()}>{busy ? 'Guardando…' : 'Guardar proyecto'}</button>
      {saved && <p className="success-note" role="status">Proyecto guardado. No se ha asignado una calificación automática.</p>}
    </form>
  </section>
}
