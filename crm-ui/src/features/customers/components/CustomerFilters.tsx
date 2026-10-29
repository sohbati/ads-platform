"use client";

import { Select, Space } from "antd";
import { SearchInput } from "@/components/shared/SearchInput";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { CustomerStatus } from "../constants/customer.constants";
import type { CustomerFiltersState } from "../types/customer.types";

interface CustomerFiltersProps {
  value: CustomerFiltersState;
  onChange: (next: CustomerFiltersState) => void;
}

export function CustomerFilters({ value, onChange }: CustomerFiltersProps) {
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
        options={[
          { value: CustomerStatus.ACTIVE, label: "ACTIVE" },
          { value: CustomerStatus.BLOCKED, label: "BLOCKED" },
          { value: CustomerStatus.SUSPENDED, label: "SUSPENDED" },
        ]}
      />
    </Space>
  );
}
