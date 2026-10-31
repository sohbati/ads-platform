import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { SESSION_COOKIE } from "@/lib/auth/constants";

const PUBLIC_PATHS = ["/login"];

function isPublic(pathname: string): boolean {
  return (
    PUBLIC_PATHS.includes(pathname) ||
    pathname.startsWith("/api/crm/auth/login") ||
    pathname.startsWith("/api/crm/locale") ||
    pathname.startsWith("/_next") ||
    pathname.startsWith("/logo-mark")
  );
}

// Only checks that a cookie exists; the server validates it. A stale cookie therefore must not
// bounce /login to /dashboard here, or it would loop. The login page redirects valid sessions.
export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const hasSession = Boolean(request.cookies.get(SESSION_COOKIE)?.value);

  if (pathname === "/") {
    const url = request.nextUrl.clone();
    url.pathname = hasSession ? "/dashboard" : "/login";
    return NextResponse.redirect(url);
  }

  if (!hasSession && !isPublic(pathname)) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  return NextResponse.next({ request });
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};
