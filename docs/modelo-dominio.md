# Modelo de dominio

## Entidades principales

```text
Guardian ──< FamilyGuardian >── Family ──< Device
                                      ├──< DeviceModule
                                      ├──< Location
                                      ├──< Geofence
                                      ├──< Command
                                      └──< AuditEntry
```

### Guardian

Identidad autenticable del tutor. Contiene email normalizado, hash de contraseña,
estado de cuenta y timestamps. Nunca se devuelve el hash.

### Family

Agregado de autorización. Todas las operaciones sobre dispositivos comienzan
resolviendo la familia del actor.

### Device

Representa una instalación Android vinculada. Tiene alias, clave pública,
estado (`pending`, `active`, `revoked`) y última conexión.

### DeviceModule

Estado de una capacidad por dispositivo:
`disabled`, `pending_permission`, `active` o `error`.

### Location

Muestra cifrada con `device_id`, latitud, longitud, precisión, `recorded_at` y
`event_id` único por dispositivo.

### Geofence

Área circular con nombre, centro, radio, estado activo y configuración de
notificaciones.

### Command

Instrucción dirigida a un dispositivo. Estados:
`pending`, `delivered`, `acknowledged`, `expired` y `failed`.

### AuditEntry

Registro de actor, familia, dispositivo, acción, resultado, request ID y hora.
No contiene secretos ni payloads sensibles.

## Reglas de negocio

1. Un dispositivo pertenece a una sola familia activa.
2. Un código de emparejamiento es de un solo uso y expira.
3. Un dispositivo revocado no puede enviar datos ni recibir comandos.
4. Una muestra repetida por `event_id` es un éxito idempotente.
5. Un módulo desactivado no puede producir nuevas muestras.
6. La retención nunca puede superar el máximo definido por política.
7. El acceso a un recurso requiere pertenencia a la familia del actor.
8. Un comando vencido no se entrega ni se ejecuta.

## Estados y transiciones

```text
Device: pending → active → revoked
Module: disabled → pending_permission → active → disabled
Command: pending → delivered → acknowledged
                  ├→ expired
                  └→ failed
```

Las transiciones inválidas deben devolver un error de dominio explícito y no
modificar la base de datos.

