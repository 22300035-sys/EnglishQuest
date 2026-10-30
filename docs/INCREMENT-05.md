# Incremento 05 — Desbloqueos y diploma simbólico

Fecha: 30 de octubre de 2026.

## Política acordada

- A1, A2 y B1 abiertos inicialmente.
- B2, C1 y C2 se abren completando las actividades del anterior o aprobando
  el examen correspondiente. Para examinarse del siguiente debe estar abierto
  el anterior. Los accesos y el progreso preexistentes no se eliminan.
- Un examen toma 12 preguntas de un banco privado de 16 por nivel. Tiempo total
  de 20 minutos, 80% de aprobación (10/12) y 10 minutos de espera para reintentar.
- Intentos persistentes, opciones opacas, snapshots privados, expiración de
  servidor e idempotencia ante peticiones repetidas/concurrentes.
- No se entregan XP, rachas ni progreso de actividades por un examen de acceso.

## Diploma

Solo se habilita al aprobar todas las actividades publicadas de los seis niveles.
Se pide el nombre en ese momento, se valida, y se conserva la emisión con folio,
fecha y snapshot de las actividades. La confirmación es única e inmutable desde
la interfaz; volver a enviar el mismo nombre no duplica el diploma.

Diseño imprimible A4 horizontal y opción de usar «Guardar como PDF» del navegador.
No hay generador externo ni exposición pública de datos personales. El diploma
indica que es simbólico y no acredita dominio CEFR, habla o escritura libre.

## Verificación

- Backend con PostgreSQL nuevo: `go test -race -count=1 ./...`, `go vet ./...` y
  `go build ./...` aprobados.
- Tests: mitad inicial abierta, bloqueos por URL/API, acceso a proyectos protegido,
  prerrequisitos de exámenes, opciones/soluciones privadas, concurrencia,
  reanudación sin reiniciar plazo, expiración, reintento limitado y ausencia de XP.
- El recorrido completo de las seis rutas verifica desbloqueo por progreso y
  emisión única del diploma, nombre inmutable y privacidad entre cuentas.
- Frontend: lint, 27 pruebas y compilación aprobados; incluye nombre solicitado
  al final, fallo honesto de emisión, relectura, impresión y escape de texto.
- Playwright: **8/8 aprobadas**, sin omisiones, en 2,7 minutos. Se verifica la
  mitad inicial abierta, examen fallido/reintento bloqueado, acceso B2→C1→C2 por
  examen, curso completo por API real sin insertar progreso en SQL, nombre final,
  diploma persistente, y los recorridos anteriores de juegos, admin y PWA.
- Tras ajustar la impresión, se repitió solo la prueba del diploma: aprobada.
  El PDF real generado y `pdfinfo` confirman **una página A4 horizontal**;
  la prueba también comprueba el número de páginas para evitar regresiones.

La migración 9 es aditiva y se aplica solo en desarrollo local. Ningún cambio
de producción, servicio de pago, commit, push o despliegue fue realizado.
