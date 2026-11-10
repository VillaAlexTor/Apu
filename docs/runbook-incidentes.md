# Runbook de incidentes

## Principios

La prioridad es proteger a las personas y los datos, contener el impacto,
preservar evidencia mínima y comunicar con precisión. No se realizan pruebas
contra sistemas de terceros ni se incluyen coordenadas, tokens o contenido
sensible en tickets o chats.

## Severidad

| Nivel | Ejemplo | Respuesta inicial |
|---|---|---:|
| SEV1 | exposición de datos, toma de cuentas o caída total | 15 min |
| SEV2 | degradación importante o fallo de un módulo crítico | 1 h |
| SEV3 | impacto limitado con alternativa operativa | 1 día hábil |

## Flujo común

1. Abrir incidente privado y designar responsable.
2. Confirmar alcance, hora de inicio y servicios afectados.
3. Contener: revocar sesiones, deshabilitar endpoint, aislar servicio o pausar
   despliegue según corresponda.
4. Preservar logs redactados, métricas, hashes de artefactos y línea temporal.
5. Recuperar con el procedimiento de despliegue o restauración aprobado.
6. Verificar autorización, aislamiento, datos y métricas antes de reabrir.
7. Comunicar estado a las personas afectadas según la política aplicable.
8. Cerrar con causa raíz, impacto, acciones correctivas y fecha responsable.

## Escenarios

### Sospecha de filtración o acceso indebido

Revocar las credenciales afectadas, bloquear el vector, conservar evidencia y
activar la revisión de seguridad. No borrar logs ni modificar evidencia para
ocultar el incidente. Evaluar notificaciones legales antes de comunicar detalles.

### Compromiso de una cuenta

Suspender la cuenta, revocar sesiones y refresh tokens, revisar cambios de rol,
dispositivos y accesos auditados. Restablecer credenciales mediante un canal
verificado y exigir MFA cuando corresponda.

### Caída de API o base de datos

Confirmar si es despliegue, dependencia, capacidad o corrupción. Detener cambios,
consultar métricas y restaurar solo después de comprobar el punto de recuperación.
Validar aislamiento por tenant y consistencia antes de aceptar tráfico.

### FCM o comandos degradados

Marcar comandos como pendientes o fallidos, no reenviar indefinidamente y mostrar
el estado en el panel. Nunca presentar un comando como ejecutado sin acuse del
dispositivo.

## Evidencia mínima

Registrar identificador del incidente, servicios, versiones, timestamps UTC,
acciones, responsables, impacto estimado y enlaces a artefactos internos. Los
logs deben permanecer redactados.
