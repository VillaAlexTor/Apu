import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { Suspense } from "react";
import { ChangePasswordForm } from "@/components/auth/change-password-form";
import { getServerSession } from "@/lib/server-api";

export const metadata: Metadata = { title: "Seguridad de la cuenta | APU" };

export default function SecurityPage() {
  return <Suspense fallback={<main className="security-page" />}><ProtectedSecurity /></Suspense>;
}

async function ProtectedSecurity() {
  const session = await getServerSession();
  if (!session) redirect("/");
  return <main className="security-page"><section className="security-sheet"><div className="security-wordmark">APU<span /></div><ChangePasswordForm returnTo={session.redirect} /></section></main>;
}
