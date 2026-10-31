import { NextResponse } from "next/server";
import { getSession } from "@/lib/auth/session";
import { can } from "@/lib/auth/permissions";
import { requiredPermission, targetStaffId } from "@/lib/auth/panelAccess";
import { backendFetch } from "@/lib/api/backend";
import { AppApiError } from "@/lib/api/errors";

function fail(error: string, status: number) {
  return NextResponse.json({ error, statusCode: status }, { status });
}

async function proxy(request: Request, context: { params: Promise<{ path: string[] }> }) {
  try {
    const session = await getSession();
    if (!session) {
      return fail("AUTH_UNAUTHENTICATED", 401);
    }

    const { path: segments } = await context.params;
    const path = segments.join("/");
    const permission = requiredPermission(request.method, path);
    if (!permission || !can(session.permissions, permission)) {
      return fail("AUTH_FORBIDDEN", 403);
    }

    // Only a super admin may change another super admin's details, password, roles, or groups.
    const targetId = targetStaffId(request.method, path);
    if (targetId !== null && !session.isSuperAdmin) {
      const target = await backendFetch<{ is_super_admin: boolean }>({
        path: `/api/v1/panel/staff/${targetId}`,
      });
      if (target.is_super_admin) {
        return fail("AUTH_FORBIDDEN", 403);
      }
    }

    let body: unknown;
    if (request.method !== "GET" && request.method !== "HEAD") {
      const text = await request.text();
      body = text ? (JSON.parse(text) as unknown) : undefined;
    }

    const data = await backendFetch<unknown>({
      path: `/api/v1/panel/${path}`,
      method: request.method,
      body,
      search: new URL(request.url).searchParams,
    });
    return data === undefined
      ? new NextResponse(null, { status: 204 })
      : NextResponse.json(data);
  } catch (error) {
    if (error instanceof AppApiError) {
      return fail(error.code, error.status ?? 500);
    }
    if (error instanceof SyntaxError) {
      return fail("INVALID_JSON", 400);
    }
    return fail("BACKEND_UNAVAILABLE", 502);
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
