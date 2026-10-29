"use client";

import { Card, Col, Row, Typography } from "antd";
import Link from "next/link";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { useI18n } from "@/lib/i18n/I18nProvider";

export function SettingsScreen() {
  const { t } = useI18n();
  return (
    <div>
      <AppPageHeader title={t("settings")} description={t("access")} />
      <Row gutter={[16, 16]}>
        <Col xs={24} md={8}>
          <Link href="/settings/roles">
            <Card hoverable>
              <Typography.Title level={4}>{t("roles")}</Typography.Title>
            </Card>
          </Link>
        </Col>
        <Col xs={24} md={8}>
          <Link href="/settings/actions">
            <Card hoverable>
              <Typography.Title level={4}>{t("actions")}</Typography.Title>
            </Card>
          </Link>
        </Col>
        <Col xs={24} md={8}>
          <Link href="/settings/groups">
            <Card hoverable>
              <Typography.Title level={4}>{t("groups")}</Typography.Title>
            </Card>
          </Link>
        </Col>
      </Row>
    </div>
  );
}
