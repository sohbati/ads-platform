"use client";

import { Segmented } from "antd";
import { useRouter } from "next/navigation";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { Locale } from "@/types/locale.types";
import { apiClient } from "@/lib/api/client";

export function LocaleSwitcher() {
  const { locale, t } = useI18n();
  const router = useRouter();

  return (
    <Segmented<Locale>
      value={locale}
      options={[
        { label: t("persian"), value: "fa" },
        { label: t("english"), value: "en" },
      ]}
      onChange={async (next) => {
        await apiClient.post("/locale", { locale: next });
        router.refresh();
      }}
    />
  );
}
