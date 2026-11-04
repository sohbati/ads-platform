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

/** How often an open ad renews the reviewer's claim; the backend lease is 10 minutes. */
export const CLAIM_HEARTBEAT_MS = 60_000;
/** Warn the reviewer this long before their claim ends. */
export const CLAIM_WARNING_MS = 60_000;
