"use client";

import { useMemo, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { AdvertisementStatus } from "../constants/advertisement.constants";
import type { AdvertisementFiltersState } from "../types/advertisement.types";
import { useAdvertisements } from "../hooks/useAdvertisements";
import { AdvertisementDrawer } from "./AdvertisementDrawer";
import { AdvertisementFilters } from "./AdvertisementFilters";
import { AdvertisementTable } from "./AdvertisementTable";

function parseStatus(value: string | null) {
  return Object.values(AdvertisementStatus).find((status) => status === value);
}

interface AdvertisementScreenProps {
  titleKey?: "advertisements" | "review";
  /** Fixes the status filter, e.g. the review queue shows only pending ads. */
  presetStatus?: AdvertisementStatus;
}

export function AdvertisementScreen({
  titleKey = "advertisements",
  presetStatus,
}: AdvertisementScreenProps) {
  const { t } = useI18n();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [openId, setOpenId] = useState<number | null>(null);

  const filters = useMemo<AdvertisementFiltersState>(
    () => ({
      page: Number(searchParams.get("page") ?? "1"),
      pageSize: Number(searchParams.get("pageSize") ?? "20"),
      search: searchParams.get("search") ?? undefined,
      status: presetStatus ?? parseStatus(searchParams.get("status")),
    }),
    [searchParams, presetStatus],
  );
  const query = useAdvertisements(filters);

  const replaceFilters = (next: AdvertisementFiltersState) => {
    const params = new URLSearchParams();
    params.set("page", String(next.page));
    params.set("pageSize", String(next.pageSize));
    if (next.search) {
      params.set("search", next.search);
    }
    if (next.status && !presetStatus) {
      params.set("status", next.status);
    }
    router.replace(`${pathname}?${params.toString()}`);
  };

  return (
    <div>
      <AppPageHeader title={t(titleKey)} />
      <AdvertisementFilters
        value={filters}
        statusLocked={Boolean(presetStatus)}
        onChange={replaceFilters}
      />
      {query.error ? (
        <AppError error={query.error} onRetry={() => void query.refetch()} />
      ) : null}
      <AdvertisementTable
        advertisements={query.data?.items ?? []}
        loading={query.isLoading}
        page={filters.page}
        pageSize={filters.pageSize}
        total={query.data?.total ?? 0}
        queue={presetStatus === AdvertisementStatus.pending}
        onPageChange={(page, pageSize) => replaceFilters({ ...filters, page, pageSize })}
        onOpen={(item) => setOpenId(item.id)}
      />
      <AdvertisementDrawer adId={openId} onClose={() => setOpenId(null)} />
    </div>
  );
}
