import { NextResponse } from "next/server";
import { getSession } from "@/lib/auth/session";
import { emptyPage } from "@/lib/api/backend";

export async function GET(request: Request) {
  const session = await getSession();
  if (!session) {
    return NextResponse.json(
      { error: "AUTH_UNAUTHENTICATED", statusCode: 401 },
      { status: 401 },
    );
  }
  const url = new URL(request.url);
  const page = Number(url.searchParams.get("page") ?? "1");
  const pageSize = Number(url.searchParams.get("pageSize") ?? "20");
  return NextResponse.json(emptyPage(page, pageSize));
}
