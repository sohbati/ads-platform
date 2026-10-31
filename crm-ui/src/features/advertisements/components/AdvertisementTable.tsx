"use client";

import type { ColumnsType } from "antd/es/table";
import { Button, Typography } from "antd";
import { AppTable } from "@/components/ui/AppTable";
import { formatDateTime } from "@/lib/formatters/date";
import { formatDigits } from "@/lib/formatters/digits";
import { formatNumber } from "@/lib/formatters/number";
import { formatMobile } from "@/lib/formatters/phone";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { AdvertisementListItem } from "../types/advertisement.types";
import { AdvertisementStatusTag } from "./AdvertisementStatusTag";
import { formatPrice } from "./formatPrice";

interface AdvertisementTableProps {
  advertisements: AdvertisementListItem[];
  loading?: boolean;
  page: number;
  pageSize: number;
  total: number;
  /** Review queue: show how long each ad has waited instead of when it last changed. */
  queue?: boolean;
  onPageChange: (page: number, pageSize: number) => void;
  onOpen: (item: AdvertisementListItem) => void;
}

export function AdvertisementTable({
  advertisements,
  loading,
  page,
  pageSize,
  total,
  queue,
  onPageChange,
  onOpen,
}: AdvertisementTableProps) {
  const { t, locale } = useI18n();

  const columns: ColumnsType<AdvertisementListItem> = [
    {
      title: t("adId"),
      dataIndex: "id",
      render: (id: number) => formatDigits(id, locale),
    },
    {
      title: t("title"),
      dataIndex: "title",
      render: (title: string, item) => (
        <Typography.Link onClick={() => onOpen(item)}>{title}</Typography.Link>
      ),
    },
    {
      title: t("category"),
      key: "category",
      render: (_, item) => item.category_title || formatDigits(item.category_id, locale),
    },
    {
      title: t("city"),
      key: "city",
      render: (_, item) => item.city_name || formatDigits(item.city_id, locale),
    },
    {
      title: t("customer"),
      key: "customer",
      render: (_, item) => (
        <div>
          <div>{item.customer_name}</div>
          <Typography.Text type="secondary" dir="ltr">
            {formatMobile(item.customer_mobile, locale)}
          </Typography.Text>
        </div>
      ),
    },
    {
      title: t("price"),
      key: "price",
      render: (_, item) => formatPrice(item, locale, t),
    },
    {
      title: t("status"),
      key: "status",
      render: (_, item) => <AdvertisementStatusTag ad={item} />,
    },
    queue
      ? {
          title: t("waitingSince"),
          dataIndex: "pending_since",
          render: (value: string | null) => formatDateTime(value, locale),
        }
      : {
          title: t("updatedAt"),
          dataIndex: "updated_at",
          render: (value: string) => formatDateTime(value, locale),
        },
    {
      title: t("actionsColumn"),
      key: "actions",
      render: (_, item) => (
        <Button size="small" onClick={() => onOpen(item)}>
          {t("view")}
        </Button>
      ),
    },
  ];

  return (
    <AppTable<AdvertisementListItem>
      rowKey="id"
      columns={columns}
      dataSource={advertisements}
      loading={loading}
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
