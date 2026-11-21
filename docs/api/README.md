# Contrato de API

La API usa JSON y está versionada bajo `/api/v1`. El contrato formal se mantiene
en `openapi.yaml`; actualmente están implementados los endpoints de sesión y
salud. Los demás recursos de este documento describen el alcance planificado.

## Convenciones

- `Content-Type: application/json`.
- IDs: UUID en formato canónico.
- Fechas: ISO 8601 UTC.
- Coordenadas: latitud y longitud como números; validar rangos.
- Listados: `items`, `page`, `page_size`, `next_cursor`.
- Errores:

```json
{
  "error": {
    "code": "invalid_request",
    "message": "La solicitud no es válida",
    "request_id": "req_01..."
  }
}
```

`message` es seguro para mostrar; los detalles internos solo van a logs
redactados. Los códigos de error son estables y se documentan en OpenAPI.

## Autenticación

- Los usuarios humanos reciben una cookie de sesión opaca `HttpOnly`; el servidor
  guarda únicamente su hash.
- Las operaciones que modifican la sesión requieren `X-CSRF-Token`.
- Login, activación, recuperación y verificación MFA son públicas pero requieren
  credenciales o retos de un solo uso. No existe registro público.
- Los dispositivos Android usarán credenciales propias ligadas al dispositivo;
  ese mecanismo todavía no está implementado.
- El servidor responde `401` por credencial ausente/inválida y `403` por falta de
  autorización sobre el recurso.

## Idempotencia

Los envíos de datos y operaciones que crean comandos aceptan
`Idempotency-Key`. La misma clave con un payload diferente devuelve
`409 conflict`. Android además envía `event_id` por cada muestra.

## Recursos iniciales

| Método | Ruta | Autorización | Resultado |
|---|---|---|---|
| POST | `/auth/login` | pública | sesión |
| GET | `/auth/me` | usuario | identidad, rol y destino |
| POST | `/auth/logout` | usuario + CSRF | revoca la sesión |
| POST | `/auth/password/reset/request` | pública | respuesta uniforme y envío de enlace |
| POST | `/auth/password/reset/confirm` | token de un solo uso | cambia contraseña y revoca sesiones |
| POST | `/auth/password/change` | usuario + CSRF | reautentica, cambia contraseña y revoca sesiones |
| POST | `/auth/mfa/enroll` | reto después de contraseña | prepara TOTP |
| POST | `/auth/mfa/enroll/confirm` | reto + TOTP | activa MFA y entrega códigos de recuperación |
| POST | `/auth/mfa/verify` | reto + TOTP/recuperación | crea sesión con MFA verificado |
| GET | `/subscription` | tenant owner/admin | estado y límites |
| POST | `/billing/checkout` | tenant owner | inicia suscripción |
| POST | `/billing/webhook` | proveedor verificado | actualiza pago |
| GET | `/platform/summary` | platform admin | resumen operativo |
| GET/POST | `/platform/tenants` | platform admin | lista o crea tenants |
| PATCH | `/platform/tenants/{id}/status` | platform admin + CSRF | cambia estado del tenant |
| GET | `/platform/users` | platform admin | lista usuarios y membresías |
| POST | `/platform/invitations` | platform admin + CSRF | invita o asigna usuario |
| PATCH | `/platform/users/{id}/status` | platform admin + CSRF | suspende o reactiva cuenta |
| POST | `/platform/users/{id}/sessions/revoke` | platform admin + CSRF | revoca sesiones activas |
| PATCH | `/platform/tenants/{tenantId}/users/{userId}` | platform admin + CSRF | cambia rol de membresía |
| GET | `/platform/audit` | platform admin | auditoría operacional |
| POST | `/pairing/code` | tutor | código temporal |
| POST | `/pairing/confirm` | dispositivo | dispositivo vinculado |
| GET | `/devices` | tutor | dispositivos de su familia |
| GET | `/devices/{id}` | tutor | detalle autorizado |
| PUT | `/devices/{id}/modules` | tutor | estado de módulos |
| POST | `/devices/{id}/locations` | dispositivo | lote idempotente |
| GET | `/devices/{id}/locations` | tutor | historial paginado |
| GET | `/devices/{id}/events` | tutor | eventos derivados |
| POST | `/devices/{id}/commands` | tutor | comando con expiración |
| WS | `/ws` | tutor | eventos autorizados |

## Reglas para el contrato formal

Antes de escribir un endpoint:

1. definir request, response y errores en OpenAPI;
2. definir actor y comprobación de autorización;
3. definir idempotencia, límites y reintentos;
4. definir auditoría y retención;
5. agregar casos positivos y negativos a tests;
6. generar o revisar clientes tipados sin copiar modelos de base de datos.
