"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { App } from "antd";
import { queryKeys } from "@/lib/query/keys";
import { toApiError, userMessageForError } from "@/lib/api/errors";
import { useI18n } from "@/lib/i18n/I18nProvider";
import {
  approveAdvertisement,
  getAdvertisement,
  getAdvertisements,
  rejectAdvertisement,
} from "../api/advertisementApi";
import type {
  AdvertisementFiltersState,
  RejectInput,
} from "../types/advertisement.types";

export function useAdvertisements(filters: AdvertisementFiltersState) {
  return useQuery({
    queryKey: queryKeys.advertisements(filters),
    queryFn: () => getAdvertisements(filters),
  });
}

export function useAdvertisement(id: number | null) {
  return useQuery({
    queryKey: queryKeys.advertisement(id ?? 0),
    queryFn: () => getAdvertisement(id ?? 0),
    enabled: id !== null,
  });
}

export function useReviewMutations() {
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const { t, locale } = useI18n();

  // A decision changes the ad, its list position, and the review queue, so refresh all of them.
  const refresh = () =>
    queryClient.invalidateQueries({ queryKey: queryKeys.advertisementsAll });
  const failure = (error: unknown) => {
    message.error(userMessageForError(toApiError(error, t("error")), locale));
    // Another reviewer may have decided first; show the current state.
    void refresh();
  };

  const approve = useMutation({
    mutationFn: (id: number) => approveAdvertisement(id),
    onSuccess: async () => {
      await refresh();
      message.success(t("successSave"));
    },
    onError: failure,
  });

  const reject = useMutation({
    mutationFn: ({ id, input }: { id: number; input: RejectInput }) =>
      rejectAdvertisement(id, input),
    onSuccess: async () => {
      await refresh();
      message.success(t("successSave"));
    },
    onError: failure,
  });

  return { approve, reject };
}
