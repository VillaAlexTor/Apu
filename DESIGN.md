---
name: APU
description: Protección transparente expresada como un atlas topográfico andino.
colors:
  mineral-paper: "#F3EBD8"
  paper-light: "#FAF5E9"
  forest-ink: "#203B35"
  forest-deep: "#142A26"
  mineral-red: "#C75A3A"
  survey-yellow: "#E1B84B"
  river-blue: "#447D88"
  error-red: "#9F382B"
  survey-line: "rgba(32, 59, 53, 0.34)"
typography:
  display:
    fontFamily: '"Archivo Variable", Archivo, sans-serif'
    fontSize: "clamp(2rem, 3.1vw, 3.65rem)"
    fontWeight: 650
    lineHeight: 0.98
    letterSpacing: "-0.055em"
  wordmark:
    fontFamily: '"Barlow Condensed", sans-serif'
    fontSize: "clamp(4.5rem, 7vw, 7.75rem)"
    fontWeight: 700
    lineHeight: 0.75
    letterSpacing: "-0.055em"
  body:
    fontFamily: '"Archivo Variable", Archivo, sans-serif'
    fontSize: "clamp(0.88rem, 1vw, 1.02rem)"
    fontWeight: 400
    lineHeight: 1.45
  label:
    fontFamily: '"Archivo Variable", Archivo, sans-serif'
    fontSize: "0.82rem"
    fontWeight: 650
  coordinate:
    fontFamily: '"Barlow Condensed", sans-serif'
    fontSize: "0.68rem"
    fontWeight: 500
    lineHeight: 1
    letterSpacing: "0.13em"
rounded:
  square: "0px"
  marker: "50%"
spacing:
  control-x: "1.05rem"
  field-y: "0.8rem"
  legend-inset: "clamp(1.35rem, 3vw, 3rem)"
  frame-inset: "clamp(1rem, 2.2vw, 2.25rem)"
components:
  button-primary:
    backgroundColor: "{colors.forest-ink}"
    textColor: "{colors.paper-light}"
    typography: "{typography.label}"
    rounded: "{rounded.square}"
    height: "3.25rem"
    padding: "0 1.05rem 0 1.2rem"
  button-primary-hover:
    backgroundColor: "{colors.survey-yellow}"
    textColor: "{colors.forest-deep}"
    rounded: "{rounded.square}"
  input:
    backgroundColor: "rgba(255, 252, 244, 0.64)"
    textColor: "{colors.forest-deep}"
    typography: "{typography.body}"
    rounded: "{rounded.square}"
    height: "3.15rem"
    padding: "0.8rem 3rem 0.8rem 0.9rem"
  login-legend:
    backgroundColor: "rgba(250, 245, 233, 0.92)"
    textColor: "{colors.forest-ink}"
    rounded: "{rounded.square}"
    padding: "clamp(1.35rem, 2.6vw, 2.4rem) clamp(1.35rem, 3vw, 3rem)"
---

# Design System: APU

## Overview

**Creative North Star: "El Atlas Protector"**

APU se comporta como un atlas topográfico andino que acompaña y orienta. El papel mineral, las reglas de levantamiento y las marcas de coordenadas dan una presencia física y territorial; la información permanece directa, legible y explícita porque la protección nunca debe confundirse con vigilancia oculta.

La interfaz es plana, precisa y táctil. Su identidad nace del contraste entre campos tranquilos de papel y evidencia cartográfica densa: curvas, ríos, crestas, rutas y puntos de registro. Rechaza la tarjeta SaaS genérica, los degradados tecnológicos sin contexto y la decoración que no explique estado, recorrido o jerarquía.

**Key Characteristics:**

- Papel mineral cálido como material continuo, no como fondo crema genérico.
- Verde bosque estructural para texto, marcos, controles y límites.
- Acentos cartográficos escasos: amarillo de levantamiento, rojo mineral y azul río.
- Tipografía funcional con una capa condensada para coordenadas, sellos y metadatos.
- Geometría recta, muescas y marcas de registro en lugar de tarjetas redondeadas.
- Estados visibles y lenguaje transparente para una audiencia no técnica y multilingüe.

## Colors

La paleta proviene de un mapa de campo impreso: sobria en superficie y específica en sus señales.

### Primary

- **Tinta Bosque:** sostiene texto, reglas, marcos, campos y acciones principales; concentra autoridad sin recurrir al negro puro.
- **Bosque Profundo:** refuerza titulares, wordmark y estados de máximo contraste.

### Secondary

- **Rojo Mineral:** identifica sellos, registros, alertas editoriales y pequeñas marcas de procedencia; no es un color de relleno dominante.
- **Amarillo de Levantamiento:** señala rutas, foco, selección y momentos de llegada; su rareza le da significado.

### Tertiary

- **Azul Río:** comunica foco de campo, seguridad y flujo de información sin competir con la acción primaria.

### Neutral

- **Papel Mineral:** superficie base cálida y material de continuidad.
- **Papel Claro:** superficie legible para controles y leyendas sobre cartografía compleja.
- **Rojo de Error:** reserva semántica para validación y fallos.
- **Línea de Levantamiento:** divisores y trazos secundarios que preservan la textura del papel.

### Named Rules

**The One Route Rule.** El amarillo aparece como una sola ruta o señal prioritaria por contexto; no se dispersa como acento decorativo.

**The Paper Is Material Rule.** Las superficies claras deben conservar calidez, textura y relación con el campo cartográfico; no se sustituyen por blanco clínico ni crema plano genérico.

## Typography

- **Display Font:** Archivo Variable (con Archivo y sans-serif como respaldo)
- **Body Font:** Archivo Variable (con Archivo y sans-serif como respaldo)
- **Label/Coordinate Font:** Barlow Condensed (con sans-serif como respaldo)

**Character:** Archivo aporta claridad contemporánea y admite jerarquías densas sin perder humanidad. Barlow Condensed introduce la voz del levantamiento técnico en coordenadas, escalas, sellos y microetiquetas; nunca reemplaza el cuerpo de lectura.

### Hierarchy

- **Display** (650, escala fluida, 0.98): titulares cortos de identidad y orientación, con interletraje compacto.
- **Wordmark** (700, escala fluida, 0.75): uso exclusivo y contundente para APU.
- **Body** (400, escala fluida, 1.45): explicación clara, preparada para crecer en español, inglés, quechua y aimara.
- **Label** (650, 0.82rem): nombres de campos y controles, siempre visibles.
- **Coordinate** (500, 0.68rem, 0.13em): coordenadas, escala y metadatos técnicos en mayúsculas cuando la función lo exige.

### Named Rules

**The Two Voices Rule.** Archivo habla a las personas; Barlow Condensed anota el mapa. No se usa la condensada para párrafos, errores ni instrucciones esenciales.

## Layout

El sistema utiliza campos continuos antes que colecciones de tarjetas. Las superficies de Operate deben separar información mediante marcos, reglas, ritmo y diferencias tonales, manteniendo una lectura clara aun cuando la cartografía esté presente.

En el login, la composición específica pertenece a su brief: mapa continuo, masa montañosa a la izquierda, identidad arriba a la derecha y leyenda de acceso abajo a la derecha. En otras superficies, se conserva la gramática —campo cartográfico, zona de lectura tranquila, coordenadas y puntos de registro— sin copiar esa distribución literalmente.

La densidad se reduce antes que el contenido: en pantallas estrechas desaparecen coordenadas y marcas secundarias, mientras identidad, explicación, campos, errores y acciones mantienen prioridad. Los contenedores no dependen de altura fija y los controles táctiles conservan al menos 2.75rem de área interactiva; los campos principales alcanzan 3.15rem o más.

**The Continuous Field Rule.** El fondo y la estructura espacial forman un solo territorio; no se fragmentan en un mosaico de tarjetas flotantes.

## Elevation & Depth

APU es plano por principio. La profundidad proviene de capas de papel translúcido, densidad cartográfica, reglas, desenfoque contenido y contraste tonal; no hay sombras de tarjeta. El foco visible y los cambios de color expresan interacción sin simular elevación.

### Named Rules

**The No SaaS Shadow Rule.** Las superficies de trabajo y leyendas se delimitan con trazos, muescas y material, nunca con la sombra difusa de una tarjeta genérica.

## Shapes

La forma dominante es recta y dibujada: bordes sin radio, reglas finas, esquinas recortadas y siluetas de leyenda. Los círculos se reservan para puntos de registro, destinos y marcadores de ruta; los rombos señalan actividad o verificación. Las muescas deben parecer parte de una pieza cartográfica, no un adorno repetido en todos los contenedores.

**The Registered Geometry Rule.** Cada círculo, rombo o muesca indica una coordenada, un estado o un punto de llegada; la geometría especial nunca es relleno ornamental.

## Components

### Buttons

- **Shape:** bloque rectangular de borde firme y esquinas cuadradas.
- **Primary:** campo verde bosque, texto papel claro, altura cómoda y marca de levantamiento en rombo.
- **Hover / Focus:** el hover invierte hacia amarillo de levantamiento; el foco conserva un contorno amarillo visible y el estado activo desplaza un píxel.
- **Loading / Disabled:** la etiqueta explica la acción en curso, el rombo puede rotar y la reducción de movimiento detiene la animación.
- **Secondary / Text:** acción transparente con subrayado estructural; el hover cambia a rojo mineral.

### Cards / Containers

- **Corner Style:** esquinas cuadradas o recortadas mediante muescas deliberadas.
- **Background:** papel claro translúcido cuando se superpone a cartografía; papel mineral en superficies tranquilas.
- **Shadow Strategy:** sin sombra; borde, línea interior y cambio tonal construyen separación.
- **Border:** trazo verde bosque semitransparente.
- **Internal Padding:** fluido y generoso en leyendas; compacto en sellos y metadatos.

### Inputs / Fields

- **Style:** trazo verde bosque de 2px, fondo papel claro translúcido, esquinas cuadradas y etiqueta persistente.
- **Focus:** borde azul río y halo azul tenue, además del foco global cuando corresponde.
- **Error / Disabled:** rojo de error con mensaje reservado debajo; opacidad y cursor comunican espera sin ocultar la etiqueta.

### Navigation

La navegación futura usa Archivo para destinos y Barlow Condensed para coordenadas o contexto. El estado activo debe leerse como posición registrada mediante línea, marca o campo tonal, no como una píldora genérica.

### Cartographic Frame

El marco cartográfico es la firma reusable: una regla perimetral con coordenadas, ticks, cruces y escala. Es contextual y se simplifica responsivamente; nunca debe bloquear interacción ni competir con contenido esencial.

### Login Legend

El formulario de acceso se presenta como leyenda cartográfica recortada. La ruta llega a un único marcador, la etiqueta condensada establece el contexto y los controles conservan semántica HTML, validación visible y estados accesibles.

## Do's and Don'ts

### Do:

- **Do** usar papel mineral, reglas y densidad cartográfica para construir jerarquía antes de añadir contenedores.
- **Do** reservar amarillo, rojo y azul para rutas, registros y estados con significado concreto.
- **Do** mantener Archivo en contenido humano y Barlow Condensed en anotación técnica.
- **Do** conservar errores, foco, carga, permisos y destino autorizado como información explícita.
- **Do** retirar marcas cartográficas secundarias antes que contenido cuando disminuye el viewport.

### Don't:

- **Don't** introducir tarjetas SaaS centradas, sombras ambientales, radios blandos o píldoras como lenguaje por defecto.
- **Don't** convertir los acentos minerales en una paleta multicolor decorativa.
- **Don't** recrear la cartografía principal con curvas SVG improvisadas cuando existe un raster de procedencia aprobada.
- **Don't** usar un selector de rol en el acceso compartido; los permisos y el backend determinan el destino.
- **Don't** fijar alturas que impidan expansión de texto o localización futura.
