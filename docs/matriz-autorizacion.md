# Matriz de autorización

La autorización se evalúa en el backend. El cliente no puede conceder permisos
por enviar un `tenant_id`, `family_id` o `device_id`.

| Acción | Owner/Admin tenant | Guardian asignado | Dispositivo | Otro tenant | Operador |
|---|---:|---:|---:|---:|---:|
| Ver su familia | Sí | Sí | No | No | No |
| Generar código | Sí | Sí | No | No | No |
| Confirmar emparejamiento | No | No | Sí, con código | No | No |
| Listar sus dispositivos | Sí | Sí | No | No | No |
| Activar módulo | Sí | Según rol | No | No | No |
| Enviar ubicación | No | No | Sí, activo | No | No |
| Leer historial | Sí, auditado | Sí, auditado | No | No | No |
| Crear geocerca | Sí | Según rol | No | No | No |
| Enviar comando | Sí, confirmado | Según rol | No | No | No |
| Revocar dispositivo | Sí | No | No | No | No |
| Leer coordenadas descifradas | Sí, auditado | Sí, auditado | No | No | No |
| Administrar suscripción | Owner | No | No | No | No |
| Operar infraestructura | No | No | No | No | Sí, sin datos |

## Reglas obligatorias

- Toda ruta con dispositivo debe cargar `tenant_id` y `family_id`.
- Las lecturas de ubicación requieren auditoría.
- El dispositivo solo escribe datos propios y con credencial vigente.
- Un tutor no puede cambiar la pertenencia de un dispositivo desde una ruta
  ordinaria.
- El operador no tiene una ruta de negocio para leer datos
  sensibles.
- Los WebSocket se filtran por familia y se revalidan tras reconectar.

## Pruebas mínimas

Cada acción debe probarse con actor permitido, actor de otra familia, credencial
revocada, recurso inexistente y payload alterado.
