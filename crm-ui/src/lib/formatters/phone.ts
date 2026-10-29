import { formatDigits, toLatinDigits } from "./digits";

export function formatMobile(mobile: string, locale: "fa" | "en"): string {
  const digits = toLatinDigits(mobile).replace(/\D/g, "");
  const normalized = digits.startsWith("98") ? `0${digits.slice(2)}` : digits;
  const grouped =
    normalized.length === 11
      ? `${normalized.slice(0, 4)} ${normalized.slice(4, 7)} ${normalized.slice(7)}`
      : normalized;
  return formatDigits(grouped, locale);
}
