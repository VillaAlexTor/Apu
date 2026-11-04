# Suscripciones y facturación

**Nota:** este documento define comportamiento técnico, no constituye una oferta
comercial ni asesoría tributaria.

## Ciclo de suscripción

```text
prospect → trial → active → past_due → suspended → terminated
                         └→ cancelled (vigente hasta fin de periodo)
```

El proveedor de pagos es autoridad para el estado de cobro. APU mantiene una
réplica mínima: proveedor, customer reference, subscription reference, plan,
estado, periodo actual y timestamps.

## Reglas

- No almacenar números completos de tarjeta ni códigos de seguridad.
- Procesar webhooks con firma, timestamp y deduplicación por `event_id`.
- No activar beneficios por una respuesta de navegador no verificada.
- Una renovación fallida inicia un periodo de gracia configurable.
- Cancelar impide la siguiente renovación, pero conserva acceso hasta el fin
  del periodo pagado salvo fraude o riesgo de seguridad.
- Cambios de plan deben registrar precio, moneda, fecha efectiva y actor.
- Reembolsos y ajustes requieren autorización del operador.

## Integración

La interfaz de facturación debe encapsular el proveedor para poder cambiar entre
Stripe, Mercado Pago u otro servicio sin contaminar el dominio. El webhook
actualiza una máquina de estados idempotente y publica eventos internos.

## Datos y privacidad

Solo se almacenan referencias necesarias para operar la cuenta. Las facturas y
datos fiscales se conservan según obligaciones legales aplicables y separados de
los datos de ubicación.

## Estados de acceso

El backend debe consultar el estado local verificado antes de permitir nuevas
operaciones. Los jobs de reconciliación detectan divergencias entre proveedor y
APU y las dejan visibles para soporte.

