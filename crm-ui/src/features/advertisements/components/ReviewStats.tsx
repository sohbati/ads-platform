"use client";

import { Collapse, Table, Typography } from "antd";
import { can } from "@/lib/auth/permissions";
import { useSession } from "@/lib/auth/SessionProvider";
import { formatNumber } from "@/lib/formatters/number";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { Permission } from "@/lib/permissions/codes";
import { useReviewStats } from "../hooks/useAdvertisements";
import type { ReviewerStats } from "../types/advertisement.types";

/** Today's decisions: the signed-in reviewer's own count, and everyone's for queue managers. */
export function ReviewStats() {
  const { t, locale } = useI18n();
  const session = useSession();
  const reviewer = can(session.permissions, Permission.advertisementReview);
  const manager = can(session.permissions, Permission.advertisementReviewManage);
  const query = useReviewStats(reviewer || manager);
  const rows = query.data ?? [];
  const mine = rows.find((row) => row.staff_id === session.staffId);

  if (!reviewer && !manager) {
    return null;
  }

  const count = (value: number) => formatNumber(value, locale);

  return (
    <div className="app-filters">
      {reviewer ? (
        <Typography.Text type="secondary">
          {`${t("youToday")}: ${count(mine?.approved ?? 0)} ${t("approvedCount")} · ${count(
            mine?.rejected ?? 0,
          )} ${t("rejectedCount")}`}
        </Typography.Text>
      ) : null}
      {manager && rows.length > 0 ? (
        <Collapse
          size="small"
          style={{ marginTop: 8 }}
          items={[
            {
              key: "team",
              label: t("reviewedToday"),
              children: (
                <Table<ReviewerStats>
                  size="small"
                  rowKey="staff_id"
                  pagination={false}
                  dataSource={rows}
                  columns={[
                    { title: t("name"), dataIndex: "name" },
                    {
                      title: t("approvedCount"),
                      dataIndex: "approved",
                      render: (value: number) => count(value),
                    },
                    {
                      title: t("rejectedCount"),
                      dataIndex: "rejected",
                      render: (value: number) => count(value),
                    },
                  ]}
                />
              ),
            },
          ]}
        />
      ) : null}
    </div>
  );
}
