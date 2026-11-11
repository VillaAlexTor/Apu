import { forwardRef, type InputHTMLAttributes, type ReactNode } from "react";

type FormFieldProps = InputHTMLAttributes<HTMLInputElement> & {
  id: string;
  label: string;
  error?: string;
  action?: ReactNode;
};

export const FormField = forwardRef<HTMLInputElement, FormFieldProps>(function FormField(
  { id, label, error, action, className, ...props },
  ref,
) {
  const errorId = `${id}-error`;
  return (
    <div className={`form-field${error ? " has-error" : ""}`}>
      <label htmlFor={id}>{label}</label>
      <div className="field-control">
        <input ref={ref} id={id} className={className} aria-invalid={Boolean(error)} aria-describedby={error ? errorId : undefined} {...props} />
        {action}
      </div>
      <p id={errorId} className="field-error" aria-live="polite">{error ?? ""}</p>
    </div>
  );
});
