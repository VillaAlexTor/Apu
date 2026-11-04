# Convenciones de código

## Generales

- Preferir nombres explícitos sobre abreviaturas.
- Mantener funciones pequeñas y con una responsabilidad.
- Validar entradas en los bordes y reglas en el dominio.
- No ocultar errores con valores por defecto.
- Evitar comentarios obvios; documentar solo decisiones no evidentes.

## Go

- `gofmt`, `go vet`, `staticcheck` y `gosec` en CI.
- Paquetes internos por dominio, no por tipo técnico global.
- `context.Context` como primer parámetro de operaciones con I/O.
- Errores de dominio tipados y traducción HTTP en un único lugar.
- DTOs separados de entidades de persistencia.
- Repositorios detrás de interfaces usadas por los casos de uso.
- Logs estructurados; nunca secretos, tokens ni coordenadas.

## Kotlin/Android

- Kotlin idiomático y coroutines estructuradas.
- Compose solo presenta estado; la lógica vive en ViewModel/casos de uso.
- Room para cola y estado local; WorkManager para trabajo diferible.
- Permisos solicitados desde el flujo que activa la capacidad.
- Foreground service únicamente para una función visible y justificada.
- Tests para ViewModel, repositorios y workers.

## TypeScript/Next.js

- TypeScript estricto.
- Cliente HTTP centralizado y respuestas tipadas.
- Componentes presentacionales separados de acceso a datos.
- Zustand solo para estado global necesario; no duplicar caché del servidor.
- Validar datos externos antes de usarlos en componentes.
- Accesibilidad y estados de carga/error como parte del componente.

## APIs y persistencia

- Nombres de recursos en plural y rutas versionadas.
- Fechas en UTC y campos con sufijos claros.
- DTOs explícitos; no serializar modelos de base directamente.
- Toda mutación sensible define idempotencia, autorización y auditoría.
- Migraciones pequeñas, numeradas y no editables después de aplicarse.

## Tests

- Nombres que describan comportamiento.
- Un caso de fallo por cada regla de autorización.
- Fixtures sintéticos y deterministas.
- Cada defecto corregido agrega una regresión.

