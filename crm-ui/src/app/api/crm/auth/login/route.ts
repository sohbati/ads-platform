import { NextResponse } from "next/server";
import {
  demoSession,
  isDemoCredentials,
  setSessionCookie,
} from "@/lib/auth/session";

export async function POST(request: Request) {
  const body = (await request.json()) as { mobile?: string; password?: string };
  const mobile = body.mobile?.trim() ?? "";
  const password = body.password ?? "";

  if (!isDemoCredentials(mobile, password)) {
    return NextResponse.json(
      { error: "AUTH_INVALID_CREDENTIALS", statusCode: 401 },
      { status: 401 },
    );
  }

  await setSessionCookie(demoSession());
  return NextResponse.json({ ok: true });
}
