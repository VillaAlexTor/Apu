# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Panel web previsto con Next.js, React, TypeScript, Tailwind CSS, shadcn/ui y Zustand. El backend debe implementarse en Go (Golang), con autenticación mediante JWT y refresh tokens rotativos.

## Users

- Madres, padres y tutores legales que administran familias y dispositivos Android autorizados.
- Responsables de un tenant (`tenant_owner` y `tenant_admin`) que gestionan usuarios, configuración y, según su rol, facturación.
- Personal interno (`platform_admin` y `support_agent`) que opera tenants, soporte, facturación y salud del servicio sin acceso implícito a datos familiares sensibles.
- Todos los usuarios humanos acceden mediante un login compartido; después de autenticar, el backend determina el rol y el panel correspondiente.

## Product Purpose

APU es un servicio SaaS privado de control parental para Android. Permite administrar familias y dispositivos autorizados, emparejarlos con consentimiento visible, activar módulos de supervisión y consultar información desde un panel web seguro. El MVP se concentra en emparejamiento, ubicación, historial, geocercas, auditoría y operación multi-tenant.

El éxito del MVP consiste en que un tutor pueda vincular un dispositivo, activar ubicación, tolerar pérdidas temporales de conexión, consultar el historial y revocar el dispositivo sin duplicar eventos ni exponer datos entre familias o tenants.

## Positioning

APU plantea el control parental como acompañamiento transparente: el menor conoce qué módulos están activos y qué datos se comparten. Combina autoridad del tutor con consentimiento visible, minimización de datos, retención limitada y auditoría de accesos.

## Operating Context

- El tutor utiliza un panel web para administrar su familia, dispositivos y módulos permitidos.
- El personal interno utiliza paneles operativos separados y limitados para administrar tenants, soporte, facturación y salud del servicio.
- La aplicación Android se instala inicialmente como APK de pruebas en dispositivos propios de la familia.
- El producto debe seguir siendo comprensible durante conectividad intermitente y mostrar estados pendientes, errores de permisos y sincronización.
- APU se encuentra en fase de diseño y planificación técnica y todavía no debe utilizarse con datos reales de menores.

## Capabilities and Constraints

- Autenticación con email y contraseña, sesiones cortas, refresh tokens rotativos y MFA para roles administrativos.
- Autorización y aislamiento aplicados en el backend por `tenant_id`, familia y rol; la interfaz nunca concede permisos por sí sola.
- Registro, emparejamiento temporal por código o QR, gestión de dispositivos, ubicación, historial, geocercas y auditoría.
- Los operadores no pueden consultar coordenadas ni historiales familiares sensibles por defecto.
- El login es compartido por tutores, responsables de tenant y operadores; el backend resuelve el destino autorizado después de autenticar.
- Internacionalización futura prevista para español, inglés, quechua y aimara.
- Monitoreo de mensajes, bloqueo de aplicaciones, estadísticas de uso y publicación en Google Play quedan fuera del MVP.

## Brand Commitments

- Nombre: **APU**.
- Significado: en la cosmovisión andina, un Apu es el espíritu protector de las montañas que cuida a la comunidad.
- Origen declarado: La Paz, Bolivia.
- La comunicación debe sostener transparencia, protección, proporcionalidad y acompañamiento familiar, evitando presentar el producto como vigilancia secreta.

## Evidence on Hand

- Documentación de producto, arquitectura, seguridad, privacidad, roles, flujos y operación en `README.md` y `docs/`.
- No existen todavía una implementación visual, testimonios, métricas comerciales, capturas del producto ni activos de marca confirmados fuera del nombre y su significado.
- No se deben inventar clientes, certificaciones, resultados, cumplimiento legal concluido ni garantías de seguridad.

## Product Principles

1. Hacer visible y comprensible qué se monitorea y por qué.
2. Proteger datos familiares mediante mínimo privilegio, aislamiento y auditoría.
3. Acompañar a la familia sin normalizar la vigilancia encubierta.
4. Mantener los estados, permisos y consecuencias claros para personas no técnicas.
5. Separar estrictamente la operación del SaaS del acceso a datos familiares sensibles.

## Accessibility & Inclusion

Las interfaces deben ser comprensibles, accesibles y mostrar estados y errores de forma visible. La dirección futura contempla español, inglés, quechua y aimara, por lo que los layouts no deben depender de textos de longitud fija.
