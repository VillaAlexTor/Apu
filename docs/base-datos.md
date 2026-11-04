# Base de datos

## Tecnología y responsabilidades

PostgreSQL es la fuente de verdad para tenants, identidad, suscripciones,
familias, dispositivos, configuración, ubicaciones, auditoría y estados de
comandos. La aplicación es
responsable de migraciones; el ORM no debe crear cambios implícitos en
producción.

## Convenciones

- UUID para entidades públicas.
- `BIGSERIAL` o identidad para tablas de eventos de alto volumen.
- Todas las tablas tienen timestamps UTC cuando aplique.
- Claves foráneas explícitas e índices para filtros por `family_id` y `device_id`.
- `NOT NULL` por defecto; nulabilidad solo con significado de dominio.
- CHECK constraints para estados y rangos simples.
- JSONB solo para payloads versionados, no para relaciones esenciales.

## Integridad y aislamiento

- Toda lectura sensible resuelve primero la familia autorizada.
- Las mutaciones relacionadas usan transacción.
- Los lotes de ubicación tienen restricción única por
  `(device_id, event_id)`.
- El borrado de un dispositivo debe definir explícitamente cascada, revocación o
  retención; nunca depender del comportamiento por defecto.

## Índices iniciales

- `guardians(email)` único normalizado.
- `tenant_users(tenant_id, user_id)`.
- `devices(tenant_id, family_id, status)`.
- `locations(device_id, recorded_at DESC)`.
- `audit_log(tenant_id, family_id, at DESC)`.
- `commands(device_id, status, expires_at)`.

Los índices se validan con `EXPLAIN` sobre volúmenes representativos.

## Cifrado y secretos

Las coordenadas se cifran a nivel de campo con claves fuera de PostgreSQL. Las
claves tienen versión para permitir rotación. El backend nunca registra
plaintext ni claves. La conexión a PostgreSQL usa TLS fuera del entorno local.

## Retención

Un job periódico elimina datos vencidos por tipo, en lotes acotados y con
auditoría operacional. Debe ser idempotente, reanudable y observable.

## Backups

- Backups cifrados y con acceso mínimo.
- Retención de backups documentada separadamente de retención de producto.
- Pruebas periódicas de restauración.
- No restaurar producción sobre entornos de desarrollo sin anonimización.
