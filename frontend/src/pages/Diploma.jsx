import { useState } from 'react'
import { Link } from 'react-router-dom'
import useResource from '../lib/useResource.js'
import { api, errorMessage } from '../lib/api.js'
import { PageTitle, ResourceState, Notice, Progress } from '../components/UI.jsx'
import DiplomaDocument from '../components/DiplomaDocument.jsx'

export default function Diploma() {
  const resource = useResource('/diploma')
  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const [issued, setIssued] = useState(null)
  const diploma = issued || resource.data?.diploma
  async function issue(event) {
    event.preventDefault(); setBusy(true); setError('')
    try { const { data } = await api.post('/diploma', { name }); setIssued(data.data) }
    catch (e) { setError(errorMessage(e)) }
    finally { setBusy(false) }
  }
  return <div className="certificate-page">
    <PageTitle eyebrow="TU RECONOCIMIENTO" title="Cada paso cuenta">Un diploma simbólico para celebrar la finalización de las actividades del curso. No es una acreditación CEFR ni un título oficial.</PageTitle>
    {error && <Notice>{error}</Notice>}
    {!diploma && <ResourceState resource={resource}>{resource.data && (resource.data.eligible ? <section className="card diploma-controls">
      <h2>¡Completaste todas las actividades!</h2><p>Escribe el nombre que deseas en tu diploma. Revísalo antes de confirmar: quedará guardado con la fecha y las actividades completadas y no se cambia automáticamente.</p>
      <form onSubmit={issue}><label htmlFor="diploma-name">Nombre para el diploma</label><input id="diploma-name" autoComplete="name" maxLength={120} minLength={2} required value={name} onChange={e => setName(e.target.value)} /><button className="btn" disabled={busy || !name.trim()}>{busy ? 'Emitiendo…' : 'Confirmar nombre y emitir diploma'}</button></form>
    </section> : <section className="card diploma-controls"><h2>Tu diploma te espera al terminar</h2><p>Has aprobado {resource.data.completed} de {resource.data.total} actividades publicadas. Debes completar las actividades de los seis niveles; aprobar un examen de acceso solo abre contenido.</p><Progress value={resource.data.total ? resource.data.completed / resource.data.total * 100 : 0} /><p>Los proyectos de escritura son práctica privada de autoevaluación, no una calificación docente incluida en este reconocimiento.</p><Link className="btn" to="/mapa">Continuar aprendiendo</Link></section>)}</ResourceState>}
    {diploma && <><div className="diploma-controls"><p>Tu diploma está guardado. Puedes imprimirlo o elegir «Guardar como PDF» en el diálogo de impresión del navegador. No se envía tu nombre a un servicio externo.</p><button className="btn" onClick={() => window.print()}>Imprimir o guardar PDF</button></div><DiplomaDocument diploma={diploma} /></>}
  </div>
}
