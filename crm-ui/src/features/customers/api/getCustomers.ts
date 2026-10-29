import { apiClient } from "@/lib/api/client";
import type { PagedResult } from "@/lib/api/types";
import type { CustomerFiltersState, CustomerListItem } from "../types/customer.types";

export async function getCustomers(
  params: CustomerFiltersState,
): Promise<PagedResult<CustomerListItem>> {
  const { data } = await apiClient.get<PagedResult<CustomerListItem>>("/catalog", {
    params,
  });
  return data;
}
