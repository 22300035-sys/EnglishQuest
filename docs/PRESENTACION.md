# EnglishQuest — Guion de presentación

## Qué presentar

**EnglishQuest es un prototipo educativo full-stack funcional** para practicar
inglés mediante microlecciones y diez minijuegos. React presenta la actividad;
Go valida las respuestas y PostgreSQL conserva cuentas, resultados y recompensas.
No es un curso completo ni certifica un nivel CEFR.

## Preparación

1. Desde la raíz: `bash scripts/dev.sh` (dependencias instaladas en `frontend/`).
2. Abre **http://127.0.0.1:5173**. Usa siempre ese hostname para frontend/API local.
3. Registra una cuenta propia y entra. No hay contraseñas públicas ni usuarios
   administrativos predefinidos. No expongas `.local/dev.env` en la presentación.
4. Prueba el sonido y el recorrido antes de exponer. Los retos se desbloquean
   secuencialmente; no eludas el bloqueo con modificaciones de base de datos.
5. Para mostrar administración, registra otra cuenta y promuévela únicamente
   en la base local con el comando del README. Conserva la cuenta de estudiante.
   Evita registrar/iniciar sesión repetidamente justo antes de exponer: hay un
   límite de 15 intentos de autenticación por IP durante 15 minutos.
6. Para detener los servicios locales sin borrar los datos:
   `bash scripts/stop-dev.sh`.

## Demostración de 7–10 minutos

| Tiempo | Acción | Qué demuestra |
| --- | --- | --- |
| 0–1 min | Portada, registro e inicio de sesión | Formularios reales y cuenta persistente |
| 1–3 min | Aprender → A1 → primera microlección → Word Hunter → enviar | Respuestas interactivas y corrección en el servidor |
| 3–4 min | Resultado, explicaciones y siguiente actividad | Progresión, feedback y desbloqueos |
| 4–5 min | Dashboard y recarga de la página | El progreso proviene de PostgreSQL, no de datos ficticios |
| 5–6 min | Misiones, logros e inventario de potenciadores | Recompensas persistentes sin XP ilimitada por repetir |
| 6–8 min | Cuenta previamente preparada: escuchar audio, leer historia y reto cronometrado | Diez formatos y plazos controlados por Go |
| 8–9 min | Administración: editar un borrador y volver a abrirlo | Editor visual, catálogo y autorización del backend |
| 9–10 min | Arquitectura, pruebas y pendientes | Alcance técnico verificable y límites honestos |

No cambies la contraseña de la cuenta preparada durante la exposición salvo que
quieras demostrar el cierre de todas sus sesiones. No publiques los borradores
de pruebas ni uses cuentas con privilegios para el recorrido del estudiante.

## Arquitectura que puedes explicar

```text
React + Vite + Tailwind 4 + Router + Axios
                  ↓ REST /api/v1
Go + Gin → middleware de sesión/rol → handlers → servicios → repositorios
                                                                ↓
                                                      PostgreSQL + GORM
```

- Contraseñas bcrypt y sesiones opacas en cookies HttpOnly; sin localStorage.
- Intentos con snapshot: editar una pregunta no cambia un intento ya iniciado.
- Soluciones privadas hasta enviar; puntuación, plazos y XP calculados en Go.
- Transacciones y restricciones SQL evitan duplicar premios ante envíos repetidos.
- Siete migraciones SQL explícitas; scripts de respaldo/restauración y blueprint Render.

## Preguntas frecuentes

**¿Qué contenido incluye?** 203 actividades en 41 unidades A1–C2 y 36 proyectos
de escritura con autoevaluación. Las cuentas de prueba pueden añadir borradores,
por lo que el panel administrativo puede mostrar más registros.

**¿Los audios son reales?** Son 42 archivos WAV reproducibles, con texto original
y voz sintética eSpeak NG; no grabaciones de hablantes profesionales.

**¿Funciona sin internet?** La instalación PWA se ofrece cuando el navegador lo
permite en producción. Sin conexión informa del estado; no inventa progreso ni XP.
La demostración local necesita que la API y PostgreSQL estén funcionando.

**¿Está en Render?** Hay configuración y guía; no se ha realizado ni verificado
un despliegue público. No presentes localhost como una URL pública.

**¿Por qué algunos niveles están cerrados?** A1, A2 y B1 empiezan abiertos.
B2, C1 y C2 se abren por progreso del anterior o por un examen de acceso de
12 preguntas (20 minutos, mínimo 10 aciertos). Los accesos avanzados ya iniciados
antes del cambio se conservan. El examen abre contenido; no entrega XP ni diploma.

**¿Hay diploma?** Sí, simbólico: se obtiene al aprobar todas las actividades
publicadas de los seis niveles. Se pide el nombre al final y queda guardado con
fecha, folio y contenido completado. La impresión permite guardar un PDF A4
horizontal de una hoja. No es un título oficial ni una acreditación CEFR.

**¿Qué falta?** Cursos completos con revisión docente, recuperación de contraseña
por correo, subida de medios, repaso espaciado, evaluación de pronunciación y
escritura libre, auditoría amplia de accesibilidad y validación de producción.

## Evidencias

Consulta `docs/INCREMENT-05.md` para el resultado de la última verificación.
Las pruebas Playwright interactúan con API y PostgreSQL reales, sin interceptar
peticiones. Las capturas de las pruebas están en `frontend/test-results/` tras
ejecutarlas; ese directorio es temporal y no se incluye en Git.
