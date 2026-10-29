import { apiClient } from "@/lib/api/client";
import type { PagedResult } from "@/lib/api/types";
import type {
  AdvertisementFiltersState,
  AdvertisementListItem,
} from "../types/advertisement.types";

export async function getAdvertisements(
  params: AdvertisementFiltersState,
): Promise<PagedResult<AdvertisementListItem>> {
  const { data } = await apiClient.get<PagedResult<AdvertisementListItem>>("/catalog", {
    params,
  });
  return data;
}
