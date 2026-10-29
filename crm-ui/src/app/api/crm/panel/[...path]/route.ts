import { NextResponse } from "next/server";
import { getSession } from "@/lib/auth/session";
import { backendFetch } from "@/lib/api/backend";
import { AppApiError } from "@/lib/api/errors";

async function proxy(
  request: Request,
  context: { params: Promise<{ path: string[] }> },
) {
  const session = await getSession();
  if (!session) {
    return NextResponse.json(
      { error: "AUTH_UNAUTHENTICATED", statusCode: 401 },
      { status: 401 },
    );
  }

  const { path } = await context.params;
  const url = new URL(request.url);
  let body: unknown;
  if (request.method !== "GET" && request.method !== "HEAD") {
    const text = await request.text();
    body = text ? (JSON.parse(text) as unknown) : undefined;
  }

  try {
    const data = await backendFetch<unknown>({
      path: `/api/v1/panel/${path.join("/")}`,
      method: request.method,
      body,
      search: url.searchParams,
    });
    return NextResponse.json(data);
  } catch (error) {
    if (error instanceof AppApiError) {
      return NextResponse.json(
        { error: error.code, statusCode: error.status ?? 500 },
        { status: error.status ?? 500 },
      );
    }
    return NextResponse.json(
      { error: "BACKEND_UNAVAILABLE", statusCode: 502 },
      { status: 502 },
    );
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
