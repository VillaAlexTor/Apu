# Contrato de API

La API pública usa JSON y está versionada bajo `/api/v1`. El contrato formal se
publicará en `openapi.yaml` antes de implementar los handlers.

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

- `Authorization: Bearer <access_token>`.
- Access token de vida corta.
- Refresh token rotativo, revocable y almacenado de forma segura.
- Las rutas públicas se limitan a registro, login y comienzo de emparejamiento.
- El servidor responde `401` por credencial ausente/inválida y `403` por falta de
  autorización sobre el recurso.

## Idempotencia

Los envíos de datos y operaciones que crean comandos aceptan
`Idempotency-Key`. La misma clave con un payload diferente devuelve
`409 conflict`. Android además envía `event_id` por cada muestra.

## Recursos iniciales

| Método | Ruta | Autorización | Resultado |
|---|---|---|---|
| POST | `/auth/register` | pública | crea tutor |
| POST | `/auth/login` | pública | sesión |
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

