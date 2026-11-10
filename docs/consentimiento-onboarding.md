# Consentimiento y onboarding

## Objetivo

El onboarding debe permitir que el tutor y la persona que usa el dispositivo
entiendan qué módulo está activo, qué datos se comparten, con quién y durante
cuánto tiempo. No se permite activación silenciosa ni camuflaje de funciones.

## Flujo mínimo

1. El tutor crea un código temporal desde el panel.
2. Android muestra la identidad del servicio y solicita introducir o escanear el
   código.
3. La aplicación muestra módulos, datos, retención y destinatarios.
4. La persona confirma cada módulo y Android solicita solo sus permisos.
5. Se registra versión del texto, módulos aceptados, timestamp y dispositivo.
6. La app muestra un indicador persistente de las funciones activas.

## Estados

`not_paired → pairing_pending → paired → permission_pending → active`,
con estados de error visibles y recuperación guiada.

## Cambios y revocación

Activar un módulo posteriormente requiere explicación y confirmación nuevas.
Desactivarlo detiene la recolección futura y actualiza el indicador. Desvincular
revoca credenciales, detiene sincronización y programa la eliminación conforme a
la política de datos.

El consentimiento técnico no sustituye la evaluación legal de autoridad,
edad, jurisdicción ni responsabilidad del tutor. Antes del uso real debe existir
un procedimiento legal revisado.
