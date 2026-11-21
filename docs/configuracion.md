# Configuración por entorno

## Entornos

| Entorno | Datos | Propósito |
|---|---|---|
| local | sintéticos | desarrollo |
| CI | efímeros y sintéticos | validación automática |
| staging | sintéticos | pruebas integradas y smoke tests |
| producción | reales, restringidos | servicio operativo |

Cada entorno usa credenciales, claves, bases de datos y proyectos FCM
independientes. Nunca se reutilizan secretos entre entornos.

## Reglas

- La configuración se inyecta por variables de entorno o un gestor de secretos.
- La aplicación debe validar variables obligatorias al arrancar.
- No existen secretos ni valores de producción por defecto.
- `.env` y archivos de credenciales no se versionan.
- Los cambios de configuración se revisan y auditan como cambios de código.

## Variables mínimas

| Variable | Componente | Obligatoria | Descripción |
|---|---|---:|---|
| `APU_ENV` | API | sí | `local`, `ci`, `staging` o `production` |
| `APU_HTTP_ADDR` | API | sí | dirección y puerto de escucha |
| `APU_DATABASE_URL` | API | sí | conexión PostgreSQL |
| `APU_ALLOWED_ORIGIN` | API | sí | origen exacto autorizado para el dashboard |
| `APU_COOKIE_SECURE` | API | sí | exige cookies HTTPS; obligatorio en producción |
| `APU_SESSION_TTL` | API | sí | duración de la sesión humana (`12h` por defecto local) |
| `APU_MFA_ENCRYPTION_KEY` | API | producción | clave base64 de 32 bytes para cifrar secretos TOTP |
| `APU_PASSWORD_RESET_BASE_URL` | API | sí | origen absoluto del enlace de recuperación; HTTPS en producción |
| `APU_EXPOSE_DEVELOPMENT_TOKENS` | API | solo local | permite mostrar el token de recuperación para pruebas locales |
| `APU_SMTP_ADDRESS` | API | producción | servidor y puerto SMTP con STARTTLS |
| `APU_SMTP_USERNAME` | API | según proveedor | usuario SMTP |
| `APU_SMTP_PASSWORD` | API | según proveedor | contraseña SMTP; siempre desde gestor de secretos |
| `APU_SMTP_FROM` | API | producción | remitente válido de recuperación |
| `APU_FIELD_ENCRYPTION_KEY` | API | producción | clave de cifrado de campos |
| `APU_FCM_CREDENTIALS` | API | producción | referencia al secreto FCM |
| `NEXT_PUBLIC_API_BASE_URL` | dashboard | sí | URL pública de la API |
| `APU_INTERNAL_API_URL` | dashboard SSR | sí | URL interna de la API accesible desde el servidor Next.js |
| `APU_ANDROID_API_BASE_URL` | Android | sí | URL de la API para la aplicación |

Los nombres son el contrato inicial; cada aplicación debe incluir un
`.env.example` sin valores reales y documentar variables adicionales junto con
su tipo, valor permitido y responsable.

El comando `apu-admin` requiere temporalmente `APU_BOOTSTRAP_EMAIL`,
`APU_BOOTSTRAP_NAME` y `APU_BOOTSTRAP_PASSWORD`. No se guardan en `.env`; deben
eliminarse del entorno inmediatamente después de crear el primer administrador.

## Rotación

Rotar secretos inmediatamente ante exposición y de forma periódica según la
política operativa. La rotación debe soportar una ventana de transición cuando
la validación de tokens o cifrado lo requiera. Registrar fecha, responsable,
alcance y resultado, nunca el valor del secreto.
