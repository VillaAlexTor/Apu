import type { Metadata } from "next";
import { MFAForm } from "@/components/auth/mfa-form";

export const metadata: Metadata = { title: "Verificación en dos pasos | APU", referrer: "no-referrer" };

export default function MFAPage() {
  return <main className="security-page"><section className="security-sheet"><div className="security-wordmark">APU<span /></div><MFAForm /></section></main>;
}
