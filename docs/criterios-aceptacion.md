# Criterios de aceptación del MVP

El MVP se considera aceptado cuando todos los criterios siguientes pasan en
staging con datos sintéticos y pruebas automatizadas.

## Identidad y aislamiento

- Un usuario puede iniciar sesión, renovar y cerrar sesiones.
- Un usuario solo ve tenants, familias y dispositivos asignados.
- Un intento entre tenants o familias devuelve `403` o una respuesta equivalente
  sin filtrar existencia.
- Los cambios de rol y acciones sensibles quedan auditados.

## Emparejamiento

- El código expira, es de un solo uso y tiene límite de intentos.
- Android muestra consentimiento antes de activar módulos.
- La clave privada no abandona Keystore.
- Revocar un dispositivo invalida sus credenciales y comandos pendientes.

## Ubicación

- Solo se recolecta con módulo y permiso activos.
- Los lotes validan rango, firma, tamaño, antigüedad e `event_id`.
- Reintentar un lote no duplica datos.
- El panel muestra historial paginado, estado de sincronización y errores.
- La retención purga datos vencidos y registra el resultado sin coordenadas.

## Operación

- La API rechaza configuración obligatoria ausente y secretos por defecto.
- Logs y errores no contienen tokens, coordenadas ni secretos.
- Backups se restauran en un entorno controlado y se verifica el aislamiento.
- El pipeline ejecuta formato, lint, tests, análisis de dependencias y secretos.

## Definición de evidencia

Cada criterio debe enlazar a una prueba, captura sintética o resultado de CI.
Una pantalla que funciona pero no tiene autorización, auditoría y prueba
negativa no cumple el criterio.
