import type { Metadata } from "next";
import Link from "next/link";
import { Suspense } from "react";
import { ResetPasswordForm } from "@/components/auth/reset-password-form";

export const metadata: Metadata = { title: "Restablecer contraseña | APU", referrer: "no-referrer" };

export default function ResetPasswordPage({ searchParams }: { searchParams: Promise<{ token?: string }> }) {
  return <Suspense fallback={<main className="security-page" />}><ResetRoute searchParams={searchParams} /></Suspense>;
}

async function ResetRoute({ searchParams }: { searchParams: Promise<{ token?: string }> }) {
  const token = (await searchParams).token?.trim() ?? "";
  return <main className="security-page"><section className="security-sheet"><div className="security-wordmark">APU<span /></div>{token ? <ResetPasswordForm token={token} /> : <div className="security-complete"><h1>Enlace incompleto</h1><p>Solicita un nuevo enlace de recuperación.</p><Link href="/recuperar">Solicitar otro enlace</Link></div>}</section></main>;
}
