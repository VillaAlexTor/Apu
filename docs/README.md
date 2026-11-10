# Documentación de APU

Este directorio contiene la documentación de producto, arquitectura, seguridad
y operación. Los documentos están escritos para poder evolucionar junto con el
código; cuando una decisión cambia, debe actualizarse la documentación afectada
en el mismo cambio.

## Mapa de documentación

### Producto y dominio

- [Requisitos](./requisitos.md)
- [Criterios de aceptación](./criterios-aceptacion.md)
- [Modelo de dominio](./modelo-dominio.md)
- [Glosario](./glosario-dominio.md)
- [Flujos](./flujos.md)
- [Usuarios y roles](./usuarios-roles.md)
- [Planes comerciales](./planes-comerciales.md)
- [Suscripciones y facturación](./suscripciones-facturacion.md)

### Arquitectura y contratos

- [Arquitectura general](./arquitectura.md)
- [Arquitectura SaaS](./arquitectura-saas.md)
- [Arquitectura del backend](./arquitectura-backend.md)
- [Arquitectura del dashboard](./arquitectura-dashboard.md)
- [Arquitectura Android](./arquitectura-android.md)
- [Base de datos](./base-datos.md)
- [Eventos](./eventos.md)
- [Contrato de API](./api/README.md)
- [OpenAPI](./api/openapi.yaml)
- [Decisiones ADR](./decisiones-adr.md)
- [Plantilla ADR](./plantilla-adr.md)

### Seguridad, privacidad y cumplimiento

- [Modelo de amenazas](./modelo-amenazas.md)
- [Matriz de autorización](./matriz-autorizacion.md)
- [Aislamiento de tenants](./tenant-isolation.md)
- [Autenticación y sesiones](./autenticacion-sesiones.md)
- [Privacidad y ciclo de vida de datos](./privacidad.md)
- [Consentimiento y onboarding](./consentimiento-onboarding.md)
- [Política de privacidad](./politica-privacidad.md)
- [Términos de servicio](./terminos-servicio.md)
- [Seguridad del repositorio](../SECURITY.md)

### Desarrollo y calidad

- [Guía de desarrollo](./guia-desarrollo.md)
- [Convenciones de código](./convenciones-codigo.md)
- [Configuración por entorno](./configuracion.md)
- [Estrategia de pruebas](./estrategia-pruebas.md)
- [Plan de implementación](./plan-implementacion.md)

### Operación

- [Despliegue](./despliegue.md)
- [Runbook de incidentes](./runbook-incidentes.md)
- [Backups y recuperación](./backup-recuperacion.md)
- [Observabilidad](./observabilidad.md)
- [Release y rollback](./release.md)
- [SLA y soporte](./sla-soporte.md)
- [Operación de datos](./operacion-datos.md)
- [Panel de administración](./panel-administracion.md)

## Estado de los documentos

La documentación describe el diseño previsto. Un documento no implica que la
funcionalidad exista: la implementación debe confirmarse mediante código,
pruebas y el estado del plan de implementación.

## Fuente de verdad

1. Los contratos públicos de API se definen en `api/openapi.yaml`.
2. Los permisos se definen en `matriz-autorizacion.md` y se verifican en el
   backend.
3. La retención y el tratamiento de datos se definen en `privacidad.md`.
4. Las decisiones que afecten varios componentes se registran en
   `decisiones-adr.md`.
