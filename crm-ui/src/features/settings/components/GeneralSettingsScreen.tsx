"use client";

import { useState } from "react";
import { Button, Card, Form, InputNumber, Space, Typography } from "antd";
import { AppError } from "@/components/ui/AppError";
import { AppLoading } from "@/components/ui/AppLoading";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { can } from "@/lib/auth/permissions";
import { useSession } from "@/lib/auth/SessionProvider";
import { formatDateTime } from "@/lib/formatters/date";
import { formatDigits } from "@/lib/formatters/digits";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { Permission } from "@/lib/permissions/codes";
import { useSettings, useUpdateSetting } from "../hooks/useSettings";
import { SettingKey, type AppSetting } from "../types/settings.types";

export function GeneralSettingsScreen() {
  const { t } = useI18n();
  const query = useSettings();
  const pendingDays = query.data?.find(
    (item) => item.key === SettingKey.adsPendingMaxDays,
  );

  return (
    <div>
      <AppPageHeader title={t("generalSettings")} />
      {query.isLoading ? <AppLoading /> : null}
      {query.error ? (
        <AppError error={query.error} onRetry={() => void query.refetch()} />
      ) : null}
      {pendingDays ? <PendingDaysCard setting={pendingDays} /> : null}
    </div>
  );
}

function PendingDaysCard({ setting }: { setting: AppSetting }) {
  const { t, locale } = useI18n();
  const session = useSession();
  const update = useUpdateSetting();
  const canEdit = can(session.permissions, Permission.settingsUpdate);
  // null = untouched, so the saved value shows.
  const [draft, setDraft] = useState<number | null>(null);
  const saved = Number(setting.value);
  const value = draft ?? saved;

  return (
    <Card style={{ maxWidth: 640 }}>
      <Form layout="vertical">
        <Form.Item
          label={t("pendingMaxDays")}
          help={`${t("pendingMaxDaysHint")} ${t("defaultValue")}: ${formatDigits(setting.default, locale)}`}
        >
          <Space wrap>
            <InputNumber
              min={setting.min}
              max={setting.max}
              precision={0}
              value={value}
              disabled={!canEdit}
              onChange={(next) => setDraft(typeof next === "number" ? next : null)}
              addonAfter={t("days")}
            />
            {canEdit ? (
              <Button
                type="primary"
                loading={update.isPending}
                disabled={draft === null || draft === saved}
                onClick={() =>
                  update.mutate(
                    { key: setting.key, value: String(value) },
                    { onSuccess: () => setDraft(null) },
                  )
                }
              >
                {t("save")}
              </Button>
            ) : null}
          </Space>
        </Form.Item>
      </Form>
      {setting.updated_at ? (
        <Typography.Text type="secondary">
          {t("lastChanged")}: {formatDateTime(setting.updated_at, locale)}
        </Typography.Text>
      ) : null}
    </Card>
  );
}
