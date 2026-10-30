import { useState } from 'react'
import { api, errorMessage } from '../lib/api.js'
import { Notice } from './UI.jsx'

export default function AdminCatalog({ data, onSaved }) {
  const [section, setSection] = useState(null)
  const [form, setForm] = useState({})
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  function edit(kind, value) { setSection(kind); setForm({ ...value }); setError('') }
  function field(key, value) { setForm(prev => ({ ...prev, [key]: value })) }
  async function save(event) {
    event.preventDefault(); setBusy(true); setError('')
    try {
      if (section === 'level') await api.put(`/admin/levels/${form.code}`, { title: form.title, description: form.description, available: form.available })
      if (section === 'unit') await api[form.id ? 'put' : 'post'](form.id ? `/admin/units/${form.id}` : '/admin/units', { title: form.title, description: form.description, published: form.published, level_code: form.level_code })
      if (section === 'resource') await api.put('/admin/resources', { key: form.key, title: form.title, kind: form.kind, url: form.url, attribution: form.attribution, active: form.active })
      setSection(null); onSaved('Catálogo guardado. Los cambios se aplican a las nuevas sesiones de aprendizaje.')
    } catch (e) { setError(errorMessage(e)) } finally { setBusy(false) }
  }
  return <>
    <section className="card"><h2>Niveles CEFR</h2><p className="fine-print">Publicar un nivel requiere contenido disponible. Habilitarlo no certifica su cobertura académica.</p><div className="admin-list">{data.levels.map(level => <div key={level.code}><span><strong>{level.code} · {level.title}</strong><small>{level.available ? 'Disponible' : 'No disponible'}</small></span><button className="btn secondary small" onClick={() => edit('level', level)}>Editar nivel {level.code}</button></div>)}</div></section>
    <section className="card"><div className="section-heading"><h2>Unidades</h2><button className="btn small" onClick={() => edit('unit', { title: '', description: '', level_code: 'A1', published: false })}>Nueva unidad</button></div><div className="admin-list">{data.units.map(unit => <div key={unit.id}><span><strong>{unit.level_code} · {unit.title}</strong><small>Posición {unit.position} · {unit.published ? 'Publicada' : 'Oculta'}</small></span><button className="btn secondary small" onClick={() => edit('unit', unit)}>Editar unidad</button></div>)}</div></section>
    <section className="card"><div className="section-heading"><h2>Recursos educativos</h2><button className="btn small" onClick={() => edit('resource', { key: '', title: '', kind: 'image', url: '', attribution: '', active: true, isNew: true })}>Nuevo recurso</button></div><p className="fine-print">Archivos propios incluidos en el Static Site o enlaces HTTPS autorizados por el docente. No se suben archivos al disco efímero del backend ni se garantiza disponibilidad de enlaces externos.</p><div className="admin-list">{data.resources.map(resource => <div key={resource.key}><span><strong>{resource.title}</strong><small>{resource.kind} · {resource.key} · {resource.active ? 'Activo' : 'Inactivo'}</small></span><button className="btn secondary small" onClick={() => edit('resource', resource)}>Editar recurso</button></div>)}</div></section>
    {section && <section className="card"><h2>{section === 'level' ? `Editar nivel ${form.code}` : section === 'unit' ? form.id ? 'Editar unidad' : 'Crear unidad' : 'Configurar recurso'}</h2>{error && <Notice>{error}</Notice>}<form onSubmit={save}>
      {section === 'unit' && <label>Nivel de la unidad<select value={form.level_code} disabled={!!form.id} onChange={e => field('level_code', e.target.value)}>{data.levels.map(level => <option value={level.code} key={level.code}>{level.code} · {level.title}</option>)}</select></label>}
      {section === 'resource' && <><label>Clave del recurso<input required pattern="[a-z][a-z0-9_-]{0,39}" maxLength={40} readOnly={!form.isNew} value={form.key} onChange={e => field('key', e.target.value)} /></label><label>Tipo de recurso<select value={form.kind} disabled={!form.isNew} onChange={e => field('kind', e.target.value)}><option value="image">Imagen</option><option value="audio">Audio</option></select></label></>}
      <label>Título del catálogo<input value={form.title} required maxLength={section === 'level' ? 80 : 150} onChange={e => field('title', e.target.value)} /></label>
      {section !== 'resource' && <label>Descripción<textarea value={form.description || ''} maxLength={2000} onChange={e => field('description', e.target.value)} /></label>}
      {section === 'resource' && <><label>URL o ruta del archivo<input value={form.url} required maxLength={500} placeholder="https://... o /images/archivo.svg" onChange={e => field('url', e.target.value)} /></label><label>Autoría, licencia o atribución<textarea value={form.attribution} required minLength={3} maxLength={2000} onChange={e => field('attribution', e.target.value)} /></label></>}
      <label className="checkbox-label"><input type="checkbox" checked={section === 'level' ? form.available : section === 'unit' ? form.published : form.active} onChange={e => field(section === 'level' ? 'available' : section === 'unit' ? 'published' : 'active', e.target.checked)} />{section === 'level' ? 'Habilitar nivel' : section === 'unit' ? 'Publicar unidad' : 'Recurso activo'}</label>
      <div className="hero-actions"><button className="btn" disabled={busy}>{busy ? 'Guardando…' : 'Guardar catálogo'}</button><button type="button" className="btn secondary" onClick={() => setSection(null)}>Cancelar</button></div>
    </form></section>}
  </>
}
