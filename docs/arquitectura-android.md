# Arquitectura Android

## Responsabilidad

La aplicación Android muestra consentimiento, solicita permisos, recolecta solo
los módulos activos y sincroniza datos con tolerancia a conectividad
intermitente. El dispositivo no usa rutas destinadas a usuarios del panel.

## Capas

```text
Compose UI -> ViewModel -> casos de uso -> repositorios
                                      -> API / Room / Android APIs
```

- La UI muestra módulo, permiso, sincronización y errores.
- Los casos de uso aplican estados y coordinan permisos.
- Room conserva una cola local limitada y protegida.
- WorkManager sincroniza con backoff y límites.
- El foreground service se reserva para ubicación activa y notificación visible.

## Identidad y emparejamiento

La clave privada se genera y permanece en Android Keystore. El backend recibe
solo la clave pública. El código temporal se usa una sola vez; perder la clave,
desinstalar o desvincular exige nuevo emparejamiento.

## Recolección y sincronización

Cada evento incluye `event_id`, timestamp, precisión y versión de esquema.
Validar antes de encolar, agrupar con límites y eliminar solo tras confirmación
del servidor. Los reintentos deben ser seguros ante respuestas repetidas,
reinicios y reloj desfasado.

## Privacidad y estados

No iniciar un recolector con el módulo desactivado o el permiso ausente. Mostrar
qué se recolecta y durante cuánto tiempo. Al desvincular, revocar credenciales,
detener workers y borrar la cola local.

## Pruebas

Cubrir permisos denegados, pérdida de red, reinicio, batería baja, lotes
duplicados, firma inválida, revocación y migraciones de Room. Usar emuladores o
dispositivos propios con datos sintéticos.
