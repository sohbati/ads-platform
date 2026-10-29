"use client";

import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/query/keys";
import { getCustomers } from "../api/getCustomers";
import type { CustomerFiltersState } from "../types/customer.types";

export function useCustomers(filters: CustomerFiltersState) {
  return useQuery({
    queryKey: queryKeys.customers(filters),
    queryFn: () => getCustomers(filters),
  });
}
