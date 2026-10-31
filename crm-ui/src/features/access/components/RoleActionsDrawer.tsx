"use client";

import { useState } from "react";
import { Button, Checkbox, Drawer, Space, Typography } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useActions, useRoleActions } from "../hooks/useAccess";
import type { IamRole } from "../types/access.types";

interface RoleActionsDrawerProps {
  role: IamRole | null;
  canEdit: boolean;
  onClose: () => void;
}

export function RoleActionsDrawer({ role, canEdit, onClose }: RoleActionsDrawerProps) {
  const { t } = useI18n();
  const allActions = useActions();
  const { query, save } = useRoleActions(role?.id ?? null);
  // null = not edited yet, so the saved actions show; reset when the drawer closes.
  const [draft, setDraft] = useState<number[] | null>(null);
  const selected = draft ?? (query.data ?? []).map((action) => action.id);
  const close = () => {
    setDraft(null);
    onClose();
  };

  return (
    <Drawer
      open={role !== null}
      title={role ? `${t("roleActions")}: ${role.name}` : undefined}
      size="large"
      onClose={close}
      extra={
        canEdit ? (
          <Space>
            <Button onClick={close}>{t("cancel")}</Button>
            <Button
              type="primary"
              loading={save.isPending}
              disabled={query.isLoading}
              onClick={() => save.mutate(selected, { onSuccess: close })}
            >
              {t("save")}
            </Button>
          </Space>
        ) : null
      }
    >
      <Typography.Paragraph type="secondary">{t("roleActionsHint")}</Typography.Paragraph>
      <Checkbox.Group
        className="app-checkbox-list"
        value={selected}
        disabled={!canEdit || query.isLoading || allActions.isLoading}
        onChange={(values) => setDraft(values as number[])}
        options={(allActions.data ?? []).map((action) => ({
          value: action.id,
          label: (
            <span>
              {action.name}{" "}
              <Typography.Text type="secondary" dir="ltr">
                ({action.code})
              </Typography.Text>
            </span>
          ),
        }))}
      />
    </Drawer>
  );
}
