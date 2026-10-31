import { cache } from "react";
import { cookies } from "next/headers";
import { env } from "@/config/env";
import { backendFetch } from "@/lib/api/backend";
import { AppApiError } from "@/lib/api/errors";
import { SESSION_COOKIE } from "./constants";
import { decodeToken, encodeToken, SESSION_TTL_SECONDS } from "./token";
import type { AdminSession, BackendStaffSession } from "./session.types";

export { SESSION_COOKIE };
export type { AdminSession };

export function toAdminSession(data: BackendStaffSession): AdminSession {
  return {
    staffId: data.staff.id,
    name: data.staff.name,
    mobile: data.staff.mobile,
    isSuperAdmin: data.staff.is_super_admin,
    permissions: data.permissions,
  };
}

/**
 * Resolves the signed-in staff member, reloading their permissions from the backend so that
 * deactivation and role changes apply on the next request. Cached per request.
 * Returns null when there is no valid session; throws if the backend is unreachable.
 */
export const getSession = cache(async (): Promise<AdminSession | null> => {
  const store = await cookies();
  const token = decodeToken(store.get(SESSION_COOKIE)?.value, env.sessionSecret);
  if (!token) {
    return null;
  }
  try {
    const data = await backendFetch<BackendStaffSession>({
      path: `/api/v1/panel/auth/session/${token.sid}`,
    });
    return toAdminSession(data);
  } catch (error) {
    if (error instanceof AppApiError && error.status === 401) {
      return null;
    }
    throw error;
  }
});

export async function setSessionCookie(staffId: number): Promise<void> {
  const store = await cookies();
  const exp = Math.floor(Date.now() / 1000) + SESSION_TTL_SECONDS;
  store.set(SESSION_COOKIE, encodeToken({ sid: staffId, exp }, env.sessionSecret), {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
    secure: process.env.NODE_ENV === "production",
    maxAge: SESSION_TTL_SECONDS,
  });
}

export async function clearSessionCookie(): Promise<void> {
  const store = await cookies();
  store.delete(SESSION_COOKIE);
}
