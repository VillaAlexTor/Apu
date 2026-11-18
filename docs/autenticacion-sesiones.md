# Autenticación y sesiones

## Flujo

1. El usuario inicia sesión con email y contraseña.
2. La API valida credenciales y estado de la cuenta.
3. La API genera un identificador de sesión opaco de alta entropía.
4. PostgreSQL almacena únicamente el hash SHA-256 del identificador.
5. El navegador recibe la sesión en una cookie `HttpOnly`; cada petición vuelve
   a comprobar expiración, revocación y estado del usuario.

Las contraseñas se almacenan con Argon2id. Los identificadores de sesión y
tokens CSRF nunca se almacenan en texto plano ni aparecen en logs. Los tokens
de dispositivos Android serán credenciales distintas y no reutilizarán la
sesión humana del panel.

## Revocación

Revocar todas las sesiones ante cambio de contraseña, desactivación, sospecha de
compromiso o cambio sensible de rol. Cerrar sesión revoca inmediatamente el
registro correspondiente. Revocar también la credencial del dispositivo al
desvincularlo.

## MFA y recuperación

MFA TOTP será obligatorio para roles administrativos y recomendado para
`tenant_owner`. Todavía no está implementado, por lo que esta versión no debe
habilitar cuentas administrativas en producción. La recuperación requerirá un
flujo con token de un solo uso, expiración corta, respuesta indistinguible para
emails inexistentes y auditoría.

## Controles

- rate limiting y backoff para login y recuperación;
- mensajes que no revelen si una cuenta existe;
- cookie de sesión `HttpOnly`, `Secure` y `SameSite`; en desarrollo local se
  permite HTTP y se elimina el prefijo `__Host-`;
- token CSRF independiente, verificado contra el hash de la sesión;
- reautenticación para roles, facturación y acciones destructivas.

El TTL y los códigos de error deben quedar sincronizados con
`api/openapi.yaml` y la configuración del entorno. La sesión inicial dura 12
horas y es configurable entre 15 minutos y 7 días.
