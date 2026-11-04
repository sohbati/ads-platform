import { NextResponse } from "next/server";
import { getSession } from "@/lib/auth/session";
import { env } from "@/config/env";

interface CdnCategory {
  id: number;
  parent: number | null;
  order: number;
  title: string;
  slug: string;
  isLeaf: boolean;
}

export async function GET() {
  const session = await getSession();
  if (!session) {
    return NextResponse.json(
      { error: "AUTH_UNAUTHENTICATED", statusCode: 401 },
      { status: 401 },
    );
  }

  const response = await fetch(`${env.cdnBaseUrl}/api/categories`, { cache: "no-store" });
  if (!response.ok) {
    return NextResponse.json(
      { error: "CATEGORIES_UNAVAILABLE", statusCode: 502 },
      { status: 502 },
    );
  }
  const categories = (await response.json()) as CdnCategory[];
  return NextResponse.json(
    categories.map(({ id, parent, order, title, slug, isLeaf }) => ({
      id,
      parent,
      order,
      title,
      slug,
      isLeaf,
    })),
  );
}
