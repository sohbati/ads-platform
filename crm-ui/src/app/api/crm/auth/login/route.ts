import { NextResponse } from "next/server";
import { backendFetch } from "@/lib/api/backend";
import { AppApiError } from "@/lib/api/errors";
import { toLatinDigits } from "@/lib/formatters/digits";
import { createRateLimiter } from "@/lib/auth/rateLimit";
import { setSessionCookie } from "@/lib/auth/session";
import type { BackendStaffSession } from "@/lib/auth/session.types";

const MAX_ATTEMPTS = 5;
const WINDOW_MS = 15 * 60 * 1000;
const limiter = createRateLimiter(MAX_ATTEMPTS, WINDOW_MS);

function clientIp(request: Request): string {
  return request.headers.get("x-forwarded-for")?.split(",")[0]?.trim() || "unknown";
}

function fail(error: string, status: number, headers?: HeadersInit) {
  return NextResponse.json({ error, statusCode: status }, { status, headers });
}

export async function POST(request: Request) {
  let body: { mobile?: unknown; password?: unknown };
  try {
    body = (await request.json()) as typeof body;
  } catch {
    return fail("AUTH_INVALID_BODY", 400);
  }
  const mobile = typeof body.mobile === "string" ? body.mobile.trim() : "";
  const password = typeof body.password === "string" ? body.password : "";
  if (!mobile || !password) {
    return fail("AUTH_INVALID_CREDENTIALS", 401);
  }

  const key = `${clientIp(request)}|${toLatinDigits(mobile).replace(/\D/g, "")}`;
  const retryAfter = limiter.retryAfter(key);
  if (retryAfter > 0) {
    return fail("AUTH_TOO_MANY_ATTEMPTS", 429, { "Retry-After": String(retryAfter) });
  }

  try {
    const session = await backendFetch<BackendStaffSession>({
      path: "/api/v1/panel/auth/login",
      method: "POST",
      body: { mobile, password },
    });
    limiter.reset(key);
    await setSessionCookie(session.staff.id);
    return NextResponse.json({ ok: true });
  } catch (error) {
    if (error instanceof AppApiError && error.status === 401) {
      limiter.fail(key);
      return fail("AUTH_INVALID_CREDENTIALS", 401);
    }
    return fail("BACKEND_UNAVAILABLE", 502);
  }
}
