"use client";

import { useMemo } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import type { ColumnsType } from "antd/es/table";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { AppTable } from "@/components/ui/AppTable";
import { SearchInput } from "@/components/shared/SearchInput";
import { formatDate } from "@/lib/formatters/date";
import { formatMobile } from "@/lib/formatters/phone";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { StaffFiltersState, StaffListItem } from "../types/staff.types";
import { useStaff } from "../hooks/useStaff";

export function StaffScreen() {
  const { t, locale } = useI18n();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const filters = useMemo<StaffFiltersState>(
    () => ({
      page: Number(searchParams.get("page") ?? "1"),
      pageSize: Number(searchParams.get("pageSize") ?? "20"),
      search: searchParams.get("search") ?? undefined,
    }),
    [searchParams],
  );
  const query = useStaff(filters);

  const replaceFilters = (next: StaffFiltersState) => {
    const params = new URLSearchParams();
    params.set("page", String(next.page));
    params.set("pageSize", String(next.pageSize));
    if (next.search) {
      params.set("search", next.search);
    }
    router.replace(`${pathname}?${params.toString()}`);
  };

  const columns: ColumnsType<StaffListItem> = [
    { title: t("name"), dataIndex: "name" },
    {
      title: t("mobile"),
      dataIndex: "mobile",
      render: (mobile: string) => formatMobile(mobile, locale),
    },
    {
      title: t("createdAt"),
      dataIndex: "createdAt",
      render: (value: string) => formatDate(value, locale),
    },
  ];

  return (
    <div>
      <AppPageHeader title={t("adminUsers")} description={t("backendMissing")} />
      <div className="app-filters">
        <SearchInput
          value={filters.search}
          onChange={(search) => replaceFilters({ ...filters, page: 1, search })}
        />
      </div>
      {query.error ? <AppError error={query.error} onRetry={() => void query.refetch()} /> : null}
      <AppTable<StaffListItem>
        rowKey="id"
        columns={columns}
        dataSource={query.data?.items ?? []}
        loading={query.isLoading}
        emptyText={t("backendMissing")}
        pagination={{
          current: filters.page,
          pageSize: filters.pageSize,
          total: query.data?.total ?? 0,
          onChange: (page, pageSize) => replaceFilters({ ...filters, page, pageSize }),
        }}
      />
    </div>
  );
}
