"use client";

import { useState } from "react";
import type { ColumnsType } from "antd/es/table";
import { Button, Space } from "antd";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { AppTable } from "@/components/ui/AppTable";
import { ConfirmDialog } from "@/components/shared/ConfirmDialog";
import { PermissionGuard } from "@/components/shared/PermissionGuard";
import { formatDate } from "@/lib/formatters/date";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useSession } from "@/lib/auth/SessionProvider";
import type { PermissionCode } from "@/lib/permissions/codes";
import type { NamedEntity, NamedEntityInput } from "../types/access.types";
import { NamedEntityForm } from "./NamedEntityForm";

interface NamedEntityScreenProps {
  title: string;
  items: NamedEntity[] | undefined;
  loading: boolean;
  error: unknown;
  writePermission: PermissionCode;
  saving?: boolean;
  deleting?: boolean;
  onRetry: () => void;
  onCreate: (input: NamedEntityInput, done: () => void) => void;
  onUpdate: (id: number, input: NamedEntityInput, done: () => void) => void;
  /** Whether codes may change after creation. */
  codeEditable?: boolean;
  onDelete: (id: number) => void;
  /** Extra per-row buttons shown before Edit/Delete. */
  rowActions?: (item: NamedEntity) => React.ReactNode;
}

export function NamedEntityScreen({
  title,
  items,
  loading,
  error,
  writePermission,
  saving,
  deleting,
  onRetry,
  onCreate,
  onUpdate,
  codeEditable,
  onDelete,
  rowActions,
}: NamedEntityScreenProps) {
  const { t, locale } = useI18n();
  const session = useSession();
  // `undefined` = closed, `null` = creating, item = editing.
  const [editing, setEditing] = useState<NamedEntity | null | undefined>(undefined);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const close = () => setEditing(undefined);

  const columns: ColumnsType<NamedEntity> = [
    {
      title: t("code"),
      dataIndex: "code",
      render: (code: string) => <span dir="ltr">{code}</span>,
    },
    { title: t("name"), dataIndex: "name" },
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
          {rowActions?.(item)}
          <PermissionGuard permission={writePermission} permissions={session.permissions}>
            <Button size="small" onClick={() => setEditing(item)}>
              {t("edit")}
            </Button>
            <Button danger size="small" onClick={() => setDeletingId(item.id)}>
              {t("delete")}
            </Button>
          </PermissionGuard>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <AppPageHeader
        title={title}
        extra={
          <PermissionGuard permission={writePermission} permissions={session.permissions}>
            <Button type="primary" onClick={() => setEditing(null)}>
              {t("create")}
            </Button>
          </PermissionGuard>
        }
      />
      {error ? <AppError error={error} onRetry={onRetry} /> : null}
      <AppTable<NamedEntity>
        rowKey="id"
        columns={columns}
        dataSource={items ?? []}
        loading={loading}
        pagination={false}
      />
      <NamedEntityForm
        open={editing !== undefined}
        item={editing}
        codeEditable={codeEditable}
        loading={saving}
        onCancel={close}
        onSubmit={(values) => {
          if (editing) {
            onUpdate(editing.id, values, close);
          } else {
            onCreate(values, close);
          }
        }}
      />
      <ConfirmDialog
        open={deletingId !== null}
        title={t("confirmDelete")}
        loading={deleting}
        onCancel={() => setDeletingId(null)}
        onConfirm={() => {
          if (deletingId !== null) {
            onDelete(deletingId);
          }
          setDeletingId(null);
        }}
      />
    </div>
  );
}
