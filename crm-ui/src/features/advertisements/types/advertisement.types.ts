import type { AdvertisementStatus } from "../constants/advertisement.constants";

export interface AdvertisementListItem {
  id: number;
  title: string;
  status: AdvertisementStatus;
  category_id: number;
  category_title?: string;
  city_id: number;
  city_name?: string;
  price_amount: number | null;
  price_type: string;
  currency: string;
  customer_id: number;
  customer_name: string;
  customer_mobile: string;
  pending_since: string | null;
  published_at: string | null;
  reviewed_at: string | null;
  /** The last review was a rejection: the ad is back in draft, or was resubmitted. */
  was_rejected: boolean;
  created_at: string;
  updated_at: string;
}

export interface AdvertisementMedia {
  url?: string;
  thumb?: string;
  is_cover?: boolean;
}

export interface AdvertisementDetail extends AdvertisementListItem {
  description: string;
  attrs: Record<string, unknown> | null;
  media: AdvertisementMedia[] | null;
  contact: Record<string, unknown> | null;
  location: { lat?: number; lng?: number; neighborhood?: string } | null;
  rejection_reason_customer: string | null;
  rejection_reason_internal: string | null;
  reviewed_by_staff_id: number | null;
  reviewed_by_name: string | null;
}

export interface AdvertisementFiltersState {
  page: number;
  pageSize: number;
  search?: string;
  status?: AdvertisementStatus;
}

export interface RejectInput {
  customer_reason: string;
  internal_reason?: string;
}
