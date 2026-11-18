import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { Suspense } from "react";
import { PlatformDashboard } from "@/components/platform/platform-dashboard";
import { getServerPlatformSnapshot, getServerSession } from "@/lib/server-api";
import "./platform.css";

export const metadata: Metadata = {
  title: "Operación de plataforma | APU",
  description: "Administración de organizaciones, usuarios y auditoría de APU.",
};

export default function PlatformPage() {
  return <Suspense fallback={<PlatformLoading />}><ProtectedPlatform /></Suspense>;
}

async function ProtectedPlatform() {
  const session = await getServerSession();
  if (!session) redirect("/");
  if (session.redirect !== "/platform" || session.user.platformRole !== "platform_admin") redirect(session.redirect);
  const snapshot = await getServerPlatformSnapshot();
  if (!snapshot) return <PlatformUnavailable />;
  return <PlatformDashboard user={session.user} initialSnapshot={snapshot} />;
}

function PlatformUnavailable() {
  return (
    <main className="platform-unavailable">
      <section>
        <h1>Operación no disponible</h1>
        <p>No pudimos cargar la plataforma. Tu sesión sigue protegida; reintenta cuando el servicio de administración esté disponible.</p>
        <a href="/platform">Reintentar</a>
      </section>
    </main>
  );
}

function PlatformLoading() {
  return <main className="platform-loading" aria-label="Cargando operación de plataforma"><div className="loading-rail" /><div className="loading-field"><span /><span /><span /><span /></div></main>;
}
