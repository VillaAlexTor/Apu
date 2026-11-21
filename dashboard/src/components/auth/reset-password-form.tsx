"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { KeyRound } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { APIError, confirmPasswordReset } from "@/lib/api";

const schema = z.object({
  password: z.string().min(15, "Usa al menos 15 caracteres.").max(1024, "La contraseña es demasiado larga."),
  confirmation: z.string(),
}).refine((values) => values.password === values.confirmation, { path: ["confirmation"], message: "Las contraseñas no coinciden." });
type Values = z.infer<typeof schema>;

export function ResetPasswordForm({ token }: { token: string }) {
  const [complete, setComplete] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Values>({ resolver: zodResolver(schema) });

  async function onSubmit(values: Values) {
    setError(null);
    try {
      await confirmPasswordReset(token, values.password);
      setComplete(true);
    } catch (resetError) {
      setError(resetError instanceof APIError ? resetError.message : "No pudimos restablecer la contraseña.");
    }
  }

  if (complete) return <div className="security-complete"><KeyRound aria-hidden="true" /><h1>Contraseña actualizada</h1><p>Todas las sesiones anteriores fueron revocadas. Inicia sesión nuevamente.</p><Link href="/">Volver al acceso</Link></div>;
  return <form className="security-form" onSubmit={handleSubmit(onSubmit)} noValidate><h1>Crea una nueva contraseña</h1><p>Usa una frase única de al menos 15 caracteres.</p><label>Nueva contraseña<input type="password" autoComplete="new-password" disabled={isSubmitting} {...register("password")} />{errors.password && <small>{errors.password.message}</small>}</label><label>Confirma la contraseña<input type="password" autoComplete="new-password" disabled={isSubmitting} {...register("confirmation")} />{errors.confirmation && <small>{errors.confirmation.message}</small>}</label>{error && <div className="activation-error" role="alert">{error}</div>}<button type="submit" disabled={isSubmitting}>{isSubmitting ? "Actualizando…" : "Restablecer contraseña"}</button></form>;
}
