# Arquitectura del backend

## Responsabilidad

La API Go es la autoridad para identidad, autorización, estado de suscripción,
emparejamiento, ubicación, reglas, auditoría y eventos. Los clientes no deciden
pertenencia, permisos, retención ni transiciones de estado.

## Capas

```text
HTTP/middleware -> DTO y validación -> caso de uso -> dominio
                                      -> repositorio/servicio
```

- **Transporte:** rutas, autenticación, límites, serialización y errores.
- **Aplicación:** casos de uso, autorización y transacciones.
- **Dominio:** entidades, invariantes y estados sin dependencias externas.
- **Infraestructura:** PostgreSQL, FCM, reloj, criptografía y observabilidad.

La lógica de dominio no importa Gin/Chi, GORM, Firebase ni tipos HTTP.

## Módulos previstos

`auth`, `tenant`, `family`, `device`, `pairing`, `location`, `geofence`,
`subscription`, `billing`, `audit`, `push`, `storage` y `platform`.

Cada módulo debe definir sus endpoints, autorización, eventos, migraciones,
retención, métricas y pruebas. Los handlers no deben acceder directamente a
tablas de otro módulo.

## Reglas de seguridad

- Resolver el actor desde la sesión autenticada, nunca desde IDs confiados.
- Comprobar `tenant_id`, familia, dispositivo, rol y módulo en el caso de uso.
- Usar consultas parametrizadas y transacciones para cambios relacionados.
- Redactar tokens, coordenadas, secretos y payloads sensibles de logs.
- Hacer idempotentes los lotes, webhooks y comandos reintentables.
- Propagar errores explícitos con códigos estables y `request_id`.

## Persistencia y eventos

PostgreSQL es la fuente de verdad. Los eventos derivados son notificadores y
deben poder reconstruirse desde el estado persistido. WebSocket y FCM no
conceden permisos ni sustituyen una lectura autorizada por REST.

## Pruebas

Cada caso de uso requiere pruebas de autorización positiva y negativa,
idempotencia, transiciones de estado y errores. Los repositorios se prueban con
PostgreSQL real de pruebas; los servicios externos usan adaptadores verificables.
