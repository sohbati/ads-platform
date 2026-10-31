"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { App } from "antd";
import { queryKeys } from "@/lib/query/keys";
import { toApiError, userMessageForError } from "@/lib/api/errors";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { getSettings, updateSetting } from "../api/settingsApi";

export function useSettings() {
  return useQuery({ queryKey: queryKeys.settingsList, queryFn: getSettings });
}

export function useUpdateSetting() {
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const { t, locale } = useI18n();
  return useMutation({
    mutationFn: ({ key, value }: { key: string; value: string }) =>
      updateSetting(key, value),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.settingsList });
      message.success(t("successSave"));
    },
    onError: (error) =>
      message.error(userMessageForError(toApiError(error, t("error")), locale)),
  });
}
