"use client";

import { LogOut, ShieldCheck } from "lucide-react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useState } from "react";
import { APIError, logout, type AuthenticatedUser } from "@/lib/api";

export function SessionLanding({ user, title, description }: { user: AuthenticatedUser; title: string; description: string }) {
  const router = useRouter();
  const [closing, setClosing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  async function closeSession() {
    setClosing(true);
    setError(null);
    try {
      await logout();
      router.replace("/");
      router.refresh();
    } catch (logoutError) {
      setError(logoutError instanceof APIError ? logoutError.message : "No se pudo cerrar la sesión. Intenta nuevamente.");
      setClosing(false);
    }
  }
  return <main className="session-landing"><section><div className="session-wordmark">APU<span /></div><ShieldCheck className="session-icon" aria-hidden="true" /><h1>{title}</h1><p>{description}</p><dl><div><dt>Usuario</dt><dd>{user.displayName}</dd></div><div><dt>Correo</dt><dd>{user.email}</dd></div></dl>{error && <div className="activation-error" role="alert">{error}</div>}<div className="session-actions"><Link href="/seguridad">Cambiar contraseña</Link><button type="button" onClick={closeSession} disabled={closing}><LogOut aria-hidden="true" />{closing ? "Cerrando…" : "Cerrar sesión"}</button></div></section></main>;
}
