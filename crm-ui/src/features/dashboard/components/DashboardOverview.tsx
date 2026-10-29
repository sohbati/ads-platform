"use client";

import { Card, Col, Row, Typography } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { formatNumber } from "@/lib/formatters/number";
import { AppError } from "@/components/ui/AppError";
import { AppLoading } from "@/components/ui/AppLoading";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { useDashboardStats } from "../hooks/useDashboardStats";

export function DashboardOverview() {
  const { t, locale } = useI18n();
  const stats = useDashboardStats();

  return (
    <div>
      <AppPageHeader title={t("dashboard")} description={t("welcome")} />
      {stats.isLoading ? <AppLoading /> : null}
      {stats.error ? <AppError error={stats.error} onRetry={() => void stats.refetch()} /> : null}
      {stats.data ? (
        <>
          {!stats.data.connected ? (
            <Typography.Paragraph type="secondary">{t("backendMissing")}</Typography.Paragraph>
          ) : null}
          <Row gutter={[16, 16]}>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Typography.Text type="secondary">{t("pendingAds")}</Typography.Text>
                <Typography.Title level={3}>
                  {formatNumber(stats.data.pendingAds, locale)}
                </Typography.Title>
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Typography.Text type="secondary">{t("activeAds")}</Typography.Text>
                <Typography.Title level={3}>
                  {formatNumber(stats.data.activeAds, locale)}
                </Typography.Title>
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Typography.Text type="secondary">{t("totalCustomers")}</Typography.Text>
                <Typography.Title level={3}>
                  {formatNumber(stats.data.totalCustomers, locale)}
                </Typography.Title>
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Typography.Text type="secondary">{t("staffCount")}</Typography.Text>
                <Typography.Title level={3}>
                  {formatNumber(stats.data.staffCount, locale)}
                </Typography.Title>
              </Card>
            </Col>
          </Row>
        </>
      ) : null}
    </div>
  );
}
