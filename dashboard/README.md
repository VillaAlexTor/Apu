# APU Dashboard

Frontend web de APU. La primera pantalla es un acceso compartido para tutores, clientes y personal operador; el backend determinará el dashboard y los permisos después de autenticar la sesión.

## Desarrollo

Requiere Node.js 24 y pnpm 10.

```bash
pnpm install
pnpm dev
```

La aplicación queda disponible en [http://localhost:3000](http://localhost:3000).

## Verificación

```bash
pnpm lint
pnpm build
```

El build usa Webpack porque la versión actual de Turbopack presenta un fallo interno al procesar CSS en Windows.

## Backend pendiente

El formulario valida los datos y representa sus estados de carga y respuesta, pero todavía no crea una sesión. La integración prevista es con la API de Go del proyecto.
