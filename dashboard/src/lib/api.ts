const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8081/api/v1";

type APIErrorBody = {
  error?: {
    code?: string;
    message?: string;
  };
};

export type AuthenticatedUser = {
  id: string;
  email: string;
  displayName: string;
  platformRole?: "platform_admin" | "support_agent";
  memberships: Array<{
    tenantId: string;
    tenantName: string;
    role: "tenant_owner" | "tenant_admin" | "guardian";
  }>;
};

export type LoginResult = {
  user: AuthenticatedUser;
  redirect: "/platform" | "/admin" | "/familia" | "/acceso-pendiente";
  csrfToken: string;
};

export class APIError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly code: string,
  ) {
    super(message);
  }
}

export async function login(email: string, password: string): Promise<LoginResult> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    const body = (await response.json().catch(() => ({}))) as APIErrorBody;
    throw new APIError(
      body.error?.message ?? "No pudimos iniciar sesión.",
      response.status,
      body.error?.code ?? "unknown_error",
    );
  }

  return (await response.json()) as LoginResult;
}
