# SLA y soporte

**Estado:** propuesta operativa; los objetivos finales dependen del plan.

## Objetivos iniciales

| Plan | Disponibilidad objetivo mensual | Primera respuesta |
|---|---:|---:|
| Básico | 99.0% | 2 días hábiles |
| Familiar | 99.5% | 1 día hábil |
| Profesional | 99.9% | 4 horas hábiles |

Son objetivos de diseño, no compromisos contractuales hasta publicarse en los
términos del servicio.

## Severidad

- **P0:** riesgo activo de seguridad, fuga entre tenants o servicio totalmente
  indisponible.
- **P1:** función principal degradada para muchos clientes.
- **P2:** fallo de una función con alternativa disponible.
- **P3:** consulta, mejora o problema cosmético.

## Soporte

El ticket incluye tenant, versión, hora, request ID, impacto y evidencia
sintética. Nunca se solicitan contraseñas, tokens, coordenadas ni contenido de
mensajes por canales de soporte.

## Mantenimiento

Las ventanas planificadas se anuncian con anticipación. Cambios urgentes por
seguridad pueden ejecutarse fuera de ventana y deben quedar documentados.

## Medición

La disponibilidad excluye mantenimiento anunciado y fallos causados por red,
permisos, batería o proveedores del cliente, según contrato. Las métricas se
calculan desde monitores independientes y se conservan sin datos sensibles.

