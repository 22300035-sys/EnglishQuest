# Incremento 04 — Rutas guiadas de A1 a C2

Fecha: 30 de octubre de 2026.

## Entrega

- Seis rutas publicadas A1, A2, B1, B2, C1 y C2 con seis unidades nuevas por nivel.
- 180 lecciones nuevas y 360 preguntas: lectura, transformación gramatical,
  interacción contextual, escucha y repaso de unidad.
- 36 audios originales sintéticos nuevos, con archivos reales incluidos.
- 36 proyectos de producción escrita con rúbrica de autoevaluación y borrador
  privado persistente en PostgreSQL. Sin calificaciones automáticas ficticias.
- Conteos de unidades/actividades en niveles y selección de las seis rutas.
- Migración 8 aditiva; contenido anterior, cuentas e intentos preservados.
- Catálogo total semilla: 203 actividades, 41 unidades y 42 audios.

## Verificación ejecutada

- `go test -race -count=1 ./...`: aprobado con una base PostgreSQL nueva de pruebas,
  no la base de demostración ni la instancia preexistente del sistema.
- Integración nueva: valida las 180 lecciones y sus soluciones; recorre todos los
  pasos publicados de las seis rutas mediante start/snapshot/submit, verifica
  desbloqueos, persistencia de progreso y logro de ruta sin duplicación.
- Portfolio: guardado, recuperación, límites, aislamiento entre estudiantes y
  ausencia de XP ficticia probados con PostgreSQL.
- `go vet ./...` y `go build ./...`: aprobados.
- Frontend: lint, 22 pruebas de componentes y compilación aprobados, incluidos
  guardado de proyecto y manejo honesto de un fallo de conexión.
- Python: ocho pruebas de herramientas/contenido aprobadas; compilación de la
  migración reproducible y contenido/audio no vacío verificados.
- Playwright: **6/6 aprobadas**, sin omisiones, en 2,4 minutos. Incluye una primera
  unidad completa de B1, B2, C1 y C2, reproducción real de audio, repaso, desbloqueo
  de la unidad siguiente, proyecto C2 guardado y recuperado tras recargar, además
  de las pruebas anteriores de estudiante, administración, juegos y PWA.

## Límites

Disponible como ruta de práctica de seis niveles, **no como curso CEFR exhaustivo
validado ni certificación**. Los repasos recuperan material entrenado; el portfolio
requiere revisión humana para evaluar la escritura libre. No se registra ni se
califica habla espontánea. La revisión docente y el resto de pendientes de
producto/producción continúan según `CURRICULUM.md`.

La migración se aplicó solo en desarrollo local. No hubo cambios de producción,
servicios de pago, commits, pushes ni despliegue Render.
