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
| `APU_JWT_ISSUER` | API | sí | emisor de tokens |
| `APU_ACCESS_TOKEN_TTL` | API | sí | duración máxima del access token |
| `APU_REFRESH_TOKEN_PEPPER` | API | sí | secreto para proteger refresh tokens |
| `APU_FIELD_ENCRYPTION_KEY` | API | producción | clave de cifrado de campos |
| `APU_FCM_CREDENTIALS` | API | producción | referencia al secreto FCM |
| `NEXT_PUBLIC_API_BASE_URL` | dashboard | sí | URL pública de la API |
| `APU_ANDROID_API_BASE_URL` | Android | sí | URL de la API para la aplicación |

Los nombres son el contrato inicial; cada aplicación debe incluir un
`.env.example` sin valores reales y documentar variables adicionales junto con
su tipo, valor permitido y responsable.

## Rotación

Rotar secretos inmediatamente ante exposición y de forma periódica según la
política operativa. La rotación debe soportar una ventana de transición cuando
la validación de tokens o cifrado lo requiera. Registrar fecha, responsable,
alcance y resultado, nunca el valor del secreto.
