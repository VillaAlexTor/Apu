# Operación de datos

## Solicitudes

Antes del despliegue público se deben definir canal, identidad, responsable y
plazos para solicitudes de acceso, corrección, exportación y eliminación. Cada
solicitud recibe un identificador y se registra sin copiar datos sensibles en el
ticket.

## Exportación

La exportación debe:

- verificar al solicitante y su alcance de tenant/familia;
- producir un formato documentado y consistente;
- limitar tamaño, frecuencia y tiempo de disponibilidad;
- cifrar el archivo durante transporte y almacenamiento temporal;
- registrar quién la generó y cuándo expiró.

## Eliminación

Desvincular un dispositivo revoca credenciales, detiene comandos y borra la cola
local. El backend purga datos según tipo y retención; la eliminación lógica solo
es válida si existe una purga posterior verificable. Los backups se tratan
según su ciclo técnico documentado.

## Corrección y retención

Las correcciones conservan la auditoría necesaria, sin reintroducir datos
eliminados. Los jobs de retención son idempotentes, reanudables, observables y
registran cantidades y tipos, nunca contenido.

## Acceso operativo

Soporte y operadores trabajan con datos mínimos y enmascarados. Cualquier acceso
excepcional a datos sensibles requiere autorización específica, motivo,
duración limitada y auditoría. No se usan datos de producción en desarrollo,
fixtures, capturas ni demos.
