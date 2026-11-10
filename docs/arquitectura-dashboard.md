# Arquitectura del dashboard

## Responsabilidad

El dashboard permite a tutores y responsables operar recursos autorizados.
Muestra estados, errores, permisos y sincronización de forma comprensible, pero
no es una frontera de seguridad: toda autorización ocurre en la API.

## Estructura

```text
dashboard/src/
├─ app/              # rutas, layouts y carga protegida
├─ features/         # auth, devices, locations, settings, billing
├─ components/       # componentes visuales compartidos
├─ lib/              # cliente HTTP, sesión, validación y configuración
└─ stores/           # estado de UI y sesión; no fuente de permisos
```

El cliente HTTP tipado centraliza base URL, headers, renovación controlada,
`request_id`, errores y cancelación. Las vistas no construyen URLs ni tokens
manualmente.

## Sesión y autorización

La interfaz oculta acciones no disponibles para mejorar UX, pero debe asumir
que cualquier respuesta puede ser `401` o `403`. Un `401` inicia una renovación
única; un `403` se muestra como falta de autorización y no se reintenta.
Cambios sensibles requieren confirmación explícita y reflejan el estado del
servidor.

## Datos sensibles

No persistir coordenadas ni respuestas sensibles en almacenamiento permanente
del navegador salvo una decisión documentada. Limpiar estado al cerrar sesión o
cambiar de tenant. Las herramientas de desarrollo y errores no deben mostrar
tokens ni payloads privados.

## Tiempo real y accesibilidad

WebSocket solo invalida o actualiza cachés ya autorizadas; la API REST sigue
siendo la fuente de verdad. La interfaz debe funcionar con conectividad
intermitente, estados de carga/vacío/error visibles, teclado, lectores de
pantalla y textos que puedan crecer por internacionalización.
