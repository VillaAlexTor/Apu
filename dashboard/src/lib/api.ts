const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8081/api/v1";

type APIErrorBody = { error?: { code?: string; message?: string } };

export type Membership = {
  tenantId: string;
  tenantName: string;
  role: "tenant_owner" | "tenant_admin" | "guardian";
};

export type AuthenticatedUser = {
  id: string;
  email: string;
  displayName: string;
  platformRole?: "platform_admin" | "support_agent";
  memberships: Membership[];
};

export type SessionResult = {
  user: AuthenticatedUser;
  redirect: "/platform" | "/admin" | "/familia" | "/acceso-pendiente";
};

export type LoginSuccess = SessionResult & { csrfToken: string; recoveryCodes?: string[] };
export type MFAChallengeResult = { mfaRequired: true; mode: "enroll" | "verify"; challenge: string };
export type LoginResult = LoginSuccess | MFAChallengeResult;

export type PlatformSummary = {
  tenantsTotal: number;
  tenantsActive: number;
  usersActive: number;
  usersInvited: number;
  usersSuspended: number;
  pendingInvitations: number;
};

export type TenantRecord = {
  id: string;
  slug: string;
  name: string;
  status: "active" | "suspended" | "deactivated";
  ownerName?: string;
  ownerEmail?: string;
  userCount: number;
  activeSessions: number;
  createdAt: string;
};

export type PlatformUser = {
  id: string;
  email: string;
  displayName: string;
  status: "invited" | "active" | "suspended" | "deactivated";
  platformRole?: "platform_admin" | "support_agent";
  memberships: Membership[];
  activeSessions: number;
  createdAt: string;
};

export type AuditEvent = {
  id: number;
  actorEmail?: string;
  tenantName?: string;
  action: string;
  resourceType: string;
  resourceId?: string;
  metadata: Record<string, unknown>;
  occurredAt: string;
};

export type PlatformSnapshot = {
  summary: PlatformSummary;
  tenants: TenantRecord[];
  users: PlatformUser[];
  audit: AuditEvent[];
};

export class APIError extends Error {
  constructor(message: string, readonly status: number, readonly code: string) {
    super(message);
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      ...(init?.body ? { "Content-Type": "application/json" } : {}),
      ...init?.headers,
    },
  });
  if (!response.ok) {
    const body = (await response.json().catch(() => ({}))) as APIErrorBody;
    throw new APIError(
      body.error?.message ?? "No pudimos completar la solicitud.",
      response.status,
      body.error?.code ?? "unknown_error",
    );
  }
  if (response.status === 204) return undefined as T;
  return (await response.json()) as T;
}

function csrfToken(): string {
  if (typeof document === "undefined") return "";
  const cookie = document.cookie
    .split("; ")
    .find((item) => item.startsWith("apu_csrf=") || item.startsWith("__Host-apu_csrf="));
  return cookie ? decodeURIComponent(cookie.slice(cookie.indexOf("=") + 1)) : "";
}

function mutation(method: "POST" | "PATCH", body?: unknown): RequestInit {
  return {
    method,
    headers: { "X-CSRF-Token": csrfToken() },
    ...(body === undefined ? {} : { body: JSON.stringify(body) }),
  };
}

export function login(email: string, password: string): Promise<LoginResult> {
  return request("/auth/login", { method: "POST", body: JSON.stringify({ email, password }) });
}

export function startMFAEnrollment(challenge: string): Promise<{ secret: string; provisioningUri: string }> {
  return request("/auth/mfa/enroll", { method: "POST", body: JSON.stringify({ challenge }) });
}

export function confirmMFAEnrollment(challenge: string, code: string): Promise<LoginSuccess> {
  return request("/auth/mfa/enroll/confirm", { method: "POST", body: JSON.stringify({ challenge, code }) });
}

export function verifyMFA(challenge: string, code: string): Promise<LoginSuccess> {
  return request("/auth/mfa/verify", { method: "POST", body: JSON.stringify({ challenge, code }) });
}

export function requestPasswordReset(email: string): Promise<{ status: string; developmentResetToken?: string }> {
  return request("/auth/password/reset/request", { method: "POST", body: JSON.stringify({ email }) });
}

export function confirmPasswordReset(token: string, newPassword: string): Promise<{ status: string }> {
  return request("/auth/password/reset/confirm", { method: "POST", body: JSON.stringify({ token, newPassword }) });
}

export function changePassword(currentPassword: string, newPassword: string): Promise<void> {
  return request("/auth/password/change", mutation("POST", { currentPassword, newPassword }));
}

export function activateInvitation(token: string, password: string): Promise<{ status: string }> {
  return request("/auth/activate", { method: "POST", body: JSON.stringify({ token, password }) });
}

export function logout(): Promise<void> {
  return request("/auth/logout", mutation("POST"));
}

export async function getPlatformSnapshot(): Promise<PlatformSnapshot> {
  const [summary, tenantResponse, userResponse, auditResponse] = await Promise.all([
    request<PlatformSummary>("/platform/summary"),
    request<{ items: TenantRecord[] }>("/platform/tenants"),
    request<{ items: PlatformUser[] }>("/platform/users"),
    request<{ items: AuditEvent[] }>("/platform/audit"),
  ]);
  return { summary, tenants: tenantResponse.items, users: userResponse.items, audit: auditResponse.items };
}

export function createTenant(input: { name: string; slug: string }): Promise<TenantRecord> {
  return request("/platform/tenants", mutation("POST", input));
}

export function setTenantStatus(tenantId: string, status: TenantRecord["status"]): Promise<{ status: string }> {
  return request(`/platform/tenants/${tenantId}/status`, mutation("PATCH", { status }));
}

export function createInvitation(input: { tenantId: string; email: string; displayName: string; role: Membership["role"] }): Promise<{ user: PlatformUser; requiresActivation: boolean; activationToken?: string }> {
  return request("/platform/invitations", mutation("POST", input));
}

export function setUserStatus(userId: string, status: "active" | "suspended" | "deactivated"): Promise<{ status: string }> {
  return request(`/platform/users/${userId}/status`, mutation("PATCH", { status }));
}

export function revokeUserSessions(userId: string): Promise<{ revokedSessions: number }> {
  return request(`/platform/users/${userId}/sessions/revoke`, mutation("POST"));
}

export function setMembershipRole(tenantId: string, userId: string, role: Membership["role"]): Promise<{ role: string }> {
  return request(`/platform/tenants/${tenantId}/users/${userId}`, mutation("PATCH", { role }));
}
