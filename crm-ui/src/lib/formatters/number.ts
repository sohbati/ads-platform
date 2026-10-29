import { formatDigits } from "./digits";

export function formatNumber(
  value: number,
  locale: "fa" | "en",
  options?: Intl.NumberFormatOptions,
): string {
  const formatted = new Intl.NumberFormat(
    locale === "fa" ? "fa-IR" : "en-US",
    options,
  ).format(value);
  return formatDigits(formatted, locale);
}
