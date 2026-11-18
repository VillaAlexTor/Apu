# Panel de administración

## Propósito

El panel interno permite operar el SaaS sin convertirse en una vía para leer
datos familiares sensibles.

## Capacidades

- crear, activar, suspender y terminar tenants;
- invitar y desactivar usuarios;
- asignar planes y revisar estados de facturación;
- consultar salud de API, workers, FCM y sincronización;
- ver auditoría operacional;
- gestionar tickets y notas de soporte sin datos sensibles;
- ejecutar acciones excepcionales con doble confirmación.

## Estado implementado

El panel `/platform`, autorizado en este corte para `platform_admin`, ya permite
consultar el resumen operativo, crear y cambiar el estado de tenants, invitar
usuarios, asignar roles, suspender, reactivar o desactivar cuentas, revocar
sesiones y consultar los últimos registros de auditoría. Las acciones se
autorizan nuevamente en el backend y requieren protección CSRF. Suspender o
desactivar una cuenta revoca sus sesiones activas.

Facturación, salud de workers, FCM, tickets y acciones excepcionales permanecen
fuera de este corte.

## Restricciones

- No mostrar coordenadas ni historial de ubicación en listados administrativos.
- No permitir consultas SQL libres desde la aplicación.
- No permitir impersonación silenciosa.
- El acceso excepcional debe indicar motivo, duración, aprobador y alcance.
- Cada acción administrativa genera auditoría inmutable.
- El soporte debe usar identificadores técnicos y datos sintéticos en pruebas.

## Estados visibles

El panel muestra estado del tenant, plan, pago, usuarios, dispositivos,
última sincronización, errores y retención, pero no contenido sensible.

## Seguridad

MFA obligatorio, sesiones cortas, allowlist de red cuando sea viable, revisión
periódica de roles y alertas ante exportaciones o accesos masivos.

