"use client";

import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/query/keys";
import { getCurrentStaff } from "../api/login";

export function useCurrentStaff() {
  return useQuery({
    queryKey: queryKeys.me,
    queryFn: getCurrentStaff,
    retry: false,
  });
}
