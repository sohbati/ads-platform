import { formatNumber } from "./number";

/** Ad prices are stored in whole tomans under the unofficial code IRT; IRR (rial) is labelled as such. */
export function formatCurrency(
  amount: number,
  locale: "fa" | "en",
  currency = "IRT",
): string {
  const number = formatNumber(amount, locale);
  if (currency === "IRT") {
    return locale === "fa" ? `${number} تومان` : `${number} Toman`;
  }
  if (currency === "IRR") {
    return locale === "fa" ? `${number} ریال` : `${number} IRR`;
  }
  return formatNumber(amount, locale, { style: "currency", currency });
}
