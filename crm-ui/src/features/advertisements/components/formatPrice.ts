import { formatCurrency } from "@/lib/formatters/currency";
import type { MessageKey } from "@/lib/i18n/messages";
import type { AdvertisementListItem } from "../types/advertisement.types";

export function formatPrice(
  ad: Pick<AdvertisementListItem, "price_amount" | "price_type" | "currency">,
  locale: "fa" | "en",
  t: (key: MessageKey) => string,
): string {
  if (ad.price_type === "free") {
    return t("priceFree");
  }
  if (ad.price_amount === null || ad.price_type === "negotiable") {
    return t("priceNegotiable");
  }
  return formatCurrency(ad.price_amount, locale, ad.currency);
}
