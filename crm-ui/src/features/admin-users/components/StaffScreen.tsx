"use client";

import { useMemo, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import type { ColumnsType } from "antd/es/table";
import { Button, Space, Tag } from "antd";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { AppTable } from "@/components/ui/AppTable";
import { ConfirmDialog } from "@/components/shared/ConfirmDialog";
import { PermissionGuard } from "@/components/shared/PermissionGuard";
import { SearchInput } from "@/components/shared/SearchInput";
import { can } from "@/lib/auth/permissions";
import { useSession } from "@/lib/auth/SessionProvider";
import { formatDate } from "@/lib/formatters/date";
import { formatDigits } from "@/lib/formatters/digits";
import { formatMobile } from "@/lib/formatters/phone";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { Permission } from "@/lib/permissions/codes";
import type { StaffFiltersState, StaffListItem } from "../types/staff.types";
import { useStaff, useStaffMutations } from "../hooks/useStaff";
import { StaffAccessDrawer } from "./StaffAccessDrawer";
import { StaffForm } from "./StaffForm";
import { StaffPasswordForm } from "./StaffPasswordForm";

export function StaffScreen() {
  const { t, locale } = useI18n();
  const session = useSession();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const canEdit = can(session.permissions, Permission.staffUpdate);

  // `undefined` = closed, `null` = creating, item = editing.
  const [editing, setEditing] = useState<StaffListItem | null | undefined>(undefined);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [accessFor, setAccessFor] = useState<StaffListItem | null>(null);
  const [passwordFor, setPasswordFor] = useState<StaffListItem | null>(null);

  // Only a super admin may change another super admin.
  const canManage = (item: StaffListItem) =>
    canEdit && (!item.is_super_admin || session.isSuperAdmin);

  const filters = useMemo<StaffFiltersState>(
    () => ({
      page: Number(searchParams.get("page") ?? "1"),
      pageSize: Number(searchParams.get("pageSize") ?? "20"),
      search: searchParams.get("search") ?? undefined,
    }),
    [searchParams],
  );
  const query = useStaff(filters);
  const mutations = useStaffMutations();

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
      render: (mobile: string) => <span dir="ltr">{formatMobile(mobile, locale)}</span>,
    },
    {
      title: t("nationalId"),
      dataIndex: "national_id",
      render: (value: string) => formatDigits(value, locale),
    },
    {
      title: t("status"),
      key: "status",
      render: (_, item) => (
        <Space size={4} wrap>
          <Tag color={item.is_active ? "green" : "default"}>
            {item.is_active ? t("active") : t("inactive")}
          </Tag>
          {item.is_super_admin ? <Tag color="gold">{t("superAdmin")}</Tag> : null}
          {item.is_reviewer ? <Tag color="purple">{t("reviewer")}</Tag> : null}
        </Space>
      ),
    },
    {
      title: t("createdAt"),
      dataIndex: "created_at",
      render: (value: string) => formatDate(value, locale),
    },
    {
      title: t("actionsColumn"),
      key: "actions",
      render: (_, item) => (
        <Space size="small" wrap>
          {item.is_super_admin ? null : (
            <Button size="small" onClick={() => setAccessFor(item)}>
              {t("staffAccess")}
            </Button>
          )}
          {canManage(item) ? (
            <>
              <Button size="small" onClick={() => setEditing(item)}>
                {t("edit")}
              </Button>
              <Button size="small" onClick={() => setPasswordFor(item)}>
                {t("setPassword")}
              </Button>
            </>
          ) : null}
          {canEdit && !item.is_super_admin ? (
            <Button danger size="small" onClick={() => setDeletingId(item.id)}>
              {t("delete")}
            </Button>
          ) : null}
        </Space>
      ),
    },
  ];

  const saving = mutations.create.isPending || mutations.update.isPending;

  return (
    <div>
      <AppPageHeader
        title={t("adminUsers")}
        extra={
          <PermissionGuard
            permission={Permission.staffUpdate}
            permissions={session.permissions}
          >
            <Button type="primary" onClick={() => setEditing(null)}>
              {t("addStaff")}
            </Button>
          </PermissionGuard>
        }
      />
      <div className="app-filters">
        <SearchInput
          value={filters.search}
          onChange={(search) => replaceFilters({ ...filters, page: 1, search })}
        />
      </div>
      {query.error ? (
        <AppError error={query.error} onRetry={() => void query.refetch()} />
      ) : null}
      <AppTable<StaffListItem>
        rowKey="id"
        columns={columns}
        dataSource={query.data?.items ?? []}
        loading={query.isLoading}
        pagination={{
          current: filters.page,
          pageSize: filters.pageSize,
          total: query.data?.total ?? 0,
          onChange: (page, pageSize) => replaceFilters({ ...filters, page, pageSize }),
        }}
      />
      <StaffForm
        open={editing !== undefined}
        staff={editing}
        loading={saving}
        onCancel={() => setEditing(undefined)}
        onSubmit={(input) => {
          const options = { onSuccess: () => setEditing(undefined) };
          if (editing) {
            mutations.update.mutate({ id: editing.id, input }, options);
          } else {
            mutations.create.mutate(input, options);
          }
        }}
      />
      <StaffPasswordForm
        staff={passwordFor}
        loading={mutations.setPassword.isPending}
        onCancel={() => setPasswordFor(null)}
        onSubmit={(password) => {
          if (passwordFor) {
            mutations.setPassword.mutate(
              { id: passwordFor.id, password },
              { onSuccess: () => setPasswordFor(null) },
            );
          }
        }}
      />
      <StaffAccessDrawer
        staff={accessFor}
        canEdit={canEdit}
        onClose={() => setAccessFor(null)}
      />
      <ConfirmDialog
        open={deletingId !== null}
        title={t("confirmDeleteStaff")}
        loading={mutations.remove.isPending}
        onCancel={() => setDeletingId(null)}
        onConfirm={() => {
          if (deletingId !== null) {
            mutations.remove.mutate(deletingId, { onSettled: () => setDeletingId(null) });
          }
        }}
      />
    </div>
  );
}
