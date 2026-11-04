# Despliegue y operación

## Entornos

- **Local:** datos sintéticos y servicios Docker.
- **CI:** servicios efímeros, sin secretos de producción.
- **Staging:** configuración similar a producción y datos sintéticos.
- **Producción:** acceso restringido, secretos gestionados y backups activos.

Cada entorno tiene base de datos, claves, credenciales y proveedores de pago
independientes. El repositorio y los artefactos de producción son privados.

## Pipeline

1. Formato y lint.
2. Tests unitarios.
3. Tests de integración y migraciones.
4. Escaneo de dependencias y secretos.
5. Build reproducible de API, dashboard y APK.
6. Despliegue a staging.
7. Smoke tests.
8. Aprobación explícita para producción.

## Configuración

La configuración se inyecta por variables de entorno o gestor de secretos.
Debe validarse al arrancar; una variable obligatoria ausente impide iniciar el
servicio. No usar secretos por defecto.

## Migraciones

Las migraciones se ejecutan como paso controlado antes de desplegar código que
las necesita. Para cambios destructivos usar expansión, migración de datos,
cambio de lectura/escritura y contracción en releases separados.

## Backups y recuperación

- Backup cifrado y automatizado.
- Restauración probada periódicamente.
- RPO/RTO se fijarán antes de producción.
- Documentar quién puede restaurar y cómo se rota el acceso.

## Rollback

Un rollback debe poder revertir la aplicación sin destruir datos. Si una
migración no es reversible, el plan debe incluir compatibilidad hacia atrás y
restauración verificada.
