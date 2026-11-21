"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Clipboard, KeyRound, ShieldCheck } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { APIError, confirmMFAEnrollment, startMFAEnrollment, verifyMFA, type LoginSuccess } from "@/lib/api";

type PendingChallenge = { challenge: string; mode: "enroll" | "verify" };
type Enrollment = { secret: string; provisioningUri: string };
const schema = z.object({ code: z.string().trim().min(6, "Ingresa el código de tu aplicación o un código de recuperación.").max(32) });
type Values = z.infer<typeof schema>;

export function MFAForm() {
  const router = useRouter();
  const [pending, setPending] = useState<PendingChallenge | null>(null);
  const [enrollment, setEnrollment] = useState<Enrollment | null>(null);
  const [recoveryCodes, setRecoveryCodes] = useState<string[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [destination, setDestination] = useState<LoginSuccess["redirect"]>("/acceso-pendiente");
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Values>({ resolver: zodResolver(schema) });

  useEffect(() => {
    let active = true;
    async function initialize() {
      await Promise.resolve();
      const raw = sessionStorage.getItem("apu_mfa_challenge");
      if (!raw) {
        if (active) setLoading(false);
        return;
      }
      try {
        const parsed = JSON.parse(raw) as PendingChallenge;
        if (!parsed.challenge || !["enroll", "verify"].includes(parsed.mode)) throw new Error("invalid challenge");
        if (active) setPending(parsed);
        if (parsed.mode === "enroll") {
          const result = await startMFAEnrollment(parsed.challenge);
          if (active) setEnrollment(result);
        }
      } catch (requestError) {
        sessionStorage.removeItem("apu_mfa_challenge");
        if (active) {
          setPending(null);
          setEnrollment(null);
          setError(requestError instanceof APIError ? requestError.message : "No pudimos preparar el segundo factor.");
        }
      } finally {
        if (active) setLoading(false);
      }
    }
    void initialize();
    return () => { active = false; };
  }, []);

  function finish(result: LoginSuccess) {
    sessionStorage.removeItem("apu_mfa_challenge");
    if (result.recoveryCodes?.length) {
      setDestination(result.redirect);
      setRecoveryCodes(result.recoveryCodes);
      return;
    }
    router.replace(result.redirect);
    router.refresh();
  }

  async function onSubmit(values: Values) {
    if (!pending) return;
    setError(null);
    try {
      const result = pending.mode === "enroll" ? await confirmMFAEnrollment(pending.challenge, values.code) : await verifyMFA(pending.challenge, values.code);
      finish(result);
    } catch (verifyError) {
      if (verifyError instanceof APIError && verifyError.code === "mfa_challenge_invalid") {
        sessionStorage.removeItem("apu_mfa_challenge");
        setPending(null);
        setEnrollment(null);
      }
      setError(verifyError instanceof APIError ? verifyError.message : "No pudimos verificar el código.");
    }
  }

  async function copyRecoveryCodes() {
    if (!recoveryCodes) return;
    try {
      await navigator.clipboard.writeText(recoveryCodes.join("\n"));
      setCopied(true);
    } catch {
      setError("No pudimos copiar los códigos. Guárdalos manualmente antes de continuar.");
    }
  }

  if (loading) return <div className="security-loading" aria-label="Preparando verificación"><span /><span /><span /></div>;
  if (recoveryCodes) return <div className="security-complete recovery-codes"><ShieldCheck aria-hidden="true" /><h1>Guarda tus códigos de recuperación</h1><p>Cada código funciona una sola vez. Guárdalos fuera de este dispositivo; no volverán a mostrarse.</p><ul>{recoveryCodes.map((code) => <li key={code}><code>{code}</code></li>)}</ul><button type="button" onClick={copyRecoveryCodes}><Clipboard aria-hidden="true" />{copied ? "Códigos copiados" : "Copiar códigos"}</button><button type="button" className="primary-security-action" onClick={() => { router.replace(destination); router.refresh(); }}>Continuar al panel</button></div>;
  if (!pending) return <div className="security-complete"><KeyRound aria-hidden="true" /><h1>La verificación expiró</h1><p>Inicia sesión nuevamente para obtener un reto de seguridad.</p><Link href="/">Volver al acceso</Link></div>;

  return <form className="security-form mfa-form" onSubmit={handleSubmit(onSubmit)} noValidate><ShieldCheck aria-hidden="true" /><h1>{pending.mode === "enroll" ? "Protege tu acceso administrativo" : "Verifica tu segundo factor"}</h1><p>{pending.mode === "enroll" ? "Agrega APU a una aplicación autenticadora y confirma el código de seis dígitos." : "Ingresa el código vigente de tu aplicación o uno de recuperación."}</p>{pending.mode === "enroll" && enrollment && <div className="mfa-setup"><a href={enrollment.provisioningUri}>Abrir aplicación autenticadora</a><span>O ingresa esta clave manualmente:</span><code>{enrollment.secret}</code></div>}<label>Código de verificación<input type="text" inputMode={pending.mode === "enroll" ? "numeric" : "text"} autoComplete="one-time-code" spellCheck={false} autoCapitalize="characters" disabled={isSubmitting || (pending.mode === "enroll" && !enrollment)} {...register("code")} />{errors.code && <small>{errors.code.message}</small>}</label>{error && <div className="activation-error" role="alert">{error}</div>}<button type="submit" disabled={isSubmitting || (pending.mode === "enroll" && !enrollment)}>{isSubmitting ? "Verificando…" : pending.mode === "enroll" ? "Activar segundo factor" : "Verificar y continuar"}</button><Link className="text-action" href="/">Cancelar</Link></form>;
}
