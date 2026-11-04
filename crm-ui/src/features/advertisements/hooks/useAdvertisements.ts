"use client";

import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { App } from "antd";
import { queryKeys } from "@/lib/query/keys";
import { toApiError, userMessageForError } from "@/lib/api/errors";
import { useI18n } from "@/lib/i18n/I18nProvider";
import {
  approveAdvertisement,
  claimAdvertisement,
  extendClaim,
  forceReleaseAdvertisement,
  getAdvertisement,
  getAdvertisements,
  getNextAdvertisement,
  getReviewStats,
  rejectAdvertisement,
  releaseAdvertisement,
} from "../api/advertisementApi";
import type {
  AdvertisementDetail,
  AdvertisementFiltersState,
  RejectInput,
} from "../types/advertisement.types";

export function useAdvertisements(
  filters: AdvertisementFiltersState,
  options?: { enabled?: boolean; refetchInterval?: number },
) {
  return useQuery({
    queryKey: queryKeys.advertisements(filters),
    queryFn: () => getAdvertisements(filters),
    enabled: options?.enabled ?? true,
    refetchInterval: options?.refetchInterval,
    // Keep the current page on screen while the next one loads.
    placeholderData: keepPreviousData,
  });
}

export function useReviewStats(enabled = true) {
  return useQuery({
    queryKey: queryKeys.reviewStats,
    queryFn: getReviewStats,
    enabled,
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

/** Taking, keeping, and handing back ads in the review queue. */
export function useClaimActions() {
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const { t, locale } = useI18n();

  // Claims show in the lists, so refresh them; the ad itself comes back in the response.
  const show = (ad: AdvertisementDetail | null) => {
    if (ad) {
      queryClient.setQueryData(queryKeys.advertisement(ad.id), ad);
    }
    void queryClient.invalidateQueries({ queryKey: queryKeys.advertisementLists });
  };
  const failure = (error: unknown) => {
    message.error(userMessageForError(toApiError(error, t("error")), locale));
    void queryClient.invalidateQueries({ queryKey: queryKeys.advertisementsAll });
  };

  const next = useMutation({
    mutationFn: getNextAdvertisement,
    onSuccess: show,
    onError: failure,
  });
  const claim = useMutation({
    mutationFn: claimAdvertisement,
    onSuccess: show,
    onError: failure,
  });
  // The heartbeat fails quietly: the drawer notices the claim is gone and says so.
  const extend = useMutation({
    mutationFn: extendClaim,
    onSuccess: (ad) => queryClient.setQueryData(queryKeys.advertisement(ad.id), ad),
    onError: (_, id) =>
      void queryClient.invalidateQueries({ queryKey: queryKeys.advertisement(id) }),
  });
  const release = useMutation({
    mutationFn: releaseAdvertisement,
    onSuccess: () => void queryClient.invalidateQueries({ queryKey: queryKeys.advertisementsAll }),
  });
  const forceRelease = useMutation({
    mutationFn: forceReleaseAdvertisement,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.advertisementsAll });
      message.success(t("successSave"));
    },
    onError: failure,
  });

  return { next, claim, extend, release, forceRelease };
}
