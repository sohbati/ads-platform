"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { App } from "antd";
import { queryKeys } from "@/lib/query/keys";
import { toApiError, userMessageForError } from "@/lib/api/errors";
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
  getRoleActions,
  getRoles,
  setRoleActions,
  updateAction,
  updateGroup,
  updateRole,
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

function useFeedback() {
  const { message } = App.useApp();
  const { t, locale } = useI18n();
  return {
    success: (key: "successSave" | "successDelete") => message.success(t(key)),
    failure: (error: unknown) =>
      message.error(userMessageForError(toApiError(error, t("error")), locale)),
  };
}

function useNamedMutations(
  key: readonly unknown[],
  createFn: (input: NamedEntityInput) => Promise<unknown>,
  updateFn: (id: number, input: NamedEntityInput) => Promise<unknown>,
  deleteFn: (id: number) => Promise<void>,
) {
  const queryClient = useQueryClient();
  const feedback = useFeedback();
  const refresh = () => queryClient.invalidateQueries({ queryKey: key });

  const create = useMutation({
    mutationFn: createFn,
    onSuccess: async () => {
      await refresh();
      feedback.success("successSave");
    },
    onError: feedback.failure,
  });

  const update = useMutation({
    mutationFn: ({ id, input }: { id: number; input: NamedEntityInput }) =>
      updateFn(id, input),
    onSuccess: async () => {
      await refresh();
      feedback.success("successSave");
    },
    onError: feedback.failure,
  });

  const remove = useMutation({
    mutationFn: deleteFn,
    onSuccess: async () => {
      await refresh();
      feedback.success("successDelete");
    },
    onError: feedback.failure,
  });

  return { create, update, remove };
}

export function useActionMutations() {
  return useNamedMutations(queryKeys.actions, createAction, updateAction, deleteAction);
}

export function useRoleMutations() {
  return useNamedMutations(queryKeys.roles, createRole, updateRole, deleteRole);
}

export function useGroupMutations() {
  return useNamedMutations(queryKeys.groups, createGroup, updateGroup, deleteGroup);
}

export function useRoleActions(roleId: number | null) {
  const queryClient = useQueryClient();
  const feedback = useFeedback();
  const id = roleId ?? 0;

  const query = useQuery({
    queryKey: queryKeys.roleActions(id),
    queryFn: () => getRoleActions(id),
    enabled: roleId !== null,
  });

  const save = useMutation({
    mutationFn: (actionIds: number[]) => setRoleActions(id, actionIds),
    onSuccess: async () => {
      // Staff permissions come from role actions, so cached staff action lists are stale too.
      await queryClient.invalidateQueries({ queryKey: queryKeys.roleActions(id) });
      await queryClient.invalidateQueries({ queryKey: queryKeys.staffAll });
      feedback.success("successSave");
    },
    onError: feedback.failure,
  });

  return { query, save };
}
