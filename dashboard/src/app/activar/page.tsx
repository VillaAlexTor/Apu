import type { Metadata } from "next";
import Link from "next/link";
import { Suspense } from "react";
import { ActivationForm } from "@/components/auth/activation-form";

export const metadata: Metadata = { title: "Activar acceso | APU" };

export default function ActivatePage({ searchParams }: { searchParams: Promise<{ token?: string }> }) {
  return <Suspense fallback={<main className="activation-page" />}><ActivationRoute searchParams={searchParams} /></Suspense>;
}

async function ActivationRoute({ searchParams }: { searchParams: Promise<{ token?: string }> }) {
  const token = (await searchParams).token?.trim() ?? "";
  return <main className="activation-page"><section className="activation-sheet"><div className="activation-brand">APU<span /></div>{token ? <ActivationForm token={token} /> : <div className="activation-complete"><h1>Invitación incompleta</h1><p>El enlace no contiene el código de activación. Solicita una nueva invitación al administrador.</p><Link href="/">Volver al acceso</Link></div>}</section></main>;
}
