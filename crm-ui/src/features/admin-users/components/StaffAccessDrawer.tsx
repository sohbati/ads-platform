"use client";

import { Drawer, Select, Typography } from "antd";
import { useGroups, useRoles } from "@/features/access/hooks/useAccess";
import type { NamedEntity } from "@/features/access/types/access.types";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useStaffAccess } from "../hooks/useStaff";
import type { StaffListItem } from "../types/staff.types";

interface StaffAccessDrawerProps {
  staff: StaffListItem | null;
  canEdit: boolean;
  onClose: () => void;
}

function toOptions(items: NamedEntity[] | undefined) {
  return (items ?? []).map((item) => ({
    value: item.id,
    label: `${item.name} (${item.code})`,
  }));
}

export function StaffAccessDrawer({ staff, canEdit, onClose }: StaffAccessDrawerProps) {
  const { t } = useI18n();
  const allRoles = useRoles();
  const allGroups = useGroups();
  const access = useStaffAccess(staff?.id ?? null);

  return (
    <Drawer
      open={staff !== null}
      title={staff ? `${t("staffAccessTitle")} ${staff.name}` : undefined}
      size="large"
      onClose={onClose}
    >
      <Typography.Title level={5}>{t("roles")}</Typography.Title>
      <Typography.Paragraph type="secondary">{t("staffRolesHint")}</Typography.Paragraph>
      <Select
        mode="multiple"
        style={{ width: "100%" }}
        placeholder={t("selectRoles")}
        disabled={!canEdit}
        loading={
          allRoles.isLoading || access.roles.isLoading || access.toggleRole.isPending
        }
        options={toOptions(allRoles.data)}
        value={(access.roles.data ?? []).map((role) => role.id)}
        optionFilterProp="label"
        onSelect={(roleId: number) =>
          access.toggleRole.mutate({ roleId, assigned: true })
        }
        onDeselect={(roleId: number) =>
          access.toggleRole.mutate({ roleId, assigned: false })
        }
      />

      <Typography.Title level={5} style={{ marginTop: 32 }}>
        {t("groups")}
      </Typography.Title>
      <Typography.Paragraph type="secondary">{t("staffGroupsHint")}</Typography.Paragraph>
      <Select
        mode="multiple"
        style={{ width: "100%" }}
        placeholder={t("selectGroups")}
        disabled={!canEdit}
        loading={
          allGroups.isLoading || access.groups.isLoading || access.toggleGroup.isPending
        }
        options={toOptions(allGroups.data)}
        value={(access.groups.data ?? []).map((group) => group.id)}
        optionFilterProp="label"
        onSelect={(groupId: number) =>
          access.toggleGroup.mutate({ groupId, assigned: true })
        }
        onDeselect={(groupId: number) =>
          access.toggleGroup.mutate({ groupId, assigned: false })
        }
      />
    </Drawer>
  );
}
