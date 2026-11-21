"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { ShieldCheck } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { APIError, changePassword } from "@/lib/api";

const schema = z.object({
  currentPassword: z.string().min(1, "Ingresa tu contraseña actual."),
  newPassword: z.string().min(15, "Usa al menos 15 caracteres.").max(1024),
  confirmation: z.string(),
}).refine((values) => values.newPassword === values.confirmation, { path: ["confirmation"], message: "Las contraseñas no coinciden." });
type Values = z.infer<typeof schema>;

export function ChangePasswordForm({ returnTo }: { returnTo: string }) {
  const [complete, setComplete] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Values>({ resolver: zodResolver(schema) });
  async function onSubmit(values: Values) {
    setError(null);
    try {
      await changePassword(values.currentPassword, values.newPassword);
      setComplete(true);
    } catch (changeError) {
      setError(changeError instanceof APIError ? changeError.message : "No pudimos cambiar la contraseña.");
    }
  }
  if (complete) return <div className="security-complete"><ShieldCheck aria-hidden="true" /><h1>Acceso renovado</h1><p>La contraseña cambió y todas tus sesiones fueron cerradas.</p><Link href="/">Iniciar sesión nuevamente</Link></div>;
  return <form className="security-form" onSubmit={handleSubmit(onSubmit)} noValidate><h1>Cambia tu contraseña</h1><p>Confirmaremos tu contraseña actual y cerraremos todas las sesiones al terminar.</p><label>Contraseña actual<input type="password" autoComplete="current-password" disabled={isSubmitting} {...register("currentPassword")} />{errors.currentPassword && <small>{errors.currentPassword.message}</small>}</label><label>Nueva contraseña<input type="password" autoComplete="new-password" disabled={isSubmitting} {...register("newPassword")} />{errors.newPassword && <small>{errors.newPassword.message}</small>}</label><label>Confirma la nueva contraseña<input type="password" autoComplete="new-password" disabled={isSubmitting} {...register("confirmation")} />{errors.confirmation && <small>{errors.confirmation.message}</small>}</label>{error && <div className="activation-error" role="alert">{error}</div>}<button type="submit" disabled={isSubmitting}>{isSubmitting ? "Protegiendo…" : "Cambiar contraseña"}</button><Link className="text-action" href={returnTo}>Cancelar y volver</Link></form>;
}
