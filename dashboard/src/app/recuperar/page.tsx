import type { Metadata } from "next";
import { RecoveryForm } from "@/components/auth/recovery-form";

export const metadata: Metadata = { title: "Recuperar acceso | APU" };

export default function RecoveryPage() {
  return <main className="security-page"><section className="security-sheet"><div className="security-wordmark">APU<span /></div><RecoveryForm /></section></main>;
}
