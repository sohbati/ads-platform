"use client";

import type { ColumnsType } from "antd/es/table";
import { Button, Space } from "antd";
import { AppTable } from "@/components/ui/AppTable";
import { AppStatusTag } from "@/components/ui/AppStatusTag";
import { PermissionGuard } from "@/components/shared/PermissionGuard";
import { Permission } from "@/lib/permissions/codes";
import { formatDate } from "@/lib/formatters/date";
import { formatMobile } from "@/lib/formatters/phone";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { CustomerListItem } from "../types/customer.types";
import { CustomerStatus } from "../constants/customer.constants";

interface CustomerTableProps {
  customers: CustomerListItem[];
  loading?: boolean;
  page: number;
  pageSize: number;
  total: number;
  permissions: readonly string[] | undefined;
  onPageChange: (page: number, pageSize: number) => void;
  onBlock: (customer: CustomerListItem) => void;
}

function toneForStatus(status: CustomerListItem["status"]) {
  if (status === CustomerStatus.ACTIVE) {
    return "success" as const;
  }
  if (status === CustomerStatus.SUSPENDED) {
    return "warning" as const;
  }
  return "error" as const;
}

export function CustomerTable({
  customers,
  loading,
  page,
  pageSize,
  total,
  permissions,
  onPageChange,
  onBlock,
}: CustomerTableProps) {
  const { t, locale } = useI18n();

  const columns: ColumnsType<CustomerListItem> = [
    { title: t("name"), dataIndex: "name" },
    {
      title: t("mobile"),
      dataIndex: "mobile",
      render: (mobile: string) => formatMobile(mobile, locale),
    },
    {
      title: t("status"),
      dataIndex: "status",
      render: (status: CustomerListItem["status"]) => (
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
      render: (_, customer) => (
        <PermissionGuard permission={Permission.customerBlock} permissions={permissions}>
          <Space>
            <Button danger size="small" onClick={() => onBlock(customer)}>
              {t("block")}
            </Button>
          </Space>
        </PermissionGuard>
      ),
    },
  ];

  return (
    <AppTable<CustomerListItem>
      rowKey="id"
      columns={columns}
      dataSource={customers}
      loading={loading}
      emptyText={t("backendMissing")}
      pagination={{
        current: page,
        pageSize,
        total,
        onChange: onPageChange,
      }}
    />
  );
}
