import { NextResponse } from "next/server";
import { getSession } from "@/lib/auth/session";
import { env } from "@/config/env";

interface CdnLocation {
  id: number;
  name: string;
  parent: number | null;
  type: string;
}

// Location types in the CDN data: 1 = province, 2 = city, 4 = neighborhood. Ads are filed by city.
const PROVINCE = "1";
const CITY = "2";

export async function GET() {
  const session = await getSession();
  if (!session) {
    return NextResponse.json(
      { error: "AUTH_UNAUTHENTICATED", statusCode: 401 },
      { status: 401 },
    );
  }

  const response = await fetch(`${env.cdnBaseUrl}/api/cities`, { cache: "no-store" });
  if (!response.ok) {
    return NextResponse.json(
      { error: "CITIES_UNAVAILABLE", statusCode: 502 },
      { status: 502 },
    );
  }
  const locations = (await response.json()) as CdnLocation[];
  const provinces = new Map(
    locations.filter((item) => item.type === PROVINCE).map((item) => [item.id, item.name]),
  );
  return NextResponse.json(
    locations
      .filter((item) => item.type === CITY)
      .map((item) => ({
        id: item.id,
        name: item.name,
        province: item.parent === null ? null : (provinces.get(item.parent) ?? null),
      })),
  );
}
