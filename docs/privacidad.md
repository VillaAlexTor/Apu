# Privacidad y ciclo de vida de datos

Este documento describe controles técnicos previstos para el SaaS. No sustituye
la [política de privacidad](./politica-privacidad.md), una política legal ni
asesoría jurídica. APU no debe usarse con datos reales de menores hasta completar
la revisión legal y de seguridad.

## 1. Clasificación

| Datos | Sensibilidad | Uso | Retención inicial |
|---|---|---|---|
| Email y cuenta del tutor | alta | autenticación y familia | mientras exista la cuenta |
| Identidad técnica del dispositivo | alta | emparejamiento y sincronización | mientras esté vinculado |
| Coordenadas | muy alta | mapa, historial y geocercas | 7 días por defecto |
| Estado de batería/conectividad | media | diagnóstico y panel | 7 días |
| Auditoría | alta | seguridad y responsabilidad | 90 días, sujeto a revisión |
| Tokens y secretos | crítica | autenticación/integraciones | nunca en logs; rotación |

Los datos se clasifican y retienen dentro del tenant correspondiente. Los datos
de facturación se mantienen separados de la ubicación y sujetos a obligaciones
legales distintas.

Los mensajes no forman parte del MVP y no deben añadirse a este esquema sin una
decisión específica.

## 2. Recolección y consentimiento

- El onboarding explica qué se recoge, para qué, durante cuánto tiempo y quién
  puede verlo.
- Cada módulo tiene estado visible: `disabled`, `pending_permission`, `active`,
  `error`.
- El permiso de Android se solicita solo después de explicar su finalidad.
- La app muestra un indicador persistente cuando una función de monitoreo está
  activa.
- Desactivar un módulo detiene la recolección futura; no borra automáticamente
  datos históricos salvo que la política del módulo lo indique.

## 3. Acceso y aislamiento

- El usuario solo accede a dispositivos de su familia dentro de su tenant.
- El backend no confía en IDs enviados por el cliente.
- Toda lectura de ubicación queda auditada con actor, dispositivo, rango de
  tiempo y resultado.
- Soporte y administradores no reciben acceso operativo a coordenadas descifradas.
- Las respuestas omiten datos no necesarios para la pantalla solicitante.

## 4. Retención y eliminación

El proceso de retención debe:

1. seleccionar datos vencidos por tipo;
2. eliminar en lotes pequeños para no bloquear el servicio;
3. registrar cantidad y tipo, nunca el contenido;
4. ser idempotente y reanudable;
5. emitir una métrica de éxito o error.

Desvincular un dispositivo revoca credenciales, detiene comandos pendientes,
borra la cola local y programa la eliminación de datos según la política
configurada. La eliminación lógica solo puede usarse si tiene un proceso de
purga posterior verificable.

## 5. Derechos y operaciones futuras

Antes del despliegue público se deben definir procedimientos para acceso,
corrección, exportación y eliminación, junto con responsables y tiempos de
respuesta. También debe definirse la jurisdicción aplicable y el rol legal del
operador del servicio.

## 6. Reglas de implementación

- No incluir datos personales en URLs, nombres de archivos, logs ni mensajes de
  error.
- Cifrar secretos mediante un gestor de secretos.
- Separar claves por entorno y rotarlas.
- No usar datos de producción en desarrollo o tests.
- Usar datos sintéticos en fixtures y capturas.
