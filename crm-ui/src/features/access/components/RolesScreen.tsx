"use client";

import { useState } from "react";
import { Button } from "antd";
import { can } from "@/lib/auth/permissions";
import { useSession } from "@/lib/auth/SessionProvider";
import { Permission } from "@/lib/permissions/codes";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useRoleMutations, useRoles } from "../hooks/useAccess";
import type { IamRole } from "../types/access.types";
import { NamedEntityScreen } from "./NamedEntityScreen";
import { RoleActionsDrawer } from "./RoleActionsDrawer";

export function RolesScreen() {
  const { t } = useI18n();
  const session = useSession();
  const query = useRoles();
  const mutations = useRoleMutations();
  const [actionsFor, setActionsFor] = useState<IamRole | null>(null);

  return (
    <>
      <NamedEntityScreen
        title={t("roles")}
        items={query.data}
        loading={query.isLoading}
        error={query.error}
        writePermission={Permission.roleUpdate}
        saving={mutations.create.isPending || mutations.update.isPending}
        deleting={mutations.remove.isPending}
        onRetry={() => void query.refetch()}
        onCreate={(input, done) => mutations.create.mutate(input, { onSuccess: done })}
        onUpdate={(id, input, done) =>
          mutations.update.mutate({ id, input }, { onSuccess: done })
        }
        codeEditable
        onDelete={(id) => mutations.remove.mutate(id)}
        rowActions={(role) => (
          <Button size="small" onClick={() => setActionsFor(role)}>
            {t("actions")}
          </Button>
        )}
      />
      <RoleActionsDrawer
        role={actionsFor}
        canEdit={can(session.permissions, Permission.roleUpdate)}
        onClose={() => setActionsFor(null)}
      />
    </>
  );
}
