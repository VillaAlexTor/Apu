# Estrategia de pruebas

## Pirámide

1. **Unitarias:** validadores, casos de uso, reglas de geocerca, retención y
   transición de estados.
2. **Integración:** repositorios, migraciones, transacciones y autorización con
   PostgreSQL real de pruebas.
3. **Contrato:** API contra OpenAPI y cliente Android/panel.
4. **End-to-end:** login, emparejamiento, activar ubicación, sincronizar y ver
   historial.
5. **Seguridad:** aislamiento entre familias, replay, rate limit, secretos y
   redacción de logs.

## Casos obligatorios del MVP

- código de emparejamiento expirado, usado dos veces y con intentos excesivos;
- tutor de otra familia intentando leer o modificar un dispositivo;
- lote duplicado, lote parcialmente inválido y lote fuera de rango;
- pérdida de red antes y después de confirmar un lote;
- reloj del dispositivo desfasado;
- revocación de dispositivo y refresh token;
- vencimiento y purga de ubicación;
- geocerca con borde, precisión baja y coordenadas inválidas;
- reconexión del WebSocket sin duplicar eventos;
- migración desde base vacía y desde la versión anterior.

## Criterios de calidad

- no se aceptan tests que dependan de servicios externos reales;
- los tests de integración usan datos sintéticos;
- cada bug corregido debe añadir una regresión;
- el pipeline debe fallar ante formato, tipos, tests, dependencias vulnerables o
  secretos detectados;
- los umbrales de cobertura se fijarán después del primer incremento, separando
  cobertura de dominio de cobertura de transporte.

