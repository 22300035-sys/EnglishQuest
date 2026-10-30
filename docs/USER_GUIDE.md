# Cómo probar EnglishQuest

## Arranque local

Desde la raíz del repositorio, instala dependencias con `npm ci` dentro de
`frontend/` y ejecuta `bash scripts/dev.sh`. Abre **http://127.0.0.1:5173**.
El script crea su propio PostgreSQL en el puerto 55440; no usa ni modifica el
servidor que ya pudiera existir en el 5432. Los datos quedan en `.local/pgdata`,
y las credenciales privadas generadas en `.local/dev.env`, excluidos de Git.

## Recorrido del estudiante

1. Pulsa **Comenzar aventura**. Regístrate con nombre, correo y contraseña de
   entre 10 y 72 bytes. No hace falta un correo real para la demo: no se envían
   mensajes ni se verifica su propiedad todavía.
2. Inicia sesión. No hay cuentas de muestra ni contraseñas predefinidas.
3. En el dashboard, selecciona **Continuar aprendiendo** o abre **Aprender**.
4. Elige una ruta de A1 a C2. Lee la microlección y comienza el desafío.
5. Word Hunter: elige imágenes. Sentence Builder: ordena palabras, pudiendo
   quitarlas o reiniciar. Memory Match: revela tarjetas y propone parejas;
   se corrigen juntas al enviar. No se muestran aciertos inventados al girarlas.
6. Envía las respuestas. La API calcula puntuación, guarda el resultado, otorga
   las recompensas correspondientes y desbloquea el siguiente paso si apruebas.
7. Revisa explicaciones, mapa, misión, logros y estadísticas. Recarga la página
   o vuelve a iniciar sesión: el progreso sigue guardado en PostgreSQL.
8. Repite actividades para mejorar, no para acumular XP sin límite.

Después de las unidades iniciales encontrarás Sound Detective (escuchar audio),
Grammar Rescue (corregir frases), Spell Master (ordenar letras), Conversation
Quest (elegir una respuesta contextual), Story Explorer (leer y responder),
Word Rain (seleccionar palabras de una categoría) y English Rush (reto mixto).
Los dos últimos tienen plazo de servidor: recargar no reinicia el tiempo.
Si no hay audio disponible se informa del error, no se simula la escucha.

En **Misiones**, completa tres actividades nuevas durante tu día local para
ganar un potenciador de 15 minutos. Actívalo cuando vayas a practicar: duplica
solo la XP base de primeras aprobaciones mientras esté vigente al enviar.
No se acumula con otro ni multiplica bonos perfectos o de misiones.
En **Mi perfil** puedes cambiar tu contraseña; se cierran todas tus sesiones.

El catálogo semilla tiene 203 actividades en 41 unidades y rutas de A1 a C2.
Son rutas guiadas, no cursos CEFR exhaustivos validados ni certificación.
Consulta `CURRICULUM.md` para el temario y los límites de evaluación.
No hay mascotas, anuncios ni funciones de pago.

### Desbloqueo de niveles

Empiezas con **A1, A2 y B1** abiertos. Para abrir **B2, C1 o C2**, completa
todas las actividades publicadas del anterior o usa **Examen de acceso**.
El examen tiene 12 preguntas y 20 minutos; necesitas 10 aciertos. Primero debe
estar abierto el nivel anterior. Si no apruebas, espera 10 minutos desde el envío.
No da XP ni completa actividades. Se retoma con el mismo plazo si recargas,
pero las selecciones no enviadas se pierden. Si el plazo venció, envía el intento
para cerrarlo antes de volver a intentar. Los accesos avanzados que ya habías
iniciado antes de este cambio se conservan.

### Diploma al finalizar

Aprueba todas las actividades publicadas de los seis niveles. En el último
resultado aparecerá **Elegir nombre y obtener mi diploma**, o entra en **Diploma**
desde el menú. Escribe el nombre completo que deseas, revísalo y confirma.
El nombre queda registrado; no se cambia automáticamente después de emitirlo.
Pulsa **Imprimir o guardar PDF** y, para descargarlo, elige «Guardar como PDF»
en la impresión del navegador. El diploma es simbólico, no un título oficial
ni una certificación CEFR. Los exámenes de acceso no sustituyen la finalización
del curso, y los proyectos de escritura no reciben una nota docente automática.

En las unidades nuevas, **Ver unidad** también abre el proyecto de escritura:
lee la consigna y la rúbrica, escribe tu borrador y reflexión, y pulsa **Guardar
proyecto** antes de salir. No se guarda automáticamente. El texto es privado,
puede recuperarse al volver y no recibe una calificación ni XP inventadas.
La propuesta de ensayo oral es voluntaria y no graba el micrófono.

## Administración

Registra una cuenta propia. El propietario local puede promoverla explícitamente
siguiendo el comando `--promote-admin` de README. La cuenta verá Administración;
las API verifican el rol desde PostgreSQL, no desde la interfaz.

- Crea una lección en una unidad existente, o edita microlección, actividad,
  instrucciones, preguntas, opciones, explicaciones y soluciones.
- El editor visual permite editar preguntas, opciones, soluciones, historias y audio.
  También hay un modo JSON avanzado validado. Guarda contenido en inglés correcto.
- En Word Hunter, cada opción requiere una de las ilustraciones locales admitidas.
- Publicar/ocultar no borra intentos ni resultados. Los intentos ya iniciados
  conservan su versión original.
- Activa/desactiva estudiantes. Desactivar revoca sus sesiones; reactivar no
  restaura una sesión revocada.
- Puedes gestionar niveles, unidades y el catálogo de recursos de imagen/audio.
  Un nivel sin contenido publicado no se habilita. El catálogo referencia archivos:
  no es un servicio para subir medios. Asignar roles sigue siendo una operación CLI.

No promociones cuentas ni apliques migraciones en producción sin autorización
explícita del propietario. No compartas `.local/dev.env` ni archivos de respaldo.

## Instalación y uso sin conexión

En una compilación de producción (HTTPS o localhost), el navegador puede ofrecer
**Instalar EnglishQuest** si cumple sus requisitos. En desarrollo no se registra
el service worker. Sin conexión se muestra una página explicativa: no se guardan
respuestas pendientes, cuentas ni resultados en caché, y no se otorga XP offline.

## Pendientes reales después de esta entrega

1. **Valor educativo:** cursos completos revisados por un docente, repaso espaciado,
   evaluación de escritura libre y pronunciación.
2. **Recursos:** subida de medios con almacenamiento seguro, además del catálogo existente.
3. **Cuentas:** recuperación de contraseña con entrega segura de correo.
   Nunca simular envíos de recuperación.
4. **Accesibilidad:** auditoría de lectores de pantalla, teclado y más dispositivos.
5. **Producción:** validar cookies/CORS en los dominios reales de Render,
   revisar límites y privacidad, preparar backups y resolver la expiración de
   PostgreSQL gratuito a los 30 días sin contratar nada automáticamente.
