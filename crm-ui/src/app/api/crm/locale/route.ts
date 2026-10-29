import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { isLocale } from "@/types/locale.types";
import { LOCALE_COOKIE } from "@/lib/i18n/locale";

export async function POST(request: Request) {
  const body = (await request.json()) as { locale?: string };
  if (!isLocale(body.locale)) {
    return NextResponse.json(
      { error: "LOCALE_INVALID", statusCode: 400 },
      { status: 400 },
    );
  }
  const store = await cookies();
  store.set(LOCALE_COOKIE, body.locale, {
    path: "/",
    sameSite: "lax",
    maxAge: 60 * 60 * 24 * 365,
  });
  return NextResponse.json({ locale: body.locale });
}
