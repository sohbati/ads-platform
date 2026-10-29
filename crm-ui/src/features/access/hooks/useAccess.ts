"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { App } from "antd";
import { queryKeys } from "@/lib/query/keys";
import { useI18n } from "@/lib/i18n/I18nProvider";
import {
  createAction,
  createGroup,
  createRole,
  deleteAction,
  deleteGroup,
  deleteRole,
  getActions,
  getGroups,
  getRoles,
} from "../api/accessApi";
import type { NamedEntityInput } from "../types/access.types";

export function useActions() {
  return useQuery({ queryKey: queryKeys.actions, queryFn: getActions });
}

export function useRoles() {
  return useQuery({ queryKey: queryKeys.roles, queryFn: getRoles });
}

export function useGroups() {
  return useQuery({ queryKey: queryKeys.groups, queryFn: getGroups });
}

function useNamedMutations(
  key: readonly unknown[],
  createFn: (input: NamedEntityInput) => Promise<unknown>,
  deleteFn: (id: number) => Promise<void>,
) {
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const { t } = useI18n();

  const create = useMutation({
    mutationFn: createFn,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: key });
      message.success(t("successSave"));
    },
  });

  const remove = useMutation({
    mutationFn: deleteFn,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: key });
      message.success(t("successDelete"));
    },
  });

  return { create, remove };
}

export function useActionMutations() {
  return useNamedMutations(queryKeys.actions, createAction, deleteAction);
}

export function useRoleMutations() {
  return useNamedMutations(queryKeys.roles, createRole, deleteRole);
}

export function useGroupMutations() {
  return useNamedMutations(queryKeys.groups, createGroup, deleteGroup);
}
