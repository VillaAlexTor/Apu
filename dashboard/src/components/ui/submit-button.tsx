import type { ButtonHTMLAttributes, ReactNode } from "react";

type SubmitButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  loading?: boolean;
  children: ReactNode;
};

export function SubmitButton({ loading = false, children, disabled, ...props }: SubmitButtonProps) {
  return (
    <button className="submit-button" type="submit" disabled={disabled || loading} aria-busy={loading} {...props}>
      <span className="button-label">{loading ? "Verificando acceso…" : children}</span>
      <span className="button-survey-mark" aria-hidden="true" />
    </button>
  );
}
