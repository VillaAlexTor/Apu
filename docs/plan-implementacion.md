# Plan técnico de implementación

Este plan ordena el trabajo para no comenzar por funcionalidades sensibles ni
por pantallas sin contratos definidos.

## Fase 0 — Fundaciones y decisiones

- validar arquitectura, amenazas, privacidad y jurisdicción;
- confirmar licencia y política de contribución;
- crear estructura de monorepo y documentación OpenAPI inicial;
- configurar formato, lint, análisis estático y CI sin secretos;
- definir configuración por entorno y estrategia de migraciones.

**Salida:** repositorio compilable aunque no tenga funcionalidades de negocio.

## Fase 1 — Identidad, tenants y roles

- registro, login, refresh y revocación;
- modelo de tenants, usuarios, roles y relación tutor-familia;
- middleware de autenticación;
- auditoría base;
- pruebas de aislamiento.

**Salida:** un usuario autenticado solo puede operar dentro del tenant y las
familias autorizadas.

## Fase 2 — Suscripciones y emparejamiento

- planes y límites;
- estados de suscripción;
- webhooks de pago idempotentes;
- panel de administración básico;

- código hashado con TTL y rate limit;
- flujo Android visible y consentimiento;
- clave en Keystore;
- registro y revocación del dispositivo;
- contrato y tests de replay.

**Salida:** tenant con acceso comercial verificable y dispositivo vinculado de
manera revocable y auditable.

## Fase 3 — Ubicación y mapa

- permisos y foreground service;
- Room + WorkManager;
- endpoint idempotente de lotes;
- historial paginado;
- WebSocket de eventos;
- retención y geocercas.

**Salida:** MVP funcional con conectividad intermitente.

## Fase 4 — Tiempo de uso y límites

Se inicia solo después de revisar consumo de batería, permisos y políticas de
Android. Cada módulo tendrá su propia especificación, consentimiento y análisis
de amenazas.

## Fase 5 — Capacidades de alta sensibilidad

Mensajes, accesibilidad, VPN y bloqueo remoto requieren aprobación explícita,
revisión legal y de seguridad. No se implementan como atajos del MVP.

## Definition of Ready

Una tarea puede entrar a desarrollo solo si tiene:

- objetivo y actor definidos;
- contrato o interfaz;
- estados y errores;
- autorización y auditoría;
- impacto de datos y retención;
- casos de prueba;
- criterio de aceptación observable.

## Definition of Done

Una tarea está terminada cuando:

- compila y pasa formato, lint y tests;
- incluye migración y rollback si corresponde;
- actualiza OpenAPI y documentación;
- no introduce secretos ni datos reales;
- tiene logs seguros y errores explícitos;
- se verificó en el flujo afectado.
