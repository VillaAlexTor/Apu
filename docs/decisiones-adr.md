# Decisiones de arquitectura

**Estado:** propuesta inicial  
**Regla:** una decisión aceptada debe actualizar este documento y la
documentación afectada antes de implementarse.

## ADR-001: Monorepo con tres aplicaciones

- **Estado:** aceptada para el MVP.
- **Decisión:** mantener backend Go, dashboard Next.js y Android Kotlin en el
  mismo repositorio.
- **Motivo:** permite versionar contratos, migraciones y documentación juntos,
  manteniendo límites de compilación independientes.
- **Consecuencia:** CI debe ejecutar validaciones por carpeta y no mezclar
  dependencias de los tres ecosistemas.

## ADR-002: API REST versionada y WebSocket complementario

- **Estado:** aceptada para el MVP.
- **Decisión:** usar REST JSON bajo `/api/v1`; WebSocket solo para eventos del
  panel y no como canal de autoridad.
- **Motivo:** REST es fácil de depurar, probar y consumir desde Android y web;
  WebSocket aporta actualización de ubicación y estados.
- **Consecuencia:** cualquier evento recibido por WebSocket debe poder recuperarse
  por REST y no debe conceder permisos.

## ADR-003: Autorización por familia en casos de uso

- **Estado:** obligatoria.
- **Decisión:** la pertenencia `usuario → familia → dispositivo` se verifica en
  el caso de uso, además de los filtros del repositorio.
- **Motivo:** evita IDOR aunque un handler nuevo olvide aplicar un filtro.
- **Consecuencia:** cada caso de uso recibe un actor autenticado y debe tener
  pruebas de acceso permitido y denegado.

## ADR-004: Room como cola local de Android

- **Estado:** aceptada para el MVP.
- **Decisión:** usar Room para eventos pendientes y estado mínimo de sincronización.
- **Motivo:** WorkManager puede reanudar envíos tras reinicios o conectividad
  intermitente sin perder muestras.
- **Consecuencia:** la cola debe cifrarse, tener límite de tamaño y políticas de
  eliminación explícitas.

## ADR-005: Claves de dispositivo en Android Keystore

- **Estado:** aceptada.
- **Decisión:** generar la clave privada en Android Keystore y registrar solo la
  clave pública en el backend.
- **Motivo:** la clave privada no debe salir del dispositivo ni almacenarse en
  preferencias o archivos de aplicación.
- **Consecuencia:** reinstalar o perder la clave exige un nuevo emparejamiento.

## ADR-006: Ubicación cifrada y retención configurable

- **Estado:** propuesta para validar legalmente.
- **Decisión:** cifrar la ubicación en reposo a nivel de campo y aplicar una
  retención por defecto corta, configurable dentro de límites seguros.
- **Motivo:** una filtración de base de datos no debe revelar historiales legibles.
- **Consecuencia:** las consultas de mapa requieren descifrado controlado y
  auditoría; no se usarán coordenadas en logs.

## ADR-007: Evolución por módulos

- **Estado:** aceptada.
- **Decisión:** ubicación, uso, límites, filtro web y mensajes se modelan como
  módulos independientes.
- **Motivo:** cada módulo tiene permisos, riesgos, retención y consentimiento
  diferentes.
- **Consecuencia:** activar un módulo requiere migración, contrato, UI de estado,
  política de datos, auditoría y tests propios.

## ADR-008: Mensajes fuera del MVP

- **Estado:** aceptada.
- **Decisión:** no implementar lectura de notificaciones hasta completar revisión
  legal, modelo de amenazas, consentimiento y pruebas de minimización.
- **Motivo:** es la capacidad de mayor sensibilidad y riesgo de abuso.
- **Consecuencia:** no se deben crear permisos, tablas ni endpoints de mensajes
  como parte del primer incremento.

