"use client";

import { useMemo, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { ConfirmDialog } from "@/components/shared/ConfirmDialog";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useSession } from "@/lib/auth/SessionProvider";
import { AdvertisementStatus } from "../constants/advertisement.constants";
import type {
  AdvertisementFiltersState,
  AdvertisementListItem,
} from "../types/advertisement.types";
import { useAdvertisements } from "../hooks/useAdvertisements";
import { AdvertisementFilters } from "./AdvertisementFilters";
import { AdvertisementTable } from "./AdvertisementTable";

function parseStatus(value: string | null) {
  return Object.values(AdvertisementStatus).find((status) => status === value);
}

function parseFilters(
  params: URLSearchParams,
  preset?: AdvertisementStatus,
): AdvertisementFiltersState {
  return {
    page: Number(params.get("page") ?? "1"),
    pageSize: Number(params.get("pageSize") ?? "20"),
    search: params.get("search") ?? undefined,
    status: parseStatus(params.get("status")) ?? preset,
  };
}

interface AdvertisementScreenProps {
  titleKey?: "advertisements" | "review";
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
  const session = useSession();
  const [approveItem, setApproveItem] = useState<AdvertisementListItem | null>(null);
  const [rejectItem, setRejectItem] = useState<AdvertisementListItem | null>(null);

  const filters = useMemo(
    () => parseFilters(new URLSearchParams(searchParams.toString()), presetStatus),
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
    if (next.status) {
      params.set("status", next.status);
    }
    router.replace(`${pathname}?${params.toString()}`);
  };

  return (
    <div>
      <AppPageHeader title={t(titleKey)} description={t("backendMissing")} />
      <AdvertisementFilters value={filters} onChange={replaceFilters} />
      {query.error ? <AppError error={query.error} onRetry={() => void query.refetch()} /> : null}
      <AdvertisementTable
        advertisements={query.data?.items ?? []}
        loading={query.isLoading}
        page={filters.page}
        pageSize={filters.pageSize}
        total={query.data?.total ?? 0}
        permissions={session.permissions}
        onPageChange={(page, pageSize) => replaceFilters({ ...filters, page, pageSize })}
        onApprove={setApproveItem}
        onReject={setRejectItem}
      />
      <ConfirmDialog
        open={approveItem !== null}
        title={t("confirmApprove")}
        onCancel={() => setApproveItem(null)}
        onConfirm={() => setApproveItem(null)}
      />
      <ConfirmDialog
        open={rejectItem !== null}
        title={t("confirmReject")}
        requireReason
        onCancel={() => setRejectItem(null)}
        onConfirm={() => setRejectItem(null)}
      />
    </div>
  );
}
