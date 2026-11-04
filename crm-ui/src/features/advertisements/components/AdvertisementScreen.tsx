"use client";

import { useMemo, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { App, Button, Space, Typography } from "antd";
import { AppEmpty } from "@/components/ui/AppEmpty";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { can } from "@/lib/auth/permissions";
import { useSession } from "@/lib/auth/SessionProvider";
import { formatNumber } from "@/lib/formatters/number";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { Permission } from "@/lib/permissions/codes";
import { AdvertisementStatus } from "../constants/advertisement.constants";
import type { AdvertisementFiltersState } from "../types/advertisement.types";
import { useAdvertisements, useClaimActions } from "../hooks/useAdvertisements";
import { AdvertisementDrawer } from "./AdvertisementDrawer";
import { AdvertisementFilters } from "./AdvertisementFilters";
import { AdvertisementTable } from "./AdvertisementTable";
import { ReviewStats } from "./ReviewStats";

// The queue refreshes on its own so reviewers see which ads others have taken.
const QUEUE_REFRESH_MS = 30_000;
const PAGE_SIZE = 20;

function parseStatus(value: string | null) {
  return Object.values(AdvertisementStatus).find((status) => status === value);
}

function parseId(value: string | null) {
  const id = Number(value);
  return Number.isInteger(id) && id > 0 ? id : undefined;
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
  const { t, locale } = useI18n();
  const { message } = App.useApp();
  const session = useSession();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [openId, setOpenId] = useState<number | null>(null);

  const filters = useMemo<AdvertisementFiltersState>(
    () => ({
      pageSize: PAGE_SIZE,
      search: searchParams.get("search") ?? undefined,
      status: presetStatus ?? parseStatus(searchParams.get("status")),
      category_id: parseId(searchParams.get("category")),
      city_id: parseId(searchParams.get("city")),
    }),
    [searchParams, presetStatus],
  );

  // Pages follow cursors, kept as a stack so Previous can step back. New filters start over.
  const filtersKey = JSON.stringify(filters);
  const [paging, setPaging] = useState({ key: filtersKey, cursors: [] as string[] });
  const cursors = paging.key === filtersKey ? paging.cursors : [];
  const cursor = cursors.at(-1);

  // Outside the pending queue the list needs a search or filter; the backend refuses otherwise.
  const filtered = Boolean(
    filters.status === AdvertisementStatus.pending ||
      filters.search ||
      filters.category_id ||
      filters.city_id,
  );
  const queue = presetStatus === AdvertisementStatus.pending;
  const reviewer = queue && can(session.permissions, Permission.advertisementReview);
  const query = useAdvertisements(
    { ...filters, cursor },
    { enabled: filtered, refetchInterval: queue ? QUEUE_REFRESH_MS : undefined },
  );
  const { next } = useClaimActions();

  const openNext = () =>
    next.mutate(undefined, {
      onSuccess: (ad) => {
        setOpenId(ad?.id ?? null);
        if (!ad) {
          message.info(t("queueEmpty"));
        }
      },
    });

  const replaceFilters = (nextFilters: AdvertisementFiltersState) => {
    const params = new URLSearchParams();
    if (nextFilters.search) {
      params.set("search", nextFilters.search);
    }
    if (nextFilters.status && !presetStatus) {
      params.set("status", nextFilters.status);
    }
    if (nextFilters.category_id) {
      params.set("category", String(nextFilters.category_id));
    }
    if (nextFilters.city_id) {
      params.set("city", String(nextFilters.city_id));
    }
    router.replace(`${pathname}?${params.toString()}`);
  };

  const page = query.data;
  const total = page
    ? `${formatNumber(page.total, locale)}${page.total_capped ? "+" : ""} ${t("results")}`
    : null;

  return (
    <div>
      <AppPageHeader
        title={t(titleKey)}
        extra={
          reviewer ? (
            <Button type="primary" loading={next.isPending} onClick={openNext}>
              {t("nextAd")}
            </Button>
          ) : null
        }
      />
      {queue ? <ReviewStats /> : null}
      <AdvertisementFilters
        value={filters}
        statusLocked={Boolean(presetStatus)}
        onChange={replaceFilters}
      />
      {query.error ? (
        <AppError error={query.error} onRetry={() => void query.refetch()} />
      ) : null}
      {filtered ? (
        <>
          <AdvertisementTable
            advertisements={page?.items ?? []}
            loading={query.isLoading || query.isPlaceholderData}
            queue={queue}
            onOpen={(item) => setOpenId(item.id)}
          />
          <Space className="app-pager" wrap>
            <Typography.Text type="secondary">{total}</Typography.Text>
            <Button
              disabled={cursors.length === 0}
              onClick={() => setPaging({ key: filtersKey, cursors: cursors.slice(0, -1) })}
            >
              {t("previousPage")}
            </Button>
            <Button
              disabled={!page?.next_cursor || query.isPlaceholderData}
              onClick={() =>
                page?.next_cursor &&
                setPaging({ key: filtersKey, cursors: [...cursors, page.next_cursor] })
              }
            >
              {t("nextPage")}
            </Button>
          </Space>
        </>
      ) : (
        <AppEmpty description={t("searchToFindAds")} />
      )}
      <AdvertisementDrawer
        adId={openId}
        onClose={() => setOpenId(null)}
        onDecided={reviewer ? openNext : undefined}
      />
    </div>
  );
}
