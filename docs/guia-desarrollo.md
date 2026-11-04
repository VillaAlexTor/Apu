# Guía de desarrollo

## Requisitos locales

- Go 1.22 o superior.
- Node.js 20 o superior y pnpm.
- JDK 17 y Android Studio.
- PostgreSQL 15 o Docker.
- Git.

## Preparación

```bash
git clone https://github.com/VillaAlexTor/Apu.git
cd Apu
docker compose up -d postgres
```

Cada aplicación tendrá un `.env.example`. Los archivos `.env`, claves Firebase,
keystores y credenciales nunca se versionan.

## Flujo diario

```bash
# Backend
cd backend
go test ./...
go run ./cmd/apu-server

# Dashboard
cd dashboard
pnpm install
pnpm lint
pnpm dev

# Android
cd android
./gradlew test
./gradlew assembleDebug
```

## Migraciones

1. Crear una migración numerada y reversible.
2. Actualizar el modelo y repositorio en el mismo cambio.
3. Probar desde base vacía y desde la última versión.
4. Documentar impacto, índices y retención.
5. No editar una migración ya aplicada.

## Ramas y commits

- `main`: rama protegida y publicable.
- `feat/<nombre>`: funcionalidad.
- `fix/<nombre>`: corrección.
- `docs/<nombre>`: documentación.
- Commits con Conventional Commits.
- Pull request con pruebas, impacto de datos y riesgos.

## Datos de desarrollo

Usar únicamente datos sintéticos. No copiar emails, coordenadas, tokens,
capturas ni bases de datos reales al entorno local.

## Checklist antes de solicitar revisión

- formato y lint;
- tests unitarios e integración afectados;
- migraciones verificadas;
- OpenAPI actualizado;
- autorización y auditoría cubiertas;
- logs sin datos sensibles;
- documentación y variables de entorno actualizadas.

