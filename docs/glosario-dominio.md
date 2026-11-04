# Glosario del dominio

| Término | Definición |
|---|---|
| **Tutor** | Persona autorizada legalmente para administrar una familia en APU. |
| **Menor** | Persona cuyo dispositivo familiar está vinculado a APU. |
| **Familia** | Límite lógico de aislamiento que contiene tutores y dispositivos. |
| **Dispositivo** | Instalación de APU en Android vinculada a una familia. |
| **Módulo** | Capacidad independiente que puede activarse por dispositivo. |
| **Ubicación** | Muestra de coordenadas, precisión y hora registrada por el dispositivo. |
| **Lote** | Conjunto de muestras enviadas en una petición idempotente. |
| **Evento** | Hecho de dominio derivado, como entrada a una geocerca. |
| **Comando** | Instrucción autorizada que el backend entrega al dispositivo. |
| **Emparejamiento** | Proceso de asociación autorizada entre tutor, familia y dispositivo. |
| **Consentimiento** | Aceptación visible de los módulos y datos que se compartirán. |
| **Geocerca** | Área geográfica configurada para detectar entrada o salida. |
| **Sesión** | Relación autenticada entre un tutor y el backend. |
| **Auditoría** | Registro de una acción sensible, actor, objetivo y momento. |
| **Retención** | Tiempo máximo que un tipo de dato puede conservarse. |
| **Revocación** | Invalidación de una sesión, dispositivo, clave o comando. |

## Reglas de vocabulario

- Usar `tutor`, no `parent`, en nombres de dominio.
- Usar `device_id`, no `child_id`: el recurso técnico es el dispositivo.
- Usar `recorded_at` para el momento de captura y `created_at` para persistencia.
- Usar `family_id` en toda consulta de autorización.
- Reservar “tiempo real” para actualizaciones entregadas mientras hay conexión;
  el historial no promete latencia constante.

