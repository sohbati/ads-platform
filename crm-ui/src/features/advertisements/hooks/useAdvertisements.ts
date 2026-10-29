"use client";

import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/query/keys";
import { getAdvertisements } from "../api/getAdvertisements";
import type { AdvertisementFiltersState } from "../types/advertisement.types";

export function useAdvertisements(filters: AdvertisementFiltersState) {
  return useQuery({
    queryKey: queryKeys.advertisements(filters),
    queryFn: () => getAdvertisements(filters),
  });
}
