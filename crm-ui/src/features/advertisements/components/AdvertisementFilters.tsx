"use client";

import { useMemo } from "react";
import { Input, Select, Space, TreeSelect } from "antd";
import { useCategories } from "@/features/categories/hooks/useCategories";
import { useCities } from "@/features/locations/hooks/useCities";
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
  const categories = useCategories();
  const cities = useCities();

  const categoryTree = useMemo(
    () =>
      (categories.data ?? []).map((category) => ({
        id: category.id,
        pId: category.parent ?? 0,
        value: category.id,
        title: category.title,
      })),
    [categories.data],
  );
  const cityOptions = useMemo(
    () =>
      (cities.data ?? []).map((city) => ({
        value: city.id,
        label: city.province ? `${city.name} (${city.province})` : city.name,
      })),
    [cities.data],
  );

  return (
    <Space wrap className="app-filters">
      {/* Searches run on Enter, not on every keystroke: each one may scan a large table. */}
      <Input.Search
        key={value.search ?? ""}
        allowClear
        defaultValue={value.search}
        placeholder={t("searchAdsPlaceholder")}
        style={{ width: 300 }}
        onSearch={(search) => onChange({ ...value, search: search.trim() || undefined })}
        onChange={(event) => {
          if (event.target.value === "" && value.search) {
            onChange({ ...value, search: undefined });
          }
        }}
      />
      {statusLocked ? null : (
        <Select
          allowClear
          placeholder={t("status")}
          style={{ minWidth: 160 }}
          value={value.status}
          onChange={(status) => onChange({ ...value, status })}
          options={FILTERABLE_STATUSES.map((status) => ({
            value: status,
            label: t(STATUS_LABEL[status]),
          }))}
        />
      )}
      <TreeSelect
        allowClear
        showSearch
        treeDataSimpleMode
        placeholder={t("category")}
        style={{ minWidth: 200 }}
        loading={categories.isLoading}
        treeData={categoryTree}
        treeNodeFilterProp="title"
        value={value.category_id}
        onChange={(categoryId?: number) => onChange({ ...value, category_id: categoryId })}
      />
      <Select
        allowClear
        showSearch
        placeholder={t("city")}
        style={{ minWidth: 200 }}
        loading={cities.isLoading}
        options={cityOptions}
        optionFilterProp="label"
        value={value.city_id}
        onChange={(cityId?: number) => onChange({ ...value, city_id: cityId })}
      />
    </Space>
  );
}
