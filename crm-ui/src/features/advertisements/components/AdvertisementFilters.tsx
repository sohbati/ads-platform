"use client";

import { Select, Space } from "antd";
import { SearchInput } from "@/components/shared/SearchInput";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { AdvertisementStatus } from "../constants/advertisement.constants";
import type { AdvertisementFiltersState } from "../types/advertisement.types";

interface AdvertisementFiltersProps {
  value: AdvertisementFiltersState;
  onChange: (next: AdvertisementFiltersState) => void;
}

export function AdvertisementFilters({ value, onChange }: AdvertisementFiltersProps) {
  const { t } = useI18n();
  return (
    <Space wrap className="app-filters">
      <SearchInput
        value={value.search}
        onChange={(search) => onChange({ ...value, page: 1, search })}
      />
      <Select
        allowClear
        placeholder={t("status")}
        style={{ minWidth: 160 }}
        value={value.status}
        onChange={(status) => onChange({ ...value, page: 1, status })}
        options={Object.values(AdvertisementStatus).map((status) => ({
          value: status,
          label: status,
        }))}
      />
    </Space>
  );
}
