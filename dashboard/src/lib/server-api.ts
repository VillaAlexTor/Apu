import "server-only";

import { cookies } from "next/headers";
import type { PlatformSnapshot, SessionResult } from "@/lib/api";

const INTERNAL_API_BASE_URL = process.env.APU_INTERNAL_API_URL ?? process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8081/api/v1";

async function serverRequest<T>(path: string): Promise<{ status: number; data?: T }> {
  const cookieHeader = (await cookies()).toString();
  const response = await fetch(`${INTERNAL_API_BASE_URL}${path}`, {
    cache: "no-store",
    headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
  });
  if (!response.ok) return { status: response.status };
  return { status: response.status, data: (await response.json()) as T };
}

export async function getServerSession(): Promise<SessionResult | null> {
  const response = await serverRequest<SessionResult>("/auth/me");
  return response.status === 200 && response.data ? response.data : null;
}

export async function getServerPlatformSnapshot(): Promise<PlatformSnapshot | null> {
  const [summary, tenants, users, audit] = await Promise.all([
    serverRequest<PlatformSnapshot["summary"]>("/platform/summary"),
    serverRequest<{ items: PlatformSnapshot["tenants"] }>("/platform/tenants"),
    serverRequest<{ items: PlatformSnapshot["users"] }>("/platform/users"),
    serverRequest<{ items: PlatformSnapshot["audit"] }>("/platform/audit"),
  ]);
  if (!summary.data || !tenants.data || !users.data || !audit.data) return null;
  return { summary: summary.data, tenants: tenants.data.items, users: users.data.items, audit: audit.data.items };
}
