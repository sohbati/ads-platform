"use client";

import { Empty } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";

interface AppEmptyProps {
  description?: string;
}

export function AppEmpty({ description }: AppEmptyProps) {
  const { t } = useI18n();
  return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={description ?? t("empty")} />;
}
