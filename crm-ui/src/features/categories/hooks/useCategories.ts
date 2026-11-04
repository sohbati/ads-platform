"use client";

import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/query/keys";
import { getCategories } from "../api/getCategories";

export function useCategories() {
  return useQuery({
    queryKey: queryKeys.categories({}),
    queryFn: getCategories,
  });
}
