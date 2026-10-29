"use client";

import { useState } from "react";
import type { ColumnsType } from "antd/es/table";
import { Button } from "antd";
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
  creating?: boolean;
  deleting?: boolean;
  onRetry: () => void;
  onCreate: (input: NamedEntityInput) => void;
  onDelete: (id: number) => void;
}

export function NamedEntityScreen({
  title,
  items,
  loading,
  error,
  writePermission,
  creating,
  deleting,
  onRetry,
  onCreate,
  onDelete,
}: NamedEntityScreenProps) {
  const { t, locale } = useI18n();
  const session = useSession();
  const [creatingOpen, setCreatingOpen] = useState(false);
  const [deletingId, setDeletingId] = useState<number | null>(null);

  const columns: ColumnsType<NamedEntity> = [
    { title: t("code"), dataIndex: "code" },
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
        <PermissionGuard permission={writePermission} permissions={session.permissions}>
          <Button danger size="small" onClick={() => setDeletingId(item.id)}>
            {t("delete")}
          </Button>
        </PermissionGuard>
      ),
    },
  ];

  return (
    <div>
      <AppPageHeader
        title={title}
        extra={
          <PermissionGuard permission={writePermission} permissions={session.permissions}>
            <Button type="primary" onClick={() => setCreatingOpen(true)}>
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
        open={creatingOpen}
        loading={creating}
        onCancel={() => setCreatingOpen(false)}
        onSubmit={(values) => {
          onCreate(values);
          setCreatingOpen(false);
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
