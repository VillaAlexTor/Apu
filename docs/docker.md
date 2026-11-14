# Ejecución local con Docker

Esta guía permite iniciar el frontend, el backend o todo APU usando imágenes
separadas. Los valores incluidos son exclusivamente para desarrollo local y no
deben reutilizarse en producción.

## Estado actual

- **Dashboard:** login funcional en Next.js, todavía sin autenticación real.
- **Backend:** bootstrap Go con `/healthz`, `/readyz` y `/api/v1`; todavía no
  implementa login ni operaciones de dominio.
- **PostgreSQL:** disponible para el backend y para las migraciones futuras.
- **Android:** aún no existe una aplicación compilable y no forma parte de
  `compose.yaml`.

## Requisitos

- Docker Desktop o Docker Engine con Docker Compose.
- Puertos libres: `3001`, `8081` y `5433`.

Comprueba la instalación:

```bash
docker --version
docker compose version
```

## Configuración opcional

Compose incluye valores locales seguros por defecto. Para cambiar puertos o
credenciales de desarrollo, copia `.env.example` como `.env` en la raíz:

```powershell
Copy-Item .env.example .env
```

En Bash:

```bash
cp .env.example .env
```

El archivo `.env` está ignorado por Git. No guardes secretos reales en el
repositorio ni uses `apu_local_only` fuera de tu equipo.

## 1. Encender solo el frontend

El login actual no consume todavía la API, por lo que puede ejecutarse solo:

```bash
docker compose up --build dashboard
```

Abre <http://localhost:3001>. Para detenerlo, usa `Ctrl+C` y después:

```bash
docker compose down
```

Alternativa sin Compose:

```bash
docker build -t apu-dashboard ./dashboard
docker run --rm -p 3001:3000 apu-dashboard
```

## 2. Encender el backend

El backend necesita PostgreSQL. Compose inicia automáticamente la base de
datos al solicitar el servicio `backend`:

```bash
docker compose up --build backend
```

Comprueba sus endpoints:

```bash
curl http://localhost:8081/healthz
curl http://localhost:8081/readyz
curl http://localhost:8081/api/v1
```

En PowerShell también puedes usar:

```powershell
Invoke-RestMethod http://localhost:8081/healthz
Invoke-RestMethod http://localhost:8081/readyz
Invoke-RestMethod http://localhost:8081/api/v1
```

Para ejecutar la imagen del backend sin Compose, primero debes proporcionar
una URL de PostgreSQL accesible desde el contenedor:

```bash
docker build -t apu-backend ./backend
docker run --rm -p 8081:8080 \
  -e APU_ENV=local \
  -e APU_HTTP_ADDR=:8080 \
  -e APU_DATABASE_URL=postgres://apu:password@host.docker.internal:5432/apu?sslmode=disable \
  apu-backend
```

## 3. Encender todo el stack

Desde la raíz del repositorio:

```bash
docker compose up --build
```

Servicios disponibles:

| Servicio | Dirección | Propósito |
|---|---|---|
| Dashboard | <http://localhost:3001> | Login web de APU |
| API | <http://localhost:8081/api/v1> | Bootstrap de la API Go |
| Salud API | <http://localhost:8081/healthz> | Proceso disponible |
| Readiness API | <http://localhost:8081/readyz> | PostgreSQL alcanzable desde la API |
| PostgreSQL | `localhost:5433` | Base local de desarrollo |

Para ejecutarlo en segundo plano:

```bash
docker compose up --build -d
docker compose ps
docker compose logs -f
```

Para apagar sin borrar los datos:

```bash
docker compose down
```

Para apagar y borrar también la base local, ejecuta conscientemente:

```bash
docker compose down --volumes
```

Ese último comando elimina el volumen `postgres-data`; no se puede recuperar
desde Docker después de eliminarlo.

## Desarrollo sin Docker

### Backend Go

PowerShell:

```powershell
cd backend
$env:APU_ENV = "local"
$env:APU_HTTP_ADDR = ":8081"
$env:APU_DATABASE_URL = "postgres://apu:apu_local_only@localhost:5433/apu?sslmode=disable"
go run ./cmd/apu-api
```

Bash:

```bash
cd backend
APU_ENV=local \
APU_HTTP_ADDR=:8081 \
APU_DATABASE_URL='postgres://apu:apu_local_only@localhost:5433/apu?sslmode=disable' \
go run ./cmd/apu-api
```

### Frontend Next.js

```bash
cd dashboard
pnpm install
pnpm dev
```

## Comandos de diagnóstico

```bash
docker compose config
docker compose ps
docker compose logs backend
docker compose logs dashboard
docker compose logs postgres
```

Si cambias `NEXT_PUBLIC_API_BASE_URL`, reconstruye el dashboard porque las
variables `NEXT_PUBLIC_*` se incorporan durante el build:

```bash
docker compose build --no-cache dashboard
docker compose up -d dashboard
```
