import type { MessageKey } from "@/lib/i18n/messages";

// A rejected ad goes back to "draft" with a reason; "rejected" is kept for older data only.
export const AdvertisementStatus = {
  draft: "draft",
  pending: "pending",
  active: "active",
  rejected: "rejected",
  expired: "expired",
  deleted: "deleted",
} as const;

export type AdvertisementStatus =
  (typeof AdvertisementStatus)[keyof typeof AdvertisementStatus];

/** Statuses offered in the list filter. */
export const FILTERABLE_STATUSES: AdvertisementStatus[] = [
  AdvertisementStatus.pending,
  AdvertisementStatus.active,
  AdvertisementStatus.draft,
  AdvertisementStatus.expired,
];

export const STATUS_LABEL: Record<AdvertisementStatus, MessageKey> = {
  draft: "statusDraft",
  pending: "statusPending",
  active: "statusActive",
  rejected: "statusRejected",
  expired: "statusExpired",
  deleted: "statusDeleted",
};

export const REASON_MAX_LENGTH = 1000;
