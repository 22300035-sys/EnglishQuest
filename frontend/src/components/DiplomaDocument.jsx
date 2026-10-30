export default function DiplomaDocument({ diploma }) {
  return <article className="diploma-document" aria-label="Diploma simbólico de finalización">
    <p className="diploma-brand">EnglishQuest</p>
    <p className="eyebrow">APRENDE · PRACTICA · AVANZA</p>
    <h2>Diploma simbólico<br />de finalización</h2>
    <p>Se entrega a</p><p className="diploma-name">{diploma.recipient_name}</p>
    <p>por completar las {diploma.activity_count} actividades publicadas de</p>
    <h3>{diploma.course_title}</h3>
    <p className="diploma-levels">A1 · A2 · B1 · B2 · C1 · C2</p>
    <p>Emitido el {new Date(diploma.issued_at).toLocaleDateString('es-MX', { dateStyle: 'long', timeZone: 'UTC' })}</p>
    <p className="fine-print">Reconocimiento simbólico de práctica en EnglishQuest. No es un título oficial ni certifica dominio CEFR, producción oral o escritura libre.</p>
    <p className="diploma-folio">Folio privado: {diploma.id}</p>
  </article>
}
