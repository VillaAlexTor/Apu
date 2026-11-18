import { Suspense } from "react";
import { ProtectedLanding } from "@/components/session/protected-landing";

export default function FamilyPage() {
  return <Suspense fallback={null}><ProtectedLanding destination="/familia" title="Panel familiar" description="Tu sesión está protegida y el rol de tutor fue validado. El emparejamiento de dispositivos se construirá en el siguiente módulo." /></Suspense>;
}
