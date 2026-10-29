import { createHmac, timingSafeEqual } from "node:crypto";
import { cookies } from "next/headers";
import { env } from "@/config/env";
import { toLatinDigits } from "@/lib/formatters/digits";
import { ALL_PERMISSIONS } from "@/lib/permissions/codes";
import { SESSION_COOKIE } from "./constants";
import type { AdminSession } from "./session.types";

export { SESSION_COOKIE };
export type { AdminSession };

const DEMO_SESSION: AdminSession = {
  staffId: 1,
  name: "مدیر روآب",
  mobile: "09120000000",
  permissions: ALL_PERMISSIONS,
};

export function isDemoCredentials(mobile: string, password: string): boolean {
  const normalized = toLatinDigits(mobile).replace(/\D/g, "");
  return (
    (normalized === "09120000000" || normalized === "9120000000") &&
    password === "admin"
  );
}

export function demoSession(): AdminSession {
  return DEMO_SESSION;
}

function sign(payload: string): string {
  return createHmac("sha256", env.sessionSecret).update(payload).digest("base64url");
}

export function encodeSession(session: AdminSession): string {
  const payload = Buffer.from(JSON.stringify(session), "utf8").toString("base64url");
  return `${payload}.${sign(payload)}`;
}

export function decodeSession(value: string | undefined): AdminSession | null {
  if (!value) {
    return null;
  }
  const [payload, signature] = value.split(".");
  if (!payload || !signature) {
    return null;
  }
  const expected = sign(payload);
  const left = Buffer.from(signature);
  const right = Buffer.from(expected);
  if (left.length !== right.length || !timingSafeEqual(left, right)) {
    return null;
  }
  try {
    const parsed = JSON.parse(Buffer.from(payload, "base64url").toString("utf8")) as AdminSession;
    if (!parsed.staffId || !Array.isArray(parsed.permissions)) {
      return null;
    }
    return parsed;
  } catch {
    return null;
  }
}

export async function getSession(): Promise<AdminSession | null> {
  const store = await cookies();
  return decodeSession(store.get(SESSION_COOKIE)?.value);
}

export async function setSessionCookie(session: AdminSession): Promise<void> {
  const store = await cookies();
  store.set(SESSION_COOKIE, encodeSession(session), {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  });
}

export async function clearSessionCookie(): Promise<void> {
  const store = await cookies();
  store.delete(SESSION_COOKIE);
}
