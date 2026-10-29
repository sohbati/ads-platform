"use client";

import { Permission } from "@/lib/permissions/codes";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useGroupMutations, useGroups } from "../hooks/useAccess";
import { NamedEntityScreen } from "./NamedEntityScreen";

export function GroupsScreen() {
  const { t } = useI18n();
  const query = useGroups();
  const mutations = useGroupMutations();

  return (
    <NamedEntityScreen
      title={t("groups")}
      items={query.data}
      loading={query.isLoading}
      error={query.error}
      writePermission={Permission.staffGroupUpdate}
      creating={mutations.create.isPending}
      deleting={mutations.remove.isPending}
      onRetry={() => void query.refetch()}
      onCreate={(input) => mutations.create.mutate(input)}
      onDelete={(id) => mutations.remove.mutate(id)}
    />
  );
}
