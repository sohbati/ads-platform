"use client";

import { Select, Space } from "antd";
import { SearchInput } from "@/components/shared/SearchInput";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { FILTERABLE_STATUSES, STATUS_LABEL } from "../constants/advertisement.constants";
import type { AdvertisementFiltersState } from "../types/advertisement.types";

interface AdvertisementFiltersProps {
  value: AdvertisementFiltersState;
  /** Hides the status filter, e.g. on the review queue. */
  statusLocked?: boolean;
  onChange: (next: AdvertisementFiltersState) => void;
}

export function AdvertisementFilters({
  value,
  statusLocked,
  onChange,
}: AdvertisementFiltersProps) {
  const { t } = useI18n();
  return (
    <Space wrap className="app-filters">
      <SearchInput
        value={value.search}
        onChange={(search) => onChange({ ...value, page: 1, search })}
      />
      {statusLocked ? null : (
        <Select
          allowClear
          placeholder={t("status")}
          style={{ minWidth: 180 }}
          value={value.status}
          onChange={(status) => onChange({ ...value, page: 1, status })}
          options={FILTERABLE_STATUSES.map((status) => ({
            value: status,
            label: t(STATUS_LABEL[status]),
          }))}
        />
      )}
    </Space>
  );
}
