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
      saving={mutations.create.isPending || mutations.update.isPending}
      deleting={mutations.remove.isPending}
      onRetry={() => void query.refetch()}
      onCreate={(input, done) => mutations.create.mutate(input, { onSuccess: done })}
      onUpdate={(id, input, done) =>
        mutations.update.mutate({ id, input }, { onSuccess: done })
      }
      codeEditable
      onDelete={(id) => mutations.remove.mutate(id)}
    />
  );
}
