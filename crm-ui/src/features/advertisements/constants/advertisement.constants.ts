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
