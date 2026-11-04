"use client";

import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/query/keys";
import { getCities } from "../api/getCities";

export function useCities() {
  return useQuery({
    queryKey: queryKeys.cities,
    queryFn: getCities,
    // The city list changes rarely; load it once per session.
    staleTime: Infinity,
  });
}
