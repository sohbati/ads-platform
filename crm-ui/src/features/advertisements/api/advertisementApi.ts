import { apiClient } from "@/lib/api/client";
import type { PagedResult } from "@/lib/api/types";
import type {
  AdvertisementDetail,
  AdvertisementFiltersState,
  AdvertisementListItem,
  RejectInput,
} from "../types/advertisement.types";

export async function getAdvertisements(
  params: AdvertisementFiltersState,
): Promise<PagedResult<AdvertisementListItem>> {
  const { data } = await apiClient.get<PagedResult<AdvertisementListItem>>("/panel/ads", {
    params,
  });
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
