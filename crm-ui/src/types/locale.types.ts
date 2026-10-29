export const LOCALES = ["fa", "en"] as const;

export type Locale = (typeof LOCALES)[number];

export type Direction = "rtl" | "ltr";

export function isLocale(value: string | undefined): value is Locale {
  return value === "fa" || value === "en";
}

export function directionForLocale(locale: Locale): Direction {
  return locale === "fa" ? "rtl" : "ltr";
}
