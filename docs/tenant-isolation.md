# Aislamiento entre tenants

## Regla principal

Ningún dato de un cliente puede ser leído, inferido, actualizado o eliminado
por otro cliente. El identificador de tenant se deriva del contexto autenticado.

## Diseño de datos

Las tablas de negocio incluyen `tenant_id` directa o indirectamente:

```text
tenant
├── tenant_users
├── subscription
├── families (tenant_id)
└── devices (tenant_id, family_id)
```

Las consultas de dispositivos deben filtrar tanto `tenant_id` como la relación
de familia. Los índices compuestos empiezan por `tenant_id` cuando sea parte
habitual del filtro.

## Capas de defensa

1. Claims de sesión validados y tenant activo.
2. Middleware que construye el contexto, no que acepta el tenant del cliente.
3. Casos de uso que comprueban rol, tenant y pertenencia del recurso.
4. Repositorios con filtros obligatorios.
5. Restricciones y políticas de base de datos cuando sean viables.
6. Tests de lectura, escritura, exportación, WebSocket y tareas asíncronas.

## Riesgos específicos

- IDOR por cambiar un UUID en la URL.
- Worker que procesa un ID sin contexto de tenant.
- Caché compartida sin clave de tenant.
- Eventos WebSocket enviados a conexiones incorrectas.
- Exportaciones o backups mezclados.
- Métricas que revelan nombres o volúmenes de otro cliente.

## Criterio de aceptación

Para cada endpoint y job existe una prueba negativa que usa un recurso válido de
otro tenant y demuestra respuesta segura (`404` o `403` según política), sin
efectos secundarios ni filtración de existencia.

