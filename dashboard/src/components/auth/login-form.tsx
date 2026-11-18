"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, ShieldCheck } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormField } from "@/components/ui/form-field";
import { SubmitButton } from "@/components/ui/submit-button";
import { APIError, login } from "@/lib/api";

const loginSchema = z.object({
  email: z.string().trim().min(1, "Ingresa tu correo electrónico.").email("Escribe un correo electrónico válido."),
  password: z.string().min(1, "Ingresa tu contraseña."),
});

type LoginValues = z.infer<typeof loginSchema>;
type Feedback = { kind: "info" | "error" | "success"; text: string } | null;

const destinationLabels: Record<string, string> = {
  "/platform": "administración de APU",
  "/admin": "administración del cliente",
  "/familia": "panel familiar",
  "/acceso-pendiente": "acceso pendiente",
};

export function LoginForm() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [feedback, setFeedback] = useState<Feedback>(null);
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    mode: "onBlur",
    defaultValues: { email: "", password: "" },
  });

  async function onSubmit(values: LoginValues) {
    setFeedback(null);
    try {
      const result = await login(values.email, values.password);
      setFeedback({
        kind: "success",
        text: `Sesión iniciada para ${result.user.displayName}. Destino autorizado: ${destinationLabels[result.redirect] ?? "panel asignado"}.`,
      });
      router.replace(result.redirect);
      router.refresh();
    } catch (error) {
      setFeedback({
        kind: "error",
        text: error instanceof APIError ? error.message : "No pudimos conectar con el servicio de acceso.",
      });
    }
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
