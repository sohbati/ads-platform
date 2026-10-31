"use client";

import { Permission } from "@/lib/permissions/codes";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useActionMutations, useActions } from "../hooks/useAccess";
import { NamedEntityScreen } from "./NamedEntityScreen";

export function ActionsScreen() {
  const { t } = useI18n();
  const query = useActions();
  const mutations = useActionMutations();

  return (
    <NamedEntityScreen
      title={t("actions")}
      items={query.data}
      loading={query.isLoading}
      error={query.error}
      writePermission={Permission.actionUpdate}
      saving={mutations.create.isPending || mutations.update.isPending}
      deleting={mutations.remove.isPending}
      onRetry={() => void query.refetch()}
      onCreate={(input, done) => mutations.create.mutate(input, { onSuccess: done })}
      onUpdate={(id, input, done) =>
        mutations.update.mutate({ id, input }, { onSuccess: done })
      }
      onDelete={(id) => mutations.remove.mutate(id)}
    />
  );
}
