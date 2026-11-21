# Autenticación y sesiones

## Flujo

1. El usuario inicia sesión con email y contraseña.
2. La API valida credenciales y estado de la cuenta.
3. Si el usuario es `platform_admin`, `tenant_owner` o `tenant_admin`, la API
   exige TOTP o un código de recuperación antes de crear la sesión.
4. La API genera un identificador de sesión opaco de alta entropía.
5. PostgreSQL almacena únicamente el hash SHA-256 del identificador.
6. El navegador recibe la sesión en una cookie `HttpOnly`; cada petición vuelve
   a comprobar expiración, revocación y estado del usuario.

Las contraseñas se almacenan con Argon2id. Los identificadores de sesión y
tokens CSRF nunca se almacenan en texto plano ni aparecen en logs. Los tokens
de dispositivos Android serán credenciales distintas y no reutilizarán la
sesión humana del panel.

## Invitaciones

No existe registro público. Un `platform_admin` crea la organización y asigna
el rol inicial. Para cuentas nuevas, la API genera un token de activación de un
solo uso, almacena únicamente su hash y lo invalida al aceptarlo o al superar
72 horas. El enlace se muestra una sola vez al operador; el transporte por
correo se integrará cuando exista un proveedor configurado.

## Revocación

Revocar todas las sesiones ante cambio de contraseña, desactivación, sospecha de
compromiso o cambio sensible de rol. Cerrar sesión revoca inmediatamente el
registro correspondiente. Revocar también la credencial del dispositivo al
desvincularlo.

## MFA y recuperación

MFA TOTP es obligatorio para `platform_admin`, `tenant_owner` y `tenant_admin`.
El secreto TOTP se cifra con AES-256-GCM, cada paso TOTP solo puede utilizarse
una vez y el alta entrega diez códigos de recuperación de un solo uso. Los
retos de autenticación vencen después de cinco minutos y se bloquean al octavo
intento fallido.

La recuperación usa un token aleatorio de un solo uso, almacena únicamente su
hash, vence después de 30 minutos y responde de manera indistinguible exista o
no la cuenta. Al restablecer o cambiar la contraseña se revocan todas las
sesiones. Producción exige SMTP con STARTTLS; el token local solo puede
exponerse cuando `APU_EXPOSE_DEVELOPMENT_TOKENS=true` fuera de producción.

Las contraseñas nuevas requieren al menos 15 caracteres, aceptan frases y
Unicode, y se rechazan si pertenecen a la lista local de valores comunes. No se
fuerzan cambios periódicos sin evidencia de compromiso.

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
