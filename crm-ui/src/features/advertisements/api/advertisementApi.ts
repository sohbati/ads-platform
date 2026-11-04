import { apiClient } from "@/lib/api/client";
import type {
  AdvertisementDetail,
  AdvertisementPage,
  AdvertisementFiltersState,
  RejectInput,
  ReviewerStats,
} from "../types/advertisement.types";

export async function getAdvertisements(
  params: AdvertisementFiltersState,
): Promise<AdvertisementPage> {
  const { data } = await apiClient.get<AdvertisementPage>("/panel/ads", { params });
  return data;
}

export async function getAdvertisement(id: number): Promise<AdvertisementDetail> {
  const { data } = await apiClient.get<AdvertisementDetail>(`/panel/ads/${id}`);
  return data;
}

export async function approveAdvertisement(id: number): Promise<AdvertisementDetail> {
  const { data } = await apiClient.post<AdvertisementDetail>(`/panel/ads/${id}/approve`);
  return data;
}

export async function rejectAdvertisement(
  id: number,
  input: RejectInput,
): Promise<AdvertisementDetail> {
  const { data } = await apiClient.post<AdvertisementDetail>(
    `/panel/ads/${id}/reject`,
    input,
  );
  return data;
}

/** Claims the next waiting ad, or returns the one the reviewer already holds; null when none wait. */
export async function getNextAdvertisement(): Promise<AdvertisementDetail | null> {
  const { data } = await apiClient.post<AdvertisementDetail | "">("/panel/ads/review/next");
  return data || null;
}

export async function claimAdvertisement(id: number): Promise<AdvertisementDetail> {
  const { data } = await apiClient.post<AdvertisementDetail>(`/panel/ads/${id}/claim`);
  return data;
}

export async function extendClaim(id: number): Promise<AdvertisementDetail> {
  const { data } = await apiClient.post<AdvertisementDetail>(`/panel/ads/${id}/extend`);
  return data;
}

export async function releaseAdvertisement(id: number): Promise<void> {
  await apiClient.post(`/panel/ads/${id}/release`);
}

export async function forceReleaseAdvertisement(id: number): Promise<void> {
  await apiClient.post(`/panel/ads/${id}/force-release`);
}

export async function getReviewStats(): Promise<ReviewerStats[]> {
  const { data } = await apiClient.get<{ items: ReviewerStats[] }>("/panel/ads/review/stats");
  return data.items;
}
