import { NextResponse } from "next/server";
import { backendFetch } from "@/lib/api/backend";
import { clearSessionCookie, getSession } from "@/lib/auth/session";

export async function POST() {
  const session = await getSession().catch(() => null);
  if (session) {
    // Hand back any ad the reviewer was working on, so it does not wait for the claim to expire.
    await backendFetch({
      path: "/api/v1/panel/ads/review/release-all",
      method: "POST",
      headers: { "X-Staff-Id": String(session.staffId) },
    }).catch(() => undefined);
  }
  await clearSessionCookie();
  return NextResponse.json({ ok: true });
}
