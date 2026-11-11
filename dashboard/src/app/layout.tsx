import type { Metadata } from "next";
import "@fontsource-variable/archivo";
import "@fontsource/barlow-condensed/500.css";
import "@fontsource/barlow-condensed/700.css";
import "./globals.css";

const directionContract = `<!--
THESIS: APU turns sign-in into a transparent cartographic checkpoint and refuses the generic centered SaaS card.
OWN-WORLD: Mineral paper, forest-green survey rules, dense topographic contours, river blue, mineral red and a single surveying-yellow route.
STORY: The visitor recognizes APU, enters one shared access, understands that permissions determine the destination, and receives explicit state feedback.
FIRST VIEWPORT: A continuous map fills the screen; mountain mass owns the left, APU sits upper-right, and a notched cartographic legend anchors the form lower-right where the yellow route arrives.
FORM: Atlas topográfico andino, selected comp atlas-inset, seed login-atlas-inset.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
-->`;

export const metadata: Metadata = {
  title: "Iniciar sesión | APU",
  description: "Acceso seguro al panel de protección familiar y operación de APU.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="es">
      <body>
        <template data-impeccable-contract="login-atlas-inset" dangerouslySetInnerHTML={{ __html: directionContract }} />
        {children}
      </body>
    </html>
  );
}
