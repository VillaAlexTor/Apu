# Observabilidad

## Logs

Usar logs estructurados JSON con:

- timestamp UTC;
- nivel;
- servicio y versión;
- `request_id`;
- `trace_id` si existe;
- operación y resultado;
- duración;
- códigos técnicos no sensibles.

Nunca registrar contraseñas, tokens, claves privadas, coordenadas, payloads de
mensajes, códigos de emparejamiento ni headers de autorización.

## Métricas

Métricas iniciales:

- latencia y errores por endpoint;
- intentos de login y emparejamiento;
- lotes aceptados, duplicados y rechazados;
- tamaño y antigüedad de colas Android;
- conexiones WebSocket activas y reconexiones;
- comandos entregados, confirmados, expirados y fallidos;
- ejecución y errores del job de retención.

Las métricas no deben contener emails, IDs sensibles innecesarios ni coordenadas.

## Alertas

Alertar por:

- aumento de `5xx`;
- fallos de migración;
- rate limit anómalo;
- crecimiento sostenido de cola;
- errores de FCM;
- retención atrasada;
- backup fallido;
- intentos de acceso denegado anormales.

## Request ID y correlación

El servidor acepta un request ID válido o genera uno nuevo. Se devuelve al
cliente y se propaga a logs y eventos operacionales. Nunca se usa como secreto.

## Privacidad operacional

Producción debe separar logs de aplicación, auditoría y métricas. El acceso a
logs se limita por rol y se audita. Los datos se retienen el mínimo tiempo
necesario para diagnóstico y seguridad.

