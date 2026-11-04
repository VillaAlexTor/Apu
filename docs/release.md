# Proceso de release

## Versionado

- Backend y dashboard: SemVer.
- Android: `versionName` SemVer y `versionCode` incremental.
- Cambios incompatibles de API requieren nueva versión o periodo de compatibilidad.

## Checklist previo

- requisitos y criterios de aceptación completos;
- tests, lint y análisis de seguridad aprobados;
- migraciones probadas desde la versión soportada;
- documentación y OpenAPI actualizados;
- cambios de privacidad y consentimiento revisados;
- notas de release escritas sin datos sensibles;
- artefactos firmados y hashes publicados por canal seguro.

## Release del backend y dashboard

1. Crear tag desde `main`.
2. Construir artefactos reproducibles.
3. Aplicar migraciones compatibles.
4. Desplegar staging y ejecutar smoke tests.
5. Aprobar producción.
6. Verificar métricas y logs.
7. Registrar rollback posible.

## Release Android

1. Revisar permisos y comportamiento visible.
2. Construir APK firmado fuera del repositorio.
3. Verificar hash y firma.
4. Probar instalación, actualización, revocación y pérdida de red.
5. Distribuir solo a dispositivos autorizados durante la fase experimental.

## Notas de release

Incluir cambios funcionales, correcciones, migraciones, riesgos conocidos y
acciones requeridas. No incluir coordenadas, identificadores reales, tokens ni
información de menores.

## Incidente posterior

Si una release causa riesgo de seguridad o privacidad, detener distribución,
revocar credenciales afectadas, preservar evidencia mínima y seguir el proceso
de reporte responsable de [SECURITY.md](../SECURITY.md).

