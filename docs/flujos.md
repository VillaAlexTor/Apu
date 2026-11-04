# Flujos funcionales

## 1. Registro y primer acceso

1. El tutor envía email y contraseña.
2. El backend normaliza el email, valida la contraseña y crea la cuenta.
3. Se crea una familia inicial y se asocia el tutor.
4. Se emite una sesión de corta duración y refresh token rotativo.
5. El panel muestra el estado vacío y la opción de añadir dispositivo.

Errores: email existente, formato inválido, contraseña débil o rate limit.

## 2. Emparejamiento

1. El tutor solicita un código.
2. El backend guarda el hash, expiración e intentos permitidos.
3. El panel muestra texto y QR.
4. Android escanea o introduce el código.
5. Android muestra los módulos y el consentimiento.
6. Android genera la clave en Keystore.
7. Android confirma con código, clave pública y versión de aplicación.
8. El backend valida, crea el dispositivo y activa solo módulos consentidos.
9. El código se invalida y se audita el resultado.

Un código expirado, usado previamente o excedido en intentos no puede vincular
un dispositivo.

## 3. Activación de ubicación

1. El tutor activa el módulo desde el panel.
2. Android recibe la configuración y muestra la explicación del permiso.
3. El menor/tutor presente concede el permiso del sistema.
4. Android inicia el foreground service visible.
5. El estado pasa a `active`.
6. Si el permiso es rechazado, el módulo queda en `pending_permission` o `error`.

## 4. Sincronización

1. El recolector crea una muestra con `event_id`.
2. Room la guarda en la cola cifrada.
3. WorkManager agrupa muestras bajo límites de tamaño y antigüedad.
4. El cliente envía el lote con `Idempotency-Key`.
5. El backend autentica, verifica firma, valida rangos y persiste.
6. Las muestras aceptadas se eliminan de la cola local.
7. Duplicados se consideran confirmados; errores permanentes se muestran.
8. Errores transitorios se reintentan con backoff.

## 5. Consulta del historial

1. El panel solicita un rango acotado.
2. El backend resuelve la familia del tutor.
3. Verifica pertenencia del dispositivo.
4. Registra la lectura en auditoría.
5. Descifra solo los campos necesarios.
6. Devuelve una página con cursor.

## 6. Revocación

1. El tutor confirma la revocación.
2. El backend revoca credenciales, comandos pendientes y sesiones del dispositivo.
3. Android recibe la orden cuando vuelva a conectarse.
4. La app borra cola y datos locales, y detiene recolectores.
5. El backend programa la eliminación según retención.

## 7. Recuperación ante fallos

- Sin red: conservar una cola limitada y mostrar estado pendiente.
- Token expirado: renovar una vez; si falla, solicitar login.
- Servidor no disponible: backoff y no borrar datos pendientes.
- Payload inválido: marcar el evento y no reintentarlo indefinidamente.
- Clave perdida: exigir nuevo emparejamiento.

