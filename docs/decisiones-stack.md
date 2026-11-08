# Decisiones del stack tecnológico

**Estado:** aceptado para el MVP  
**Fecha:** 2026-11-08  
**Alcance:** panel web, API, worker, aplicación Android, datos, seguridad,
observabilidad, pruebas y despliegue de APU.

## 1. Principios de decisión

El stack debe favorecer seguridad explícita, aislamiento multi-tenant,
operación sencilla y contratos verificables. El MVP evita infraestructura que
no sea necesaria para el volumen inicial y mantiene separados los límites de
compilación del panel, backend y Android.

Principios obligatorios:

- la autorización se ejecuta en el backend y nunca se delega al cliente;
- toda operación sensible sigue `actor → tenant → familia → dispositivo`;
- REST es la fuente de verdad y WebSocket solo acelera actualizaciones;
- los datos de desarrollo son sintéticos;
- las dependencias deben mantenerse en versiones estables soportadas;
- no se incorporan servicios adicionales sin una necesidad operativa medida.

## 2. Monorepo

Se mantiene un único repositorio con aplicaciones independientes:

```text
apu/
├── dashboard/             # Next.js
├── backend/               # API y worker en Go
├── android/               # Aplicación Android en Kotlin
├── docs/                  # Producto, arquitectura y operación
├── docker-compose.yml
└── .github/workflows/
```

No se comparten dependencias entre ecosistemas. CI detecta las carpetas
afectadas y ejecuta únicamente sus validaciones, además de una verificación
integral antes de publicar.

## 3. Panel web

| Área | Decisión |
|---|---|
| Runtime | Node.js 26 LTS |
| Framework | Next.js estable con App Router |
| Lenguaje | TypeScript con modo estricto |
| UI | React y componentes propios sobre primitivas de shadcn/ui |
| Estilos | Tailwind CSS |
| Formularios | React Hook Form + Zod |
| Datos remotos | TanStack Query |
| Estado local | Zustand, únicamente para estado efímero de interfaz |
| Mapas MVP | Leaflet + React Leaflet |
| Pruebas | Vitest, Testing Library y Playwright |
| Gestor | pnpm |

### Reglas del frontend

- El login es compartido por todos los usuarios humanos y no incluye selector
  de rol.
- Después de autenticar, el panel consulta la identidad y capacidades
  autorizadas; la API decide el destino.
- TanStack Query administra estado del servidor. Zustand no duplica cachés,
  usuarios, permisos ni respuestas de la API.
- Zod valida la experiencia del formulario, pero el backend repite toda
  validación y autorización.
- La dirección visual inicial es **Atlas topográfico andino**. La composición
  aprobada en `.impeccable/mocks/decision/login-atlas.webp` es la referencia
  de implementación del login.
- La cartografía es un sistema visual y estructural; no debe convertirse en
  folclore literal ni reducirse a una textura decorativa.

## 4. Backend y API

| Área | Decisión |
|---|---|
| Lenguaje | Go 1.27.x estable |
| HTTP/router | `net/http` + `chi` |
| Contrato | REST JSON bajo `/api/v1` + OpenAPI 3.1 |
| PostgreSQL | `pgx/v5` |
| Consultas tipadas | `sqlc` |
| Migraciones | Goose, numeradas y reversibles cuando sea seguro |
| Logs | `log/slog` con salida estructurada |
| Métricas/trazas | OpenTelemetry + Prometheus |
| Tiempo real | WebSocket complementario |

### Estructura prevista

```text
backend/
├── cmd/
│   ├── apu-api/
│   └── apu-worker/
├── internal/
│   ├── auth/
│   ├── tenant/
│   ├── family/
│   ├── device/
│   ├── location/
│   ├── subscription/
│   ├── audit/
│   └── platform/
├── migrations/
└── openapi/
```

`apu-api` atiende HTTP y WebSocket. `apu-worker` ejecuta retención,
notificaciones, reconciliación de pagos y tareas programadas. Comparten
dominio y repositorios, pero pueden desplegarse por separado.

### Decisiones explícitas

- Se usa `chi` en lugar de Gin para conservar integración directa con
  `net/http`, middleware pequeño y pruebas sencillas.
- Se usa `pgx + sqlc` en lugar de GORM. Las consultas SQL deben mostrar los
  filtros de tenant y producir código tipado.
- Los handlers no acceden directamente a la base de datos: llaman casos de
  uso que reciben un actor autenticado.
- WebSocket no concede permisos ni es fuente exclusiva de eventos; todo estado
  se puede recuperar mediante REST.

## 5. Identidad, sesiones y autorización

### Usuarios humanos del panel

- Contraseñas protegidas con Argon2id y sal única por credencial.
- Sesiones opacas de alta entropía; la base almacena solo el hash del token.
- Cookie con prefijo `__Host-`, `Secure`, `HttpOnly`, `SameSite` y `Path=/`.
- Protección CSRF independiente; `SameSite` es defensa adicional, no sustituto.
- Revocación de sesión inmediata y registro de auditoría.
- Reautenticación para cambios de rol y acciones administrativas sensibles.
- MFA TOTP obligatorio para personal operador y roles administrativos definidos.
- Tokens de sesión, JWT y refresh tokens nunca se almacenan en `localStorage`
  ni `sessionStorage`.

El backend responde con la identidad y las capacidades efectivas. El frontend
redirige a una ruta conveniente, pero cada endpoint vuelve a comprobar permisos.

```text
platform_admin / support_agent → /platform
tenant_owner / tenant_admin    → /admin
guardian                       → /familia
```

### Dispositivos Android

- La clave privada se genera y permanece en Android Keystore.
- El backend registra únicamente la clave pública.
- Las credenciales de dispositivo son distintas de las sesiones humanas.
- Los tokens de acceso son breves, revocables y ligados a la identidad del
  dispositivo.
- Perder la clave o reinstalar la aplicación exige nuevo emparejamiento.

## 6. Base de datos

| Área | Decisión |
|---|---|
| Motor | PostgreSQL 18 estable |
| Aislamiento | `tenant_id` obligatorio en recursos de negocio |
| Identificadores | UUID para entidades de dominio |
| Auditoría | Tabla append-only con actor, acción, alcance y fecha |
| Datos sensibles | Cifrado a nivel de campo mediante claves externas a la BD |
| Historial | Retención automática configurable dentro de límites seguros |

No se adopta PostGIS antes de resolver el diseño de cifrado de coordenadas. El
MVP evalúa geocercas en el backend durante la ingestión y cifra las coordenadas
antes de persistirlas. Incorporar índices geoespaciales sobre datos derivados
requiere una ADR y revisión del modelo de amenazas.

## 7. Trabajos, caché y escalado

Redis no forma parte del MVP. Inicialmente:

- las sesiones se almacenan en PostgreSQL;
- los trabajos usan una tabla con bloqueos transaccionales y
  `FOR UPDATE SKIP LOCKED`;
- el rate limiting de una sola instancia permanece en memoria;
- WebSocket opera desde una sola instancia del API.

Redis se evalúa cuando existan varias instancias y sea necesario compartir
rate limits, eventos, presencia o caché. No se introduce anticipadamente.

## 8. Android

| Área | Decisión |
|---|---|
| Lenguaje | Kotlin |
| UI | Jetpack Compose |
| Persistencia local | Room con cola limitada y cifrada |
| Trabajo diferido | WorkManager |
| HTTP | Retrofit + OkHttp |
| Inyección | Hilt |
| Claves | Android Keystore |
| Push | Firebase Cloud Messaging |
| Pruebas | JUnit, pruebas de Compose e integración |

FCM notifica que existen comandos o cambios pendientes. El dispositivo obtiene
la operación autorizada desde la API; una notificación push por sí sola no es
autoridad para ejecutar una acción sensible.

## 9. Mapas y ubicación

- Leaflet y React Leaflet cubren el mapa del MVP.
- OpenStreetMap aporta datos cartográficos, pero producción debe usar un
  proveedor de tiles con capacidad y términos apropiados; no se presupone que
  el servidor público gratuito soporte el tráfico del producto.
- El historial se pagina y limita por rango temporal.
- Coordenadas, direcciones y tokens no aparecen en logs.
- Toda lectura de ubicación sensible genera auditoría.

MapLibre se reconsidera cuando el producto necesite mapas vectoriales,
renderizado avanzado, uso offline más amplio o estilos cartográficos propios.

## 10. Infraestructura y despliegue

| Área | Decisión |
|---|---|
| Desarrollo local | Docker Compose |
| Empaquetado | Imágenes Docker separadas por aplicación |
| CI/CD | GitHub Actions |
| Producción inicial | Contenedores administrados + PostgreSQL administrado |
| TLS | Terminado por plataforma o proxy administrado |
| Secretos | Gestor de secretos del entorno, nunca archivos versionados |
| Archivos futuros | Almacenamiento compatible con S3 |

El primer despliegue favorece una sola instancia de API y una de worker. La
arquitectura permite escalar horizontalmente después de añadir coordinación
distribuida y pruebas de carga.

## 11. Observabilidad

- Cada petición recibe `request_id` y propaga contexto.
- Logs JSON estructurados con redacción de datos sensibles.
- Métricas de latencia, errores, autenticación, sincronización y colas.
- Trazas OpenTelemetry entre handlers, casos de uso, repositorios y proveedores.
- Alertas por errores sostenidos, crecimiento de colas, fallos de retención,
  aislamiento denegado y anomalías administrativas.
- Nunca se registran contraseñas, tokens, coordenadas, mensajes o payloads
  sensibles completos.

## 12. Pruebas y controles de calidad

### Backend

- `go test ./...`, `go vet`, `staticcheck`, `gosec` y `govulncheck`.
- Tests de repositorios con PostgreSQL real mediante Testcontainers.
- Para cada operación: actor permitido, otra familia, otro tenant, sesión
  revocada, recurso inexistente y payload alterado.
- Pruebas de migración desde base vacía y desde la versión anterior.

### Frontend

- TypeScript estricto, lint y build de producción.
- Vitest y Testing Library para componentes y estados.
- Playwright para login, recuperación, MFA y redirección autorizada.
- Auditoría de accesibilidad en desktop y móvil.

### Android

- Tests unitarios de dominio y sincronización.
- Tests de Room, WorkManager, reintentos e idempotencia.
- Tests de permisos rechazados, conectividad intermitente, revocación y pérdida
  de clave.

## 13. Dependencias deliberadamente no adoptadas en el MVP

| Alternativa | Decisión |
|---|---|
| GORM | No; se prefiere SQL explícito con `pgx + sqlc`. |
| Gin | No; `chi` y `net/http` cubren el alcance con menor superficie. |
| Redis | No inicialmente; se incorpora solo con necesidad distribuida. |
| Microservicios | No; API y worker comparten un monolito modular. |
| GraphQL | No; REST versionado satisface web y Android. |
| Firebase Authentication | No; identidad, tenant y auditoría pertenecen a APU. |
| JWT persistido en navegador | No; se usan sesiones opacas seguras. |
| Kubernetes | No para el MVP; añade complejidad operativa prematura. |
| PostGIS | Pendiente de decisión de cifrado y modelo de amenazas. |

## 14. Orden de implementación

1. Crear estructura compilable del monorepo, Docker Compose y CI.
2. Crear PostgreSQL, migraciones, consultas `sqlc` y pruebas de aislamiento.
3. Implementar identidad, sesiones, recuperación, MFA y auditoría base.
4. Construir el login cartográfico en Next.js y conectarlo a Go.
5. Crear shells separados para `/platform`, `/admin` y `/familia`.
6. Implementar suscripciones y límites de plan.
7. Implementar emparejamiento Android y credenciales de dispositivo.
8. Implementar ubicación, historial, WebSocket, retención y geocercas.
9. Añadir observabilidad, hardening, pruebas de carga y revisión de seguridad.

## 15. Decisiones todavía abiertas

- proveedor de pagos y moneda operativa;
- proveedor de tiles y política de caché cartográfica;
- proveedor de correo transaccional;
- plataforma concreta de despliegue y región de datos;
- gestor de secretos de producción;
- parámetros finales de Argon2id después de medir el hardware objetivo;
- estrategia legal y criptográfica definitiva para coordenadas históricas;
- política de recuperación ante pérdida de MFA.

Estas decisiones requieren ADR propia antes de introducir una dependencia o
comprometer datos reales.

## 16. Referencias técnicas

- Next.js: <https://nextjs.org/docs/app/getting-started/installation>
- Go: <https://go.dev/doc/devel/release>
- PostgreSQL: <https://www.postgresql.org/support/versioning/>
- OWASP Password Storage:
  <https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html>
- OWASP Session Management:
  <https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet>
- Firebase Cloud Messaging: <https://firebase.google.com/docs/cloud-messaging>

