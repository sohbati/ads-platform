"use client";

import { Spin } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";

export function AppLoading() {
  const { t } = useI18n();
  return (
    <div className="app-loading" aria-label={t("loading")}>
      <Spin />
    </div>
  );
}
