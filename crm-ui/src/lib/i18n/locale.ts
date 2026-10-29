import {
  directionForLocale,
  isLocale,
  type Direction,
  type Locale,
} from "@/types/locale.types";

export const LOCALE_COOKIE = "ruab_locale";
export const DEFAULT_LOCALE: Locale = "fa";

export function resolveLocale(value: string | undefined): Locale {
  return isLocale(value) ? value : DEFAULT_LOCALE;
}

export function localeMeta(locale: Locale): { lang: Locale; dir: Direction } {
  return { lang: locale, dir: directionForLocale(locale) };
}
