"use client";

import { Permission } from "@/lib/permissions/codes";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useRoleMutations, useRoles } from "../hooks/useAccess";
import { NamedEntityScreen } from "./NamedEntityScreen";

export function RolesScreen() {
  const { t } = useI18n();
  const query = useRoles();
  const mutations = useRoleMutations();

  return (
    <NamedEntityScreen
      title={t("roles")}
      items={query.data}
      loading={query.isLoading}
      error={query.error}
      writePermission={Permission.roleUpdate}
      creating={mutations.create.isPending}
      deleting={mutations.remove.isPending}
      onRetry={() => void query.refetch()}
      onCreate={(input) => mutations.create.mutate(input)}
      onDelete={(id) => mutations.remove.mutate(id)}
    />
  );
}
