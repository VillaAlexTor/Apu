# Modelo de amenazas de APU

**Metodología:** STRIDE  
**Alcance:** MVP de autenticación, emparejamiento, ubicación y panel.

## Activos

- identidad y sesiones de tutores;
- relación familiar y dispositivos vinculados;
- ubicación actual e histórica;
- códigos y claves de emparejamiento;
- configuración de módulos y reglas;
- registros de auditoría;
- credenciales FCM y secretos de infraestructura.

## Límites de confianza

1. Android ↔ Internet ↔ API.
2. Panel ↔ Internet ↔ API.
3. API ↔ PostgreSQL.
4. API ↔ Firebase/FCM.
5. Tutor/menor ↔ interfaces de sus respectivos dispositivos.

## Amenazas y controles

| ID | Categoría | Amenaza | Control requerido | Verificación |
|---|---|---|---|---|
| T-01 | Spoofing | Robo de sesión del tutor | access token corto, refresh rotativo, revocación y MFA opcional | tests de sesión y revocación |
| T-02 | Spoofing | Suplantación de dispositivo | clave en Keystore, firma de solicitudes y nonce/fecha | tests de firma inválida y replay |
| T-03 | Tampering | Alteración de lotes | TLS, validación de esquema, idempotencia y firma | tests de payload y duplicados |
| T-04 | Repudiation | Negación de acceso o comando | auditoría inmutable con actor, acción, objetivo y hora | tests de eventos de auditoría |
| T-05 | Information disclosure | IDOR entre familias | autorización centralizada y consultas con `family_id` | matriz de autorización |
| T-06 | Information disclosure | Filtración por logs | redacción de tokens, coordenadas y secretos | revisión automatizada de logs |
| T-07 | Denial of service | Fuerza bruta o abuso de API | rate limit, límites de lote y timeouts | pruebas de límites |
| T-08 | Elevation of privilege | Tutor accede a funciones no habilitadas | permisos por rol, módulo y familia | tests negativos |
| T-09 | Supply chain | Dependencia vulnerable | lockfiles, escaneo CI y actualización controlada | pipeline de seguridad |
| T-10 | Abuse | Uso contra un adulto sin consentimiento | onboarding visible, indicador persistente y prohibición documentada | revisión UX/legal |

## Requisitos mínimos antes de datos reales

- revisión de seguridad de autenticación y autorización;
- pruebas de replay y de aislamiento entre familias;
- secretos fuera de Git y rotación documentada;
- backups cifrados y prueba de restauración;
- borrado verificable por retención y desvinculación;
- revisión legal específica para la jurisdicción de despliegue;
- mecanismo de reporte responsable en `SECURITY.md`.

## Fuera de alcance inicial

No se aprueba todavía monitoreo de mensajes, bloqueo remoto, accesibilidad ni
VPN. Cada uno requiere un análisis STRIDE ampliado y decisión explícita.

