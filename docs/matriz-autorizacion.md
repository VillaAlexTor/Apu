# Matriz de autorización

La autorización se evalúa en el backend. El cliente no puede conceder permisos
por enviar un `family_id` o `device_id`.

| Acción | Tutor de la familia | Dispositivo vinculado | Otra familia | Admin técnico |
|---|---:|---:|---:|---:|
| Ver su familia | Sí | No | No | No |
| Generar código | Sí | No | No | No |
| Confirmar emparejamiento | No | Sí, con código | No | No |
| Listar sus dispositivos | Sí | No | No | No |
| Activar módulo | Sí | No | No | No |
| Enviar ubicación | No | Sí, activo | No | No |
| Leer historial | Sí, auditado | No | No | No |
| Crear geocerca | Sí | No | No | No |
| Enviar comando | Sí, confirmado | No | No | No |
| Revocar dispositivo | Sí | No | No | No |
| Leer coordenadas descifradas | Sí, auditado | No | No | No |
| Operar infraestructura | No | No | No | Solo sin datos |

## Reglas obligatorias

- Toda ruta con dispositivo debe cargar el dispositivo junto con `family_id`.
- Las lecturas de ubicación requieren auditoría.
- El dispositivo solo escribe datos propios y con credencial vigente.
- Un tutor no puede cambiar la pertenencia de un dispositivo desde una ruta
  ordinaria.
- El administrador técnico no tiene una ruta de negocio para leer datos
  sensibles.
- Los WebSocket se filtran por familia y se revalidan tras reconectar.

## Pruebas mínimas

Cada acción debe probarse con actor permitido, actor de otra familia, credencial
revocada, recurso inexistente y payload alterado.

