# Planes comerciales

Los valores económicos quedan pendientes de decisión. Los límites técnicos sí
deben modelarse desde el inicio para impedir que un cliente exceda su plan.

| Capacidad | Básico | Familiar | Profesional |
|---|---:|---:|---:|
| Usuarios del tenant | 1 | 5 | 20 |
| Dispositivos | 2 | 10 | 50 |
| Historial de ubicación | 7 días | 30 días | 90 días |
| Geocercas | 3 | 10 | 50 |
| Soporte | estándar | prioritario | prioritario |
| API/WebSocket | sí | sí | sí |
| Módulos sensibles | sujeto a aprobación | sujeto a aprobación | sujeto a aprobación |

Estos nombres y límites son una propuesta, no una tarifa publicada.

## Reglas de capacidad

- El plan vigente se consulta en el backend.
- Superar un límite bloquea la creación nueva y explica cómo resolverlo.
- Reducir de plan requiere validar que los recursos actuales caben en el nuevo
  límite o programar la reducción al siguiente periodo.
- Nunca se eliminan dispositivos automáticamente por un cambio de plan.
- Las excepciones manuales tienen expiración y auditoría.

## Módulos sensibles

Mensajes, accesibilidad, VPN y bloqueo remoto no se habilitan automáticamente
por pagar un plan. Requieren aprobación legal, técnica y de políticas del
producto.

