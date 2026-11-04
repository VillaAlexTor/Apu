# Arquitectura técnica de APU

**Estado:** propuesta para revisión  
**Versión:** 0.1  
**Última actualización:** 2026-11-04

## 1. Objetivo

Definir la arquitectura técnica, los límites entre componentes y las reglas de
implementación de APU antes de escribir código. Esta propuesta prioriza:

- transparencia para el menor y el tutor;
- aislamiento estricto entre familias;
- minimización, cifrado y retención limitada de datos;
- funcionamiento tolerante a conectividad intermitente;
- posibilidad de desactivar módulos sin reinstalar la aplicación.

El MVP se limita a autenticación, emparejamiento, ubicación, historial y zonas
seguras. Uso de aplicaciones, límites, filtro web y notificaciones quedan fuera
del primer incremento ejecutable.

## 2. Contexto y límites

### Actores

- **Tutor:** administra una familia y consulta los dispositivos vinculados.
- **Menor:** utiliza el dispositivo Android y debe conocer los módulos activos.
- **Dispositivo:** aplicación Android que recolecta únicamente los datos habilitados.
- **Administrador técnico:** opera la infraestructura, pero no debe poder leer
  datos cifrados de extremo a extremo.

### Componentes

```text
Android (Kotlin/Compose)
  ├─ presentación y consentimiento
  ├─ módulos de recolección
  ├─ cola local Room
  ├─ sincronización WorkManager
  └─ identidad del dispositivo (Android Keystore)
             │ HTTPS + firma de solicitudes
             ▼
API (Go)
  ├─ HTTP/JSON versionado
  ├─ autenticación y autorización
  ├─ emparejamiento
  ├─ módulos de ubicación y reglas
  ├─ WebSocket para eventos del panel
  └─ FCM para comandos
             │
             ▼
PostgreSQL ── migraciones versionadas
Redis (opcional, fase posterior: rate limit y tareas)

Panel (Next.js/React)
  ├─ rutas protegidas
  ├─ estado remoto/local
  ├─ mapa y reportes
  └─ confirmación explícita de comandos
```

## 3. Principios de diseño

1. **Modularidad:** cada capacidad se activa por dispositivo y tiene permisos,
   recolector, endpoint, política de retención y auditoría propios.
2. **Autorización por familia:** ningún handler acepta un `device_id` sin
   comprobar que pertenece a la familia del usuario autenticado.
3. **Servidor como fuente de verdad:** el cliente puede reintentar y cachear,
   pero no decide permisos, pertenencia ni reglas de seguridad.
4. **Idempotencia:** envíos de lotes y comandos incluyen un identificador único
   para evitar duplicados al reintentar.
5. **Fallos explícitos:** errores de validación, autenticación, permisos,
   sincronización y persistencia se registran y se exponen con códigos estables.
6. **Observabilidad sin datos sensibles:** logs estructurados con IDs técnicos,
   nunca coordenadas, contenido de mensajes, tokens o secretos.

## 4. Capas del backend

Cada módulo seguirá este flujo:

```text
HTTP handler → DTO/validación → caso de uso → repositorio/servicio → respuesta
```

- **Handler:** traduce HTTP, autentica el contexto y devuelve códigos de estado.
- **DTO:** valida entrada y salida; no expone entidades de persistencia.
- **Caso de uso:** contiene reglas de negocio y autorización de dominio.
- **Repositorio:** encapsula SQL/GORM y transacciones.
- **Servicio externo:** FCM, reloj, generación de códigos o cifrado.

Estructura prevista:

```text
backend/
├─ cmd/apu-server/
├─ internal/
│  ├─ http/              # router, middleware, errores y transporte
│  ├─ auth/
│  ├─ family/
│  ├─ device/
│  ├─ pairing/
│  ├─ location/
│  ├─ geofence/
│  ├─ rules/
│  ├─ push/
│  ├─ audit/
│  ├─ storage/
│  └─ platform/          # clock, id, crypto y configuración
├─ migrations/
└─ go.mod
```

La lógica de negocio no importará paquetes HTTP, GORM ni Firebase. Esto permite
probar casos de uso sin red ni base de datos.

## 5. Android

La app usará una arquitectura por capas:

```text
Compose UI → ViewModel → caso de uso → repositorio → API/Room/Android APIs
```

- **UI:** muestra estado de consentimiento, módulos, permisos y sincronización.
- **ViewModel:** expone estado observable y eventos de usuario.
- **Casos de uso:** coordinan permisos, recolección y sincronización.
- **Repositorios:** combinan fuente local y remota.
- **Workers:** ejecutan sincronización acotada, con backoff exponencial.
- **Foreground service:** solo para ubicación activa y con notificación visible.

Reglas:

- no iniciar un recolector si el módulo está desactivado;
- pedir permisos justo antes de activar la capacidad que los necesita;
- conservar eventos pendientes localmente cifrados;
- borrar datos locales al desvincular el dispositivo;
- tratar respuestas repetidas o fuera de orden de forma segura.

## 6. Panel web

El panel se organizará por dominio:

```text
dashboard/src/
├─ app/                  # rutas y layouts
├─ features/
│  ├─ auth/
│  ├─ devices/
│  ├─ locations/
│  └─ settings/
├─ components/           # componentes visuales compartidos
├─ lib/                  # cliente HTTP, sesión y configuración
└─ stores/               # Zustand únicamente para estado de UI/sesión
```

El acceso a API se centralizará en un cliente tipado. Las vistas no construirán
URLs ni headers manualmente. El WebSocket solo actualizará cachés ya autorizadas;
no reemplaza la autorización de cada petición HTTP.

## 7. Flujos críticos

### Emparejamiento

1. El tutor solicita un código de un solo uso.
2. El servidor guarda solo un hash del código, con TTL de cinco minutos.
3. Android introduce o escanea el código y muestra consentimiento.
4. Android genera su clave en Keystore y envía la clave pública.
5. El servidor valida el código, crea el dispositivo y lo asocia a la familia.
6. Se registra la aceptación de módulos y se emite una credencial revocable.
7. El servidor invalida el código y audita el evento.

El código no se reutiliza, no se registra en logs y se limita por IP/dispositivo.

### Ubicación

1. El recolector produce una muestra solo si el módulo está activo y el permiso
   está concedido.
2. La app la encola con `event_id`, hora del dispositivo y precisión.
3. WorkManager agrupa eventos y los envía con idempotency key.
4. El backend valida dispositivo, firma, rango y antigüedad.
5. La transacción persiste el lote y publica un evento al panel.
6. El servidor aplica la política de retención y evalúa geocercas.

### Comandos

Los comandos tienen `command_id`, tipo, expiración y estado
`pending|delivered|acknowledged|expired|failed`. Nunca se ejecuta un comando
sin verificar que el tutor tiene autorización sobre el dispositivo.

## 8. Consistencia, reintentos y disponibilidad

- POST de lotes: idempotencia obligatoria por `event_id`.
- Reintentos: máximo definido por operación; backoff exponencial con jitter.
- Un `401` provoca renovación controlada; un `403` no se reintenta.
- Un `429` respeta `Retry-After`.
- El cliente conserva una cola limitada y muestra cuando hay datos pendientes.
- El backend usa transacciones para cambios de identidad, emparejamiento y reglas.

## 9. Despliegue inicial

Desarrollo:

```text
Docker Compose: API + PostgreSQL
Local: dashboard + Android emulator/device
```

Producción futura:

- API stateless detrás de TLS y proxy;
- PostgreSQL con backups cifrados, migraciones controladas y recuperación probada;
- FCM con credenciales en un gestor de secretos;
- métricas, logs y alertas sin contenido sensible;
- separación de entornos y cuentas de servicio con mínimo privilegio.

## 10. Criterios de aceptación arquitectónicos

Antes de comenzar funcionalidades posteriores al MVP:

- existe autorización automatizada por familia en cada endpoint;
- toda migración puede aplicarse desde una base vacía;
- los reintentos no duplican ubicación ni comandos;
- el consentimiento y los módulos activos son visibles en Android;
- los datos vencidos se eliminan mediante un proceso verificable;
- los tests cubren errores, permisos y conectividad intermitente.

