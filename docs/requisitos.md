# Requisitos del proyecto

**Estado:** baseline del MVP  
**Versión:** 0.1  
**Última actualización:** 2026-11-04

## 1. Objetivo del MVP

Permitir que un tutor legal vincule un dispositivo Android familiar, active el
módulo de ubicación con conocimiento del menor y consulte su ubicación reciente
desde un panel web seguro.

El MVP no debe usarse con datos reales hasta completar la revisión legal y de
seguridad indicada en [modelo-amenazas.md](./modelo-amenazas.md).

## 2. Alcance incluido

- Registro, inicio de sesión, refresh y revocación de sesión del tutor.
- Creación de una familia y asociación del tutor.
- Generación y confirmación de códigos de emparejamiento.
- Consentimiento visible antes de activar módulos en Android.
- Registro y revocación de dispositivos.
- Activación y desactivación del módulo de ubicación.
- Recolección en Android con foreground service visible.
- Cola local y sincronización tolerante a conectividad intermitente.
- Historial de ubicación con retención limitada.
- Visualización de ubicación en mapa.
- Geocercas básicas y eventos de entrada/salida.
- Auditoría de accesos, cambios y comandos.

## 3. Fuera del MVP

- Monitoreo de mensajes o contenido de notificaciones.
- Bloqueo remoto de aplicaciones.
- `AccessibilityService`.
- Filtro web mediante VPN.
- Estadísticas de uso de aplicaciones.
- Múltiples roles avanzados de tutor.
- Publicación en Google Play.
- Uso comercial o multi-tenant público.

## 4. Requisitos funcionales

| ID | Requisito | Prioridad |
|---|---|---|
| RF-01 | El tutor puede registrarse con email y contraseña segura. | Must |
| RF-02 | El tutor puede iniciar, renovar y cerrar sesión. | Must |
| RF-03 | El sistema aísla los recursos por familia. | Must |
| RF-04 | El tutor puede generar un código temporal de emparejamiento. | Must |
| RF-05 | Android muestra consentimiento y módulos antes de vincularse. | Must |
| RF-06 | Android registra una clave en Keystore y envía solo la pública. | Must |
| RF-07 | El tutor puede activar o desactivar ubicación. | Must |
| RF-08 | Android envía lotes idempotentes de ubicación. | Must |
| RF-09 | El tutor puede consultar historial paginado. | Must |
| RF-10 | El sistema aplica retención automática. | Must |
| RF-11 | El tutor puede configurar geocercas. | Should |
| RF-12 | El panel recibe eventos de ubicación por WebSocket. | Should |
| RF-13 | El tutor puede revocar un dispositivo. | Must |
| RF-14 | Toda lectura sensible queda auditada. | Must |

## 5. Requisitos no funcionales

- **Seguridad:** TLS, tokens de corta duración, autorización por familia y
  secretos fuera del repositorio.
- **Privacidad:** minimización, consentimiento visible, cifrado y retención
  limitada.
- **Disponibilidad:** el cliente debe tolerar pérdida temporal de red.
- **Integridad:** reintentos no deben duplicar muestras ni comandos.
- **Observabilidad:** cada petición debe tener `request_id`.
- **Mantenibilidad:** dominio independiente de HTTP, ORM y proveedores.
- **Accesibilidad:** interfaces comprensibles, con estados y errores visibles.
- **Rendimiento:** paginar historiales y limitar lotes y respuestas.

## 6. Criterios de aceptación del MVP

El MVP está listo cuando un tutor puede completar el flujo de emparejamiento,
activar ubicación, perder y recuperar conectividad sin perder ni duplicar
eventos, consultar el historial y revocar el dispositivo. Además, los tests
demuestran que un tutor de otra familia no puede acceder a esos datos.

