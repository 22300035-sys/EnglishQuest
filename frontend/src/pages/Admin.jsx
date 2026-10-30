import { useState } from 'react'
import { useAuth } from '../auth/context.js'
import useResource from '../lib/useResource.js'
import { api, errorMessage } from '../lib/api.js'
import { Notice, PageTitle, ResourceState } from '../components/UI.jsx'
import AdminCatalog from '../components/AdminCatalog.jsx'
import QuestionEditor from '../components/QuestionEditor.jsx'
import { gameCatalog } from '../lib/gameCatalog.js'

const blank = {
  activity: {
    title: 'Word Hunter',
    kind: 'word_hunter',
    instructions: 'Read the English word and choose its picture.',
    published: true,
    assessment: false,
    questions: [
      {
        prompt: 'Find: apple',
        explanation: 'Apple means manzana.',
        options: [
          { id: 'a', label: 'Manzana', image: 'apple' },
          { id: 'b', label: 'Plátano', image: 'banana' },
        ],
        solution: 'a',
      },
    ],
  },
  unit_id: '',
  lesson_title: '',
  objective: '',
  content: '',
}

export default function Admin() {
  const { user } = useAuth()
  // The API independently checks the role on every request.
  if (user.role !== 'ADMIN')
    return (
      <div className="card">
        <h1>Acceso restringido</h1>
        <p>Tu cuenta no tiene permisos de administración.</p>
      </div>
    )
  return <AdminWorkspace />
}

function AdminWorkspace() {
  const resource = useResource('/admin/overview')
  const [editID, setEditID] = useState('')
  const [form, setForm] = useState(null)
  const [questions, setQuestions] = useState('')
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [busy, setBusy] = useState(false)
  const [visual, setVisual] = useState(true)
  let parsedQuestions = null
  try { const value = JSON.parse(questions); if (Array.isArray(value) && value.every(q => q && typeof q.prompt === 'string' && typeof q.explanation === 'string' && Array.isArray(q.options) && q.options.every(o => o && typeof o.id === 'string' && typeof o.label === 'string'))) parsedQuestions = value } catch { /* Preserve invalid advanced JSON for repair. */ }
  function create() {
    const draft = structuredClone(blank)
    draft.unit_id = resource.data.units[0]?.id || ''
    setEditID('')
    setForm(draft)
    setQuestions(JSON.stringify(draft.activity.questions, null, 2))
    setMessage('')
    setError('')
  }
  async function edit(item) {
    setBusy(true)
    setError('')
    setMessage('')
    try {
      const { data } = await api.get(`/admin/activities/${item.id}`)
      setEditID(item.id)
      setForm({
        activity: data.data,
        unit_id: '',
        lesson_title: item.lesson_title,
        objective: item.objective,
        content: item.content,
      })
      setQuestions(JSON.stringify(data.data.questions, null, 2))
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  async function save(event) {
    event.preventDefault()
    setError('')
    setMessage('')
    let parsed
    try {
      parsed = JSON.parse(questions)
    } catch {
      setError('Las preguntas deben ser un documento JSON válido.')
      return
    }
    setBusy(true)
    try {
      const payload = {
        ...form,
        activity: { ...form.activity, questions: parsed },
      }
      if (editID) await api.put(`/admin/activities/${editID}`, payload)
      else await api.post('/admin/lessons', payload)
      setForm(null)
      setMessage(
        'Contenido guardado. Los intentos en curso conservan su versión original.',
      )
      resource.reload()
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  async function toggle(item) {
    setBusy(true)
    setError('')
    try {
      await api.patch(`/admin/users/${item.id}`, { active: !item.is_active })
      resource.reload()
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  function field(key, value) {
    setForm((previous) => ({ ...previous, [key]: value }))
  }
  function activityField(key, value) {
    if (key === 'kind') {
      if (parsedQuestions) setQuestions(JSON.stringify(parsedQuestions.map(q => ({ ...q, response: '' })), null, 2))
      setForm(previous => ({ ...previous, activity: { ...previous.activity, kind: value, time_limit_seconds: ['word_rain', 'english_rush'].includes(value) ? 120 : 0, skill: gameCatalog.find(g => g[0] === value)[3] } }))
      return
    }
    setForm((previous) => ({
      ...previous,
      activity: { ...previous.activity, [key]: value },
    }))
  }
  return (
    <>
      <PageTitle eyebrow="GESTIÓN EDUCATIVA" title="Panel de administración">
        Contenido real, permisos del servidor y cambios auditados. No hay
        cuentas administrativas públicas.
      </PageTitle>
      {error && <Notice>{error}</Notice>}
      {message && (
        <p className="success-note" role="status">
          {message}
        </p>
      )}
      <ResourceState resource={resource}>
        {resource.data && (
          <>
            <section className="stats-grid">
              <article className="stat-card">
                <div>
                  <strong>{resource.data.stats.users}</strong>
                  <span>Usuarios registrados</span>
                </div>
              </article>
              <article className="stat-card">
                <div>
                  <strong>{resource.data.stats.attempts}</strong>
                  <span>Intentos completados</span>
                </div>
              </article>
              <article className="stat-card">
                <div>
                  <strong>{resource.data.stats.completions}</strong>
                  <span>Actividades aprobadas por estudiantes</span>
                </div>
              </article>
            </section>
            <AdminCatalog data={resource.data} onSaved={text => { setMessage(text); resource.reload() }} />
            <section className="card">
              <div className="section-heading">
                <h2>Lecciones y actividades</h2>
                <button className="btn small" onClick={create} disabled={busy}>
                  Nueva lección
                </button>
              </div>
              <div className="admin-list">
                {resource.data.activities.map((item) => (
                  <div key={item.id}>
                    <span>
                      <strong>
                        {item.level_code} · {item.lesson_title}
                      </strong>
                      <small>
                        {item.unit_title} · {item.kind} ·{' '}
                        {item.published ? 'Publicada' : 'Oculta'}
                      </small>
                    </span>
                    <button
                      className="btn secondary small"
                      onClick={() => edit(item)}
                      disabled={busy}
                    >
                      Editar
                    </button>
                  </div>
                ))}
              </div>
            </section>
            {form && (
              <section className="card admin-editor">
                <h2>{editID ? 'Editar contenido' : 'Crear lección'}</h2>
                <form onSubmit={save}>
                  {!editID && (
                    <label>
                      Unidad
                      <select
                        value={form.unit_id}
                        onChange={(e) => field('unit_id', e.target.value)}
                      >
                        {resource.data.units.map((u) => (
                          <option value={u.id} key={u.id}>
                            {u.level_code} · {u.title}
                          </option>
                        ))}
                      </select>
                    </label>
                  )}
                  <label>
                    Título de la lección
                    <input
                      required
                      maxLength={150}
                      value={form.lesson_title}
                      onChange={(e) => field('lesson_title', e.target.value)}
                    />
                  </label>
                  <label>
                    Objetivo de aprendizaje
                    <input
                      required
                      minLength={5}
                      value={form.objective}
                      onChange={(e) => field('objective', e.target.value)}
                    />
                  </label>
                  <label>
                    Microlección (en inglés)
                    <textarea
                      required
                      minLength={10}
                      value={form.content}
                      onChange={(e) => field('content', e.target.value)}
                      rows={4}
                    />
                  </label>
                  <div className="form-row">
                    <label>
                      Título de la actividad
                      <input
                        required
                        value={form.activity.title}
                        onChange={(e) => activityField('title', e.target.value)}
                      />
                    </label>
                    <label>
                      Minijuego
                      <select
                        value={form.activity.kind}
                        onChange={(e) => activityField('kind', e.target.value)}
                      >
                        {gameCatalog.map(([kind, title]) => <option value={kind} key={kind}>{title}</option>)}
                      </select>
                    </label>
                  </div>
                  <p className="fine-print">Cambiar el minijuego conserva las preguntas; revisa su formato y solución antes de guardar.</p>
                  <div className="form-row"><label>Habilidad principal<select value={form.activity.skill || 'vocabulary'} onChange={e => activityField('skill', e.target.value)}>{['vocabulary', 'grammar', 'listening', 'reading', 'writing', 'mixed'].map(skill => <option value={skill} key={skill}>{skill}</option>)}</select></label><label>Palabras objetivo practicadas<input type="number" min={0} max={100} value={form.activity.vocabulary_count || 0} onChange={e => activityField('vocabulary_count', Number(e.target.value))} /></label></div>
                  {['word_rain', 'english_rush'].includes(form.activity.kind) && <label>Plazo del intento (segundos)<input type="number" min={30} max={300} required value={form.activity.time_limit_seconds || 120} onChange={e => activityField('time_limit_seconds', Number(e.target.value))} /></label>}
                  <label>
                    Instrucciones
                    <textarea
                      required
                      minLength={5}
                      value={form.activity.instructions}
                      onChange={(e) =>
                        activityField('instructions', e.target.value)
                      }
                    />
                  </label>
                  <label className="checkbox-label">
                    <input
                      type="checkbox"
                      checked={form.activity.published}
                      onChange={(e) =>
                        activityField('published', e.target.checked)
                      }
                    />{' '}
                    Publicada
                  </label>
                  <label className="checkbox-label">
                    <input
                      type="checkbox"
                      checked={form.activity.assessment}
                      onChange={(e) =>
                        activityField('assessment', e.target.checked)
                      }
                    />{' '}
                    Evaluación (mínimo 80%)
                  </label>
                  <div className="hero-actions"><button type="button" className={`btn ${visual ? '' : 'secondary'} small`} onClick={() => setVisual(true)}>Editor visual</button><button type="button" className={`btn ${visual ? 'secondary' : ''} small`} onClick={() => setVisual(false)}>JSON avanzado</button></div>
                  {visual && parsedQuestions ? <QuestionEditor questions={parsedQuestions} kind={form.activity.kind} resources={resource.data.resources} onChange={value => setQuestions(JSON.stringify(value, null, 2))} /> : <label>
                    Preguntas, opciones y soluciones (JSON)
                    <textarea
                      className="code-editor"
                      rows={18}
                      value={questions}
                      onChange={(e) => setQuestions(e.target.value)}
                    />
                  </label>}
                  <p className="fine-print">
                    En el modo avanzado: solution es el ID de opción para Word Hunter; lista ordenada de
                    IDs para Sentence Builder; lista de parejas de IDs para
                    Memory Match. Cada opción debe tener un ID único. Imágenes
                    disponibles: apple, banana, orange, book, pen, school,
                    breakfast, night. Las soluciones solo son visibles a
                    administradores y después de completar un intento.
                  </p>
                  <div className="hero-actions">
                    <button className="btn" disabled={busy}>
                      {busy ? 'Guardando…' : 'Guardar contenido'}
                    </button>
                    <button
                      type="button"
                      className="btn secondary"
                      onClick={() => setForm(null)}
                    >
                      Cancelar
                    </button>
                  </div>
                </form>
              </section>
            )}
            <section className="card">
              <h2>Usuarios · últimos 100</h2>
              <div className="admin-list">
                {resource.data.users.map((item) => (
                  <div key={item.id}>
                    <span>
                      <strong>{item.first_name}</strong>
                      <small>
                        {item.email} · {item.role} ·{' '}
                        {item.is_active ? 'Activa' : 'Desactivada'}
                      </small>
                    </span>
                    {item.role === 'STUDENT' && (
                      <button
                        className="btn secondary small"
                        disabled={busy}
                        onClick={() => toggle(item)}
                      >
                        {item.is_active ? 'Desactivar' : 'Activar'}
                      </button>
                    )}
                  </div>
                ))}
              </div>
            </section>
            <p className="fine-print">
              Puedes gestionar niveles, unidades, recursos y contenido educativo. La promoción
              administrativa es una operación explícita del propietario desde
              CLI.
            </p>
          </>
        )}
      </ResourceState>
    </>
  )
}
