"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { MailCheck } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { APIError, requestPasswordReset } from "@/lib/api";

const schema = z.object({ email: z.string().trim().email("Escribe un correo electrónico válido.") });
type Values = z.infer<typeof schema>;

export function RecoveryForm() {
  const [sent, setSent] = useState(false);
  const [developmentToken, setDevelopmentToken] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Values>({ resolver: zodResolver(schema) });

  async function onSubmit(values: Values) {
    setError(null);
    try {
      const result = await requestPasswordReset(values.email);
      setDevelopmentToken(result.developmentResetToken ?? null);
      setSent(true);
    } catch (requestError) {
      setError(requestError instanceof APIError ? requestError.message : "No pudimos enviar la solicitud. Intenta nuevamente.");
    }
  }

  if (sent) {
    return <div className="security-complete"><MailCheck aria-hidden="true" /><h1>Revisa tu correo</h1><p>Si existe una cuenta activa para ese correo, recibirás un enlace válido durante 30 minutos.</p>{developmentToken && <Link href={`/restablecer?token=${encodeURIComponent(developmentToken)}`}>Abrir enlace local de prueba</Link>}<Link className="text-action" href="/">Volver al acceso</Link></div>;
  }

  return <form className="security-form" onSubmit={handleSubmit(onSubmit)} noValidate><h1>Recupera tu acceso</h1><p>Escribe tu correo. La respuesta será la misma exista o no una cuenta.</p><label>Correo electrónico<input type="email" inputMode="email" autoComplete="email" disabled={isSubmitting} {...register("email")} />{errors.email && <small>{errors.email.message}</small>}</label>{error && <div className="activation-error" role="alert">{error}</div>}<button type="submit" disabled={isSubmitting}>{isSubmitting ? "Enviando…" : "Enviar enlace seguro"}</button><Link className="text-action" href="/">Volver al acceso</Link></form>;
}
