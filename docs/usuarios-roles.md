# Usuarios y roles

## Jerarquía

| Nivel | Actor | Alcance |
|---|---|---|
| Operación | `platform_admin` | Todos los tenants, sin datos sensibles por defecto |
| Operación | `support_agent` | Soporte técnico limitado y auditado |
| Cliente | `tenant_owner` | Un tenant; facturación, usuarios y configuración |
| Cliente | `tenant_admin` | Un tenant; usuarios, familias y dispositivos |
| Cliente | `guardian` | Familias y dispositivos asignados |
| Dispositivo | `device` | Su propia identidad y datos de sincronización |

El menor no es una cuenta de acceso al panel. Su conocimiento y consentimiento
se gestionan en el dispositivo y en la configuración de módulos.

## Reglas

- Un usuario puede pertenecer a uno o más tenants, con roles independientes.
- Un usuario solo ve familias explícitamente asignadas dentro de su tenant.
- `tenant_owner` puede transferir propiedad o cancelar la suscripción.
- `tenant_admin` no puede cambiar facturación ni elevar su propio rol.
- `support_agent` nunca obtiene acceso implícito a coordenadas descifradas.
- Los cambios de rol requieren reautenticación y auditoría.
- El dispositivo no es usuario humano y no puede usar rutas del panel.

## Ciclo de vida del usuario

`invited → active → suspended → deactivated`

Las invitaciones expiran, son de un solo uso y no revelan si un email ya existe.
Desactivar un usuario revoca sus sesiones, pero no elimina la auditoría histórica.

## Autenticación

Contraseña con Argon2id, MFA para roles administrativos, sesiones de corta
duración y refresh tokens rotativos. El proveedor de pagos no autentica usuarios
de APU: solo comunica eventos verificados.

