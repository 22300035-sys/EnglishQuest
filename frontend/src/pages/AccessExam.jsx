import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import useResource from '../lib/useResource.js'
import { api, errorMessage } from '../lib/api.js'
import { PageTitle, ResourceState, Notice } from '../components/UI.jsx'
import { AttemptClock, ChoiceQuestion } from '../components/AdvancedGames.jsx'

export function AccessExamStart() {
  const { code } = useParams()
  const resource = useResource('/learning/levels')
  const level = resource.data?.find(l => l.code === code)
  const navigate = useNavigate()
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  async function start() {
    setBusy(true); setError('')
    try { const { data } = await api.post(`/levels/${code}/access-exams`, {}); navigate(`/examen-acceso/${data.data.id}`) }
    catch (e) { setError(errorMessage(e)); setBusy(false) }
  }
  return <>
    <PageTitle eyebrow="ACCESO POR EXAMEN" title={`Abrir ${code}`}>Una alternativa al progreso del nivel anterior. No es una acreditación CEFR ni una prueba oficial de nivel.</PageTitle>
    {error && <Notice>{error}</Notice>}
    <ResourceState resource={resource}>{level ? <section className="card">
      <h2>Antes de empezar</h2><p>12 preguntas seleccionadas de un banco privado; 20 minutos controlados por el servidor. Necesitas al menos 80%: 10 respuestas correctas. Las preguntas comprueban gramática, comprensión y uso contextual escrito; no evalúan habla espontánea.</p>
      <p>Si apruebas, se abre {code}, pero no recibes XP, actividades completadas ni un diploma. Si no apruebas, puedes reintentar 10 minutos después de enviar. Un intento sin terminar se retoma con su plazo original.</p>
      <p>Las respuestas que estés eligiendo no se guardan al recargar; el examen y su plazo sí. Puedes enviar incompleto, pero las preguntas sin respuesta cuentan como incorrectas.</p>
      {level.unlocked ? <Link className="btn" to={`/mapa?level=${code}`}>Nivel abierto: ir a {code}</Link> : level.can_take_exam ? <button className="btn" disabled={busy} onClick={start}>{busy ? 'Preparando…' : 'Comenzar o retomar examen'}</button> : <p>Primero abre {level.previous_level || 'el nivel anterior'} para presentar este examen.</p>}
    </section> : !resource.loading && <Notice>Nivel no encontrado.</Notice>}</ResourceState>
  </>
}

export default function AccessExam() {
  const { id } = useParams()
  const resource = useResource(`/access-exams/${id}`)
  const [answers, setAnswers] = useState({})
  const [submitted, setSubmitted] = useState(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const exam = submitted || resource.data
  async function submit() {
    setBusy(true); setError('')
    try { const { data } = await api.post(`/access-exams/${id}/submit`, { answers }); setSubmitted(data.data) }
    catch (e) { setError(errorMessage(e)) }
    finally { setBusy(false) }
  }
  return <>
    <PageTitle eyebrow="EXAMEN DE ACCESO" title={exam ? `Acceso a ${exam.target_level}` : 'Tu examen'}>El resultado solo desbloquea contenido. No otorga XP ni completa el curso.</PageTitle>
    {error && <Notice>{error}</Notice>}
    <ResourceState resource={resource}>{exam && (exam.result ? <section className="card access-exam-result">
      <h2>{exam.result.passed ? `Nivel ${exam.target_level} desbloqueado` : exam.result.expired ? 'El plazo del examen terminó' : 'Todavía no alcanzas el umbral'}</h2>
      <p className="result-score">{exam.result.score}/100</p>
      <p>Umbral: 80%. Este resultado no certifica dominio CEFR.</p>
      {exam.result.passed ? <Link className="btn" to={`/mapa?level=${exam.target_level}`}>Practicar {exam.target_level}</Link> : <><p>Podrás iniciar otro examen a partir de {new Date(exam.retry_at).toLocaleString('es-MX')}.</p><Link className="btn secondary" to={`/acceso/${exam.target_level}`}>Volver al acceso</Link></>}
      <h3>Revisión de tus respuestas</h3>{exam.result.feedback.map((feedback, index) => <p key={feedback.question_id}>{index + 1}. {feedback.correct ? 'Correcta.' : 'Para revisar.'} {feedback.explanation}</p>)}
    </section> : <>
      <AttemptClock key={exam.id} expiresAt={exam.expires_at} serverTime={exam.server_time} />
      <p>{Object.keys(answers).length}/{exam.activity.questions.length} preguntas respondidas</p>
      {exam.activity.questions.map((question, index) => <section className="card question-card" key={question.id}><p className="eyebrow">PREGUNTA {index + 1}</p><p className="reading-passage">{question.passage}</p><h2>{question.prompt}</h2><ChoiceQuestion question={question} value={answers[question.id]} disabled={busy} onChange={value => setAnswers(previous => ({ ...previous, [question.id]: value }))} /></section>)}
      <button className="btn" disabled={busy} onClick={submit}>{busy ? 'Validando…' : 'Enviar examen'}</button>
    </>)}</ResourceState>
  </>
}
