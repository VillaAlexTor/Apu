# Backups y recuperación

## Objetivos

Los valores definitivos de RPO/RTO se fijarán por plan. Propuesta inicial:

| Servicio | RPO objetivo | RTO objetivo |
|---|---:|---:|
| Configuración y cuentas | 1 hora | 4 horas |
| Ubicaciones | 24 horas | 8 horas |
| Auditoría y facturación | 1 hora | 4 horas |

## Estrategia

- backups automáticos cifrados;
- copia separada del entorno principal;
- retención por política y obligación legal;
- claves de backup separadas y acceso mínimo;
- snapshots y backups lógicos según tamaño;
- verificación de integridad.

## Restauración

1. Declarar incidente y responsable.
2. Aislar el entorno afectado.
3. Seleccionar punto de recuperación.
4. Restaurar en entorno controlado.
5. Verificar migraciones, aislamiento por tenant y consistencia.
6. Reanudar servicios gradualmente.
7. Documentar pérdida, tiempos y acciones correctivas.

## Pruebas

Probar restauración al menos trimestralmente y después de cambios importantes.
Usar datos protegidos: nunca copiar una restauración de producción a desarrollo
sin controles de acceso y anonimización.

## Eliminación y solicitudes

Los backups pueden conservar datos hasta el fin de su ciclo técnico. La política
de eliminación debe documentar cómo se atienden solicitudes mientras existen
copias y qué excepciones legales aplican.

