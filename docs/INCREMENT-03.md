# Incremento 03 — Práctica ampliada y preparación de presentación

Fecha: 30 de octubre de 2026.

## Entregado

- Diez tipos de juego completos: Word Hunter, Sentence Builder, Memory Match,
  Sound Detective, Grammar Rescue, Spell Master, Conversation Quest,
  Story Explorer, Word Rain y English Rush.
- 23 microlecciones semilla en cinco unidades A1/A2, con seis audios sintéticos.
- Plazos de servidor, intentos inmutables y corrección de respuestas escritas,
  letras repetidas, selecciones múltiples y retos mixtos.
- Misiones por habilidad, inventario de potenciadores de 15 minutos,
  activación única, expiración y ×2 solamente para XP base.
- Logros por práctica de palabras y ruta de nivel completada, sin afirmar
  certificación ni vocabulario necesariamente distinto.
- Cambio de contraseña con revocación de sesiones.
- Editor visual para diez formatos y administración de niveles, unidades
  y referencias de recursos. Corrección de etiquetas del editor al reabrir preguntas.
- Manifest PWA, instalación condicionada al navegador y fallback offline honesto;
  sin caché de cuentas, API ni rutas privadas.
- Guion de exposición: `docs/PRESENTACION.md`; README y guías actualizados.

## Verificación

Ejecución final completa de Playwright: **5 de 5 pruebas aprobadas**, sin omisiones,
en 25,2 segundos, con administrador local autorizado y vista previa PWA activos.

- Backend: `go test -race -count=1 ./...` con PostgreSQL local en una base nueva
  separada de la demostración; aprobado. `go vet ./...` y `go build ./...`: aprobados.
- Frontend: lint, 20 pruebas de componentes y compilación de producción aprobados.
- Herramientas de respaldo: cuatro pruebas Python aprobadas.
- Navegador: recorrido del estudiante, siete juegos adicionales con audio real,
  temporizador persistente tras recarga, potenciador, cambio de contraseña,
  niveles no disponibles y fallback PWA aprobados.
- Administración: creación/edición persistente, editor visual, unidad, recurso y
  rechazo de habilitar B1 sin contenido publicado aprobados tras corregir etiquetas.

Los tests agregan cuentas y borradores a la base **local**, no a producción.
Las ejecuciones consecutivas comparten el límite de autenticación por IP (15
intentos en 15 minutos). Un `429` por reejecutar la suite no debe resolverse
desactivando esa protección: espera el plazo o reinicia solo la API de desarrollo
antes de otra ejecución completa. No realices ese procedimiento en producción.
La PWA se verifica contra `vite preview` en 4173; el desarrollo usa 5173 y no
registra service worker. Resultados históricos del incremento 02 describen una
versión anterior de tres juegos y no el alcance de esta entrega.

## Límites y pendientes

Esta entrega es presentable como **prototipo funcional**, no como producto
educativo completo. Faltan cursos íntegros revisados por docentes, recuperación
por correo, subida de archivos, repaso espaciado, evaluación de producción oral
y escritura libre, auditoría amplia de accesibilidad y validación/despliegue Render.

No se han realizado commits, pushes, despliegues externos ni migraciones de
producción. La instancia PostgreSQL preexistente del sistema no se usa ni modifica.
