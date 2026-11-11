"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, ShieldCheck } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormField } from "@/components/ui/form-field";
import { SubmitButton } from "@/components/ui/submit-button";

const loginSchema = z.object({
  email: z.string().trim().min(1, "Ingresa tu correo electrónico.").email("Escribe un correo electrónico válido."),
  password: z.string().min(1, "Ingresa tu contraseña.").min(8, "La contraseña debe tener al menos 8 caracteres."),
});

type LoginValues = z.infer<typeof loginSchema>;
type Feedback = { kind: "info" | "error"; text: string } | null;

export function LoginForm() {
  const [showPassword, setShowPassword] = useState(false);
  const [feedback, setFeedback] = useState<Feedback>(null);
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    mode: "onBlur",
    defaultValues: { email: "", password: "" },
  });

  async function onSubmit() {
    setFeedback(null);
    await new Promise((resolve) => window.setTimeout(resolve, 700));
    setFeedback({ kind: "info", text: "El formulario está listo. La sesión se activará cuando conectemos la API de Go." });
  }

  function showRecoveryStatus() {
    setFeedback({ kind: "info", text: "La recuperación de contraseña se habilitará junto con el servicio de correo del backend." });
  }

  return (
    <form className="login-form" onSubmit={handleSubmit(onSubmit)} noValidate>
      <FormField id="email" label="Correo electrónico" type="email" autoComplete="username" inputMode="email" placeholder="nombre@ejemplo.com" error={errors.email?.message} disabled={isSubmitting} {...register("email")} />
      <FormField
        id="password"
        label="Contraseña"
        type={showPassword ? "text" : "password"}
        autoComplete="current-password"
        placeholder="Ingresa tu contraseña"
        error={errors.password?.message}
        disabled={isSubmitting}
        action={
          <button className="field-action" type="button" onClick={() => setShowPassword((visible) => !visible)} aria-label={showPassword ? "Ocultar contraseña" : "Mostrar contraseña"} aria-pressed={showPassword} disabled={isSubmitting}>
            {showPassword ? <EyeOff aria-hidden="true" /> : <Eye aria-hidden="true" />}
          </button>
        }
        {...register("password")}
      />
      <SubmitButton loading={isSubmitting}>Iniciar sesión</SubmitButton>
      <button className="recovery-link" type="button" onClick={showRecoveryStatus} disabled={isSubmitting}>¿Olvidaste tu contraseña?</button>
      <div className="security-note">
        <ShieldCheck aria-hidden="true" />
        <p>Tu acceso se dirige de forma segura según tus permisos.</p>
      </div>
      <div className={`form-feedback${feedback ? ` is-${feedback.kind}` : ""}`} role="status" aria-live="polite">{feedback?.text}</div>
    </form>
  );
}
