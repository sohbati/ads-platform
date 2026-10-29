"use client";

import { useMemo, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { ConfirmDialog } from "@/components/shared/ConfirmDialog";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useSession } from "@/lib/auth/SessionProvider";
import { CustomerStatus } from "../constants/customer.constants";
import type { CustomerFiltersState, CustomerListItem } from "../types/customer.types";
import { useCustomers } from "../hooks/useCustomers";
import { CustomerFilters } from "./CustomerFilters";
import { CustomerTable } from "./CustomerTable";

function parseFilters(params: URLSearchParams): CustomerFiltersState {
  const status = params.get("status");
  return {
    page: Number(params.get("page") ?? "1"),
    pageSize: Number(params.get("pageSize") ?? "20"),
    search: params.get("search") ?? undefined,
    status:
      status === CustomerStatus.ACTIVE ||
      status === CustomerStatus.BLOCKED ||
      status === CustomerStatus.SUSPENDED
        ? status
        : undefined,
  };
}

export function CustomerScreen() {
  const { t } = useI18n();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const session = useSession();
  const [blocking, setBlocking] = useState<CustomerListItem | null>(null);

  const filters = useMemo(
    () => parseFilters(new URLSearchParams(searchParams.toString())),
    [searchParams],
  );
  const query = useCustomers(filters);

  const replaceFilters = (next: CustomerFiltersState) => {
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
      <AppPageHeader title={t("customers")} description={t("backendMissing")} />
      <CustomerFilters value={filters} onChange={replaceFilters} />
      {query.error ? <AppError error={query.error} onRetry={() => void query.refetch()} /> : null}
      <CustomerTable
        customers={query.data?.items ?? []}
        loading={query.isLoading}
        page={filters.page}
        pageSize={filters.pageSize}
        total={query.data?.total ?? 0}
        permissions={session.permissions}
        onPageChange={(page, pageSize) => replaceFilters({ ...filters, page, pageSize })}
        onBlock={setBlocking}
      />
      <ConfirmDialog
        open={blocking !== null}
        title={t("confirmBlock")}
        requireReason
        onCancel={() => setBlocking(null)}
        onConfirm={() => setBlocking(null)}
      />
    </div>
  );
}
