# APU Dashboard

Frontend web de APU. La primera pantalla es un acceso compartido para tutores, clientes y personal operador; el backend determinará el dashboard y los permisos después de autenticar la sesión.

## Desarrollo

Requiere Node.js 24 y pnpm 10.

```bash
pnpm install
pnpm dev
```

La aplicación queda disponible en [http://localhost:3001](http://localhost:3001).

## Verificación

```bash
pnpm lint
pnpm build
```

El build usa Webpack porque la versión actual de Turbopack presenta un fallo interno al procesar CSS en Windows.

## Integración con el backend

El formulario envía las credenciales a `POST /api/v1/auth/login`, conserva la
sesión mediante cookies seguras y muestra el destino que corresponde al rol.
En desarrollo, `NEXT_PUBLIC_API_URL` puede apuntar a
`http://localhost:8081/api/v1`.

Las rutas privadas consultan `/auth/me` durante el render del servidor y
redirigen al login o al panel autorizado. La sesión persiste en la cookie opaca
del backend: los roles de plataforma reciben `/platform`, `tenant_owner` y
`tenant_admin` reciben `/admin`, `guardian` recibe `/familia`, y una cuenta sin
asignación recibe `/acceso-pendiente`. En este corte, la operación completa de
`/platform` está autorizada únicamente para `platform_admin`.

`/platform` implementa organizaciones, invitaciones, usuarios, cambios de rol y
estado, revocación de sesiones y auditoría. El cierre de sesión llama a
`POST /auth/logout`, revoca la sesión actual y regresa al acceso compartido.
`/activar` permite que una persona invitada establezca su contraseña.
