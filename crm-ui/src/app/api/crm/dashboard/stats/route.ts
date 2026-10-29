import { NextResponse } from "next/server";
import { getSession } from "@/lib/auth/session";

export async function GET() {
  const session = await getSession();
  if (!session) {
    return NextResponse.json(
      { error: "AUTH_UNAUTHENTICATED", statusCode: 401 },
      { status: 401 },
    );
  }
  return NextResponse.json({
    pendingAds: 0,
    activeAds: 0,
    totalCustomers: 0,
    staffCount: 0,
    connected: false,
  });
}
