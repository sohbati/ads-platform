"use client";

import { useMemo } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import type { ColumnsType } from "antd/es/table";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { AppTable } from "@/components/ui/AppTable";
import { SearchInput } from "@/components/shared/SearchInput";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { MessageKey } from "@/lib/i18n/messages";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api/client";
import type { PagedResult } from "@/lib/api/types";

interface CatalogRow {
  id: number;
  name: string;
}

interface CatalogScreenProps {
  titleKey: MessageKey;
  queryKey: readonly unknown[];
}

export function CatalogScreen({ titleKey, queryKey }: CatalogScreenProps) {
  const { t } = useI18n();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const page = Number(searchParams.get("page") ?? "1");
  const pageSize = Number(searchParams.get("pageSize") ?? "20");
  const search = searchParams.get("search") ?? undefined;

  const query = useQuery({
    queryKey: [...queryKey, page, pageSize, search],
    queryFn: async () => {
      const { data } = await apiClient.get<PagedResult<CatalogRow>>("/catalog", {
        params: { page, pageSize, search },
      });
      return data;
    },
  });

  const columns = useMemo<ColumnsType<CatalogRow>>(
    () => [{ title: t("name"), dataIndex: "name" }],
    [t],
  );

  const replace = (nextPage: number, nextSize: number, nextSearch?: string) => {
    const params = new URLSearchParams();
    params.set("page", String(nextPage));
    params.set("pageSize", String(nextSize));
    if (nextSearch) {
      params.set("search", nextSearch);
    }
    router.replace(`${pathname}?${params.toString()}`);
  };

  return (
    <div>
      <AppPageHeader title={t(titleKey)} description={t("backendMissing")} />
      <div className="app-filters">
        <SearchInput value={search} onChange={(value) => replace(1, pageSize, value)} />
      </div>
      {query.error ? <AppError error={query.error} onRetry={() => void query.refetch()} /> : null}
      <AppTable<CatalogRow>
        rowKey="id"
        columns={columns}
        dataSource={query.data?.items ?? []}
        loading={query.isLoading}
        emptyText={t("backendMissing")}
        pagination={{
          current: page,
          pageSize,
          total: query.data?.total ?? 0,
          onChange: (nextPage, nextSize) => replace(nextPage, nextSize, search),
        }}
      />
    </div>
  );
}
