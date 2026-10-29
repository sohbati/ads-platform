"use client";

import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/query/keys";
import { getStaff } from "../api/getStaff";
import type { StaffFiltersState } from "../types/staff.types";

export function useStaff(filters: StaffFiltersState) {
  return useQuery({
    queryKey: queryKeys.staff(filters),
    queryFn: () => getStaff(filters),
  });
}
