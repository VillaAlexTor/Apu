# Autenticación y sesiones

## Flujo

1. El usuario inicia sesión con email y contraseña.
2. La API valida credenciales y estado de la cuenta.
3. Emite un access token de corta duración y un refresh token rotativo.
4. Cada petición valida emisor, audiencia, expiración y sesión revocable.
5. El refresh usado se invalida y se emite uno nuevo.

Las contraseñas se almacenan con Argon2id. Los refresh tokens se almacenan
protegidos, nunca en texto plano ni en logs.

## Revocación

Revocar todas las sesiones ante cambio de contraseña, desactivación, sospecha de
compromiso o cambio sensible de rol. Revocar también la credencial del
dispositivo al desvincularlo. Una reutilización de refresh token invalida la
familia de tokens y genera una alerta auditable.

## MFA y recuperación

MFA TOTP es obligatorio para roles administrativos y recomendado para
`tenant_owner`. La recuperación requiere un flujo con token de un solo uso,
expiración corta, respuesta indistinguible para emails inexistentes y auditoría.

## Controles

- rate limiting y backoff para login, refresh y recuperación;
- mensajes que no revelen si una cuenta existe;
- cookies `HttpOnly`, `Secure` y `SameSite` cuando el dashboard use cookies;
- protección CSRF para operaciones basadas en cookies;
- reautenticación para roles, facturación y acciones destructivas.

Los detalles de claims, TTL y códigos de error deben quedar sincronizados con
`api/openapi.yaml` y la configuración del entorno.
