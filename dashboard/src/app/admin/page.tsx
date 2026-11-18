import { Suspense } from "react";
import { ProtectedLanding } from "@/components/session/protected-landing";

export default function AdminPage() {
  return <Suspense fallback={null}><ProtectedLanding destination="/admin" title="Administración del cliente" description="Tu sesión está protegida y el rol fue validado. La gestión de tutores, familias y dispositivos será el siguiente módulo." /></Suspense>;
}
