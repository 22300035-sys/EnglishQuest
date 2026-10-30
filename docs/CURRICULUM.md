# EnglishQuest — Rutas guiadas A1–C2

## Alcance implementado

Se han añadido **36 unidades, 180 lecciones, 360 preguntas y 36 proyectos de
producción escrita**, repartidos en seis rutas. Se conservan las 23 actividades
anteriores: el catálogo semilla total tiene **203 actividades en 41 unidades**.
Hay 42 audios locales de narración sintética, incluidos los seis anteriores.

Cada unidad nueva tiene cinco pasos:

1. **Read and discover:** explicación y lectura original con comprensión.
2. **Grammar in context:** dos transformaciones o respuestas controladas.
3. **Respond with purpose:** selección de una respuesta contextual/pragmática.
4. **Listen for meaning:** un audio distinto de la lectura y dos preguntas.
5. **Unit review:** recuperación de contenidos de la lectura y la interacción;
   es un repaso de reconocimiento, no una prueba independiente de dominio.

Desde **Ver unidad** se abre además el **proyecto de producción y reflexión**.
El estudiante escribe un borrador y una autoevaluación, pulsa Guardar y puede
recuperarlos al volver. El servidor no interpreta una autoevaluación como una
calificación docente, ni otorga XP por almacenar texto.

## Programa por nivel

| Nivel | Seis unidades nuevas | Foco de práctica |
| --- | --- | --- |
| A1 | People and introductions; Days and daily routines; Homes and neighbourhoods; Food and simple requests; Interests and abilities; Shopping and near-future plans | Identidad, be, presente simple, ubicación, cantidades, peticiones, can y going to |
| A2 | Past events and short narratives; Travel and arrangements; Choices and comparisons; Health, advice and obligations; Experiences and unfinished time; Plans, predictions and invitations | Pasado, arreglos futuros, comparativos, obligación/consejo, presente perfecto y primer condicional |
| B1 | Narratives and turning points; Learning strategies and advice; Work and practical communication; Community and consequences; News, claims and sources; Decisions and possible outcomes | Narración, contexto temporal, hipótesis, pasiva, relativas, estilo indirecto y decisiones justificadas |
| B2 | Developing a balanced argument; Evidence, likelihood and causation; Negotiating workable agreements; Reporting events and perspectives; Environment and policy trade-offs; Culture, interpretation and reviews | Concesión, incertidumbre, negociación, atribución, contrafactuales, distribución de costes e interpretación |
| C1 | Stance, hedging and calibrated claims; Register and diplomatic disagreement; Synthesising sources without flattening them; Implication and institutional language; Evaluation and counterarguments; Mediating technical information | Graduación de certeza, registro, síntesis, inferencia, condiciones necesarias/suficientes y mediación |
| C2 | Rhetorical framing and presupposition; Irony, understatement and perspective; Ambiguity and scope; Lexical precision and evaluative force; Conflicting accounts and interpretive frames; Faithful reformulation and final synthesis | Presuposiciones, ironía, alcance, fuerza léxica, marcos interpretativos y reformulación fiel |

Las unidades iniciales anteriores siguen antes de la ampliación en A1 y A2.
Por eso estos niveles tienen 43 y 40 actividades; B1–C2 tienen 30 cada uno.
A1, A2 y B1 empiezan abiertos. B2, C1 y C2 requieren completar las actividades
publicadas del nivel anterior o aprobar un examen de acceso. Para presentar el
examen siguiente basta tener abierto el nivel previo; no se salta directamente
de B1 a un examen C2. Dentro de cada ruta el desbloqueo también es secuencial.
Todo se verifica en Go, no solo en botones. Los accesos avanzados ya iniciados
antes de incorporar esta política se conservan, sin borrar progreso.

## Exámenes y diploma simbólico

Cada examen selecciona 12 preguntas de un banco de 16 para el nivel solicitado.
El umbral es 80% (al menos 10 aciertos), el plazo es 20 minutos y el reintento
requiere esperar 10 minutos tras enviar. Recargar conserva preguntas y plazo,
pero no las selecciones aún no enviadas. El examen solo abre el nivel: no
entrega XP, racha, actividades completadas ni acreditación CEFR.

Al aprobar todas las actividades publicadas de los seis niveles, aparece el
acceso al diploma desde Resultados y también está la sección **Diploma**.
Solo entonces se solicita el nombre elegido. Se guarda un reconocimiento
simbólico con folio privado, fecha y snapshot del contenido completado.
La emisión es única por cuenta y el nombre confirmado no cambia automáticamente.
Puede imprimirse o guardarse como PDF mediante el navegador, sin enviar el
nombre a servicios externos. Los proyectos privados de autoevaluación no se
presentan como una evaluación docente incluida en ese diploma.

## Producción, rúbricas y evaluación

Los proyectos sugieren 40–80 palabras en A1, 80–120 en A2, 120–180 en B1,
180–250 en B2, 250–350 en C1 y 300–450 en C2. Son rangos orientativos, no
umbrales que otorguen un nivel automáticamente. Se permite guardar borradores.

La rúbrica invita a revisar cuatro dimensiones: tarea/audiencia, organización,
recursos lingüísticos/registro y precisión/revisión. Cada proyecto incluye el
contexto original y el foco lingüístico de su unidad. También propone ensayo
oral voluntario, **sin capturar el micrófono ni evaluar pronunciación**.

Las respuestas controladas verifican un objetivo específico con normalización
de mayúsculas, espacios y puntuación final; no aceptan automáticamente todas las
paráfrasis válidas. La consigna lo advierte. El docente debe revisar la producción
abierta y los posibles equivalentes que el corrector cerrado no reconozca.

## Qué significa completar una ruta

- Las actividades automáticas conservan los umbrales de práctica 60% y repaso 80%.
- El progreso y el logro de ruta cuentan actividades publicadas aprobadas, no
  una evaluación docente de los proyectos ni el dominio espontáneo del idioma.
- Los proyectos no generan notas automáticas ni premios por volumen de texto.
- No se revocan logros históricos al ampliar un catálogo que antes tenía menos
  contenido. El progreso porcentual actual sí usa el contenido publicado actual.

**Las seis rutas están implementadas, pero no se afirma cobertura académica
exhaustiva del CEFR ni que completar 30 actividades convierta a alguien en C2.**
Para considerarlas cursos completos validados faltan revisión docente externa,
ampliación/validación de cobertura de descriptores, tareas sostenidas y variadas,
interacción oral real, pronunciación, evaluación moderada de escritura libre y
pruebas de transferencia con material no entrenado. No hay certificación oficial.

## Autoría, medios y mantenimiento

- Fuente privada: `backend/curriculum/course.json`. Contiene las soluciones;
  nunca debe copiarse a `frontend/public` ni devolverse por una API pública.
- Compilador: `python3 scripts/build-curriculum.py`; `--check` comprueba que la
  migración coincide con la fuente, sin modificar archivos ni base de datos.
- Migración aditiva: `000008_guided_curriculum.up.sql`. Añade datos y portfolio,
  sin borrar contenido, intentos, usuarios, cambios de administrador ni progreso.
- Audio: `python3 scripts/generate-audio.py`, con eSpeak NG instalado. Textos
  originales y voz sintética declarada; el aumento de velocidad no se presenta
  como evidencia de dificultad CEFR validada.
- Nunca se regenera una migración ya desplegada. Futuras correcciones de una
  instalación publicada requieren una migración nueva o edición administrativa.
- Los tests validan estructura, soluciones, recursos, desbloqueos, persistencia y
  privacidad; no sustituyen la revisión pedagógica humana.
