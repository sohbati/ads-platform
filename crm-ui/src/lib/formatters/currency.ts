import { formatNumber } from "./number";

export function formatCurrency(
  amount: number,
  locale: "fa" | "en",
  currency = "IRR",
): string {
  if (currency === "IRR") {
    const number = formatNumber(amount, locale);
    return locale === "fa" ? `${number} تومان` : `${number} IRR`;
  }
  return formatNumber(amount, locale, { style: "currency", currency });
}
