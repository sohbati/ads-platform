"use client";

import type { ColumnsType } from "antd/es/table";
import { Button, Space } from "antd";
import { AppTable } from "@/components/ui/AppTable";
import { AppStatusTag } from "@/components/ui/AppStatusTag";
import { PermissionGuard } from "@/components/shared/PermissionGuard";
import { Permission } from "@/lib/permissions/codes";
import { formatDate } from "@/lib/formatters/date";
import { formatNumber } from "@/lib/formatters/number";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { AdvertisementListItem } from "../types/advertisement.types";
import { AdvertisementStatus } from "../constants/advertisement.constants";

interface AdvertisementTableProps {
  advertisements: AdvertisementListItem[];
  loading?: boolean;
  page: number;
  pageSize: number;
  total: number;
  permissions: readonly string[] | undefined;
  onPageChange: (page: number, pageSize: number) => void;
  onApprove: (item: AdvertisementListItem) => void;
  onReject: (item: AdvertisementListItem) => void;
}

function toneForStatus(status: AdvertisementListItem["status"]) {
  if (status === AdvertisementStatus.active) {
    return "success" as const;
  }
  if (status === AdvertisementStatus.pending) {
    return "processing" as const;
  }
  if (status === AdvertisementStatus.rejected || status === AdvertisementStatus.deleted) {
    return "error" as const;
  }
  return "default" as const;
}

export function AdvertisementTable({
  advertisements,
  loading,
  page,
  pageSize,
  total,
  permissions,
  onPageChange,
  onApprove,
  onReject,
}: AdvertisementTableProps) {
  const { t, locale } = useI18n();

  const columns: ColumnsType<AdvertisementListItem> = [
    { title: t("name"), dataIndex: "title" },
    {
      title: t("status"),
      dataIndex: "status",
      render: (status: AdvertisementListItem["status"]) => (
        <AppStatusTag label={status} tone={toneForStatus(status)} />
      ),
    },
    {
      title: t("createdAt"),
      dataIndex: "createdAt",
      render: (value: string) => formatDate(value, locale),
    },
    {
      title: t("actionsColumn"),
      key: "actions",
      render: (_, item) => (
        <Space>
          <PermissionGuard permission={Permission.advertisementApprove} permissions={permissions}>
            <Button size="small" type="primary" onClick={() => onApprove(item)}>
              {t("approve")}
            </Button>
          </PermissionGuard>
          <PermissionGuard permission={Permission.advertisementReject} permissions={permissions}>
            <Button size="small" danger onClick={() => onReject(item)}>
              {t("reject")}
            </Button>
          </PermissionGuard>
        </Space>
      ),
    },
  ];

  return (
    <AppTable<AdvertisementListItem>
      rowKey="id"
      columns={columns}
      dataSource={advertisements}
      loading={loading}
      emptyText={t("backendMissing")}
      pagination={{
        current: page,
        pageSize,
        total,
        showTotal: (count) => formatNumber(count, locale),
        onChange: onPageChange,
      }}
    />
  );
}
