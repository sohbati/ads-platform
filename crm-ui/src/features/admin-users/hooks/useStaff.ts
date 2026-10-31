"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { App } from "antd";
import { queryKeys } from "@/lib/query/keys";
import { toApiError, userMessageForError } from "@/lib/api/errors";
import { useI18n } from "@/lib/i18n/I18nProvider";
import {
  addStaffGroup,
  addStaffRole,
  createStaff,
  deleteStaff,
  getStaff,
  getStaffGroups,
  getStaffRoles,
  removeStaffGroup,
  removeStaffRole,
  setStaffPassword,
  updateStaff,
} from "../api/staffApi";
import type { StaffFiltersState, StaffInput } from "../types/staff.types";

export function useStaff(filters: StaffFiltersState) {
  return useQuery({
    queryKey: queryKeys.staff(filters),
    queryFn: () => getStaff(filters),
  });
}

function useFeedback() {
  const { message } = App.useApp();
  const { t, locale } = useI18n();
  return {
    t,
    success: (key: "successSave" | "successDelete") => message.success(t(key)),
    failure: (error: unknown) =>
      message.error(userMessageForError(toApiError(error, t("error")), locale)),
  };
}

export function useStaffMutations() {
  const queryClient = useQueryClient();
  const feedback = useFeedback();
  const refresh = () => queryClient.invalidateQueries({ queryKey: queryKeys.staffAll });

  const create = useMutation({
    mutationFn: (input: StaffInput) => createStaff(input),
    onSuccess: async () => {
      await refresh();
      feedback.success("successSave");
    },
    onError: feedback.failure,
  });

  const update = useMutation({
    mutationFn: ({ id, input }: { id: number; input: StaffInput }) =>
      updateStaff(id, input),
    onSuccess: async () => {
      await refresh();
      feedback.success("successSave");
    },
    onError: feedback.failure,
  });

  const remove = useMutation({
    mutationFn: (id: number) => deleteStaff(id),
    onSuccess: async () => {
      await refresh();
      feedback.success("successDelete");
    },
    onError: feedback.failure,
  });

  const setPassword = useMutation({
    mutationFn: ({ id, password }: { id: number; password: string }) =>
      setStaffPassword(id, password),
    onSuccess: () => feedback.success("successSave"),
    onError: feedback.failure,
  });

  return { create, update, remove, setPassword };
}

export function useStaffAccess(staffId: number | null) {
  const queryClient = useQueryClient();
  const feedback = useFeedback();
  const enabled = staffId !== null;
  const id = staffId ?? 0;

  const roles = useQuery({
    queryKey: queryKeys.staffRoles(id),
    queryFn: () => getStaffRoles(id),
    enabled,
  });
  const groups = useQuery({
    queryKey: queryKeys.staffGroups(id),
    queryFn: () => getStaffGroups(id),
    enabled,
  });

  const toggleRole = useMutation({
    mutationFn: ({ roleId, assigned }: { roleId: number; assigned: boolean }) =>
      assigned ? addStaffRole(id, roleId) : removeStaffRole(id, roleId),
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: queryKeys.staffRoles(id) }),
    onError: feedback.failure,
  });

  const toggleGroup = useMutation({
    mutationFn: ({ groupId, assigned }: { groupId: number; assigned: boolean }) =>
      assigned ? addStaffGroup(id, groupId) : removeStaffGroup(id, groupId),
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: queryKeys.staffGroups(id) }),
    onError: feedback.failure,
  });

  return { roles, groups, toggleRole, toggleGroup };
}
