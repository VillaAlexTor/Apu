# Eventos y mensajería

## Tipos

APU distingue:

- **Evento de dominio:** hecho persistente, por ejemplo `location.recorded`.
- **Evento de integración:** mensaje entre componentes o proveedor externo.
- **Evento de UI:** actualización no autoritativa para el panel.

## Formato

```json
{
  "id": "evt_01...",
  "type": "location.recorded",
  "version": 1,
  "occurred_at": "2026-11-04T20:00:00Z",
  "family_id": "uuid",
  "device_id": "uuid",
  "data": {}
}
```

El payload de ejemplo es intencionalmente vacío: no se deben incluir coordenadas
ni secretos en logs o canales no autorizados.

## Eventos iniciales

- `device.paired`
- `device.revoked`
- `module.changed`
- `location.recorded`
- `geofence.entered`
- `geofence.exited`
- `command.created`
- `command.acknowledged`

## Reglas

- Todo evento tiene ID único y versión.
- Consumidores deben tolerar duplicados.
- Los eventos no conceden autorización.
- El WebSocket entrega solo eventos de la familia autenticada.
- Los cambios de contrato requieren incrementar versión o compatibilidad.
- Los eventos sensibles se minimizan y se retienen según su política.

## WebSocket

La conexión se autentica al abrirse y se revalida al renovar sesión. El cliente
puede reconectar usando un cursor; si pierde eventos, recupera el estado por
REST. El servidor limita conexiones, tamaño de mensaje y frecuencia.

## FCM

FCM es un transporte de comandos, no la fuente de verdad. El backend guarda el
comando, expiración y estado; el dispositivo confirma recepción y ejecución.

