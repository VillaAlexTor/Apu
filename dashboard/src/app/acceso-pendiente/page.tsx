import { Suspense } from "react";
import { ProtectedLanding } from "@/components/session/protected-landing";

export default function PendingAccessPage() {
  return <Suspense fallback={null}><ProtectedLanding destination="/acceso-pendiente" title="Acceso pendiente" description="La cuenta es válida, pero todavía no tiene una organización o rol asignado. Comunícate con un administrador de APU." /></Suspense>;
}
