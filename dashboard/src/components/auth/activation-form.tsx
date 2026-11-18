"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { CheckCircle2, KeyRound } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { activateInvitation, APIError } from "@/lib/api";

const schema = z.object({
  password: z.string().min(12, "Usa al menos 12 caracteres.").max(1024, "La contraseña es demasiado larga."),
  confirmation: z.string(),
}).refine((values) => values.password === values.confirmation, { path: ["confirmation"], message: "Las contraseñas no coinciden." });

type Values = z.infer<typeof schema>;

export function ActivationForm({ token }: { token: string }) {
  const [complete, setComplete] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Values>({ resolver: zodResolver(schema) });

  async function onSubmit(values: Values) {
    setServerError(null);
    try {
      await activateInvitation(token, values.password);
      setComplete(true);
    } catch (error) {
      setServerError(error instanceof APIError ? error.message : "No pudimos activar la invitación.");
    }
  }

  if (complete) return <div className="activation-complete"><CheckCircle2 aria-hidden="true" /><h1>Acceso activado</h1><p>Tu contraseña quedó registrada. Ya puedes entrar al panel asignado.</p><Link href="/">Ir al inicio de sesión</Link></div>;

  return <form className="activation-form" onSubmit={handleSubmit(onSubmit)} noValidate><KeyRound aria-hidden="true" /><h1>Activa tu acceso</h1><p>Crea una contraseña para aceptar la invitación a APU.</p><label>Nueva contraseña<input type="password" autoComplete="new-password" {...register("password")} />{errors.password && <small>{errors.password.message}</small>}</label><label>Confirma la contraseña<input type="password" autoComplete="new-password" {...register("confirmation")} />{errors.confirmation && <small>{errors.confirmation.message}</small>}</label>{serverError && <div className="activation-error" role="alert">{serverError}</div>}<button type="submit" disabled={isSubmitting}>{isSubmitting ? "Activando…" : "Activar acceso"}</button></form>;
}
