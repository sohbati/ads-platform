import { apiClient } from "@/lib/api/client";
import type { PagedResult } from "@/lib/api/types";
import type { StaffFiltersState, StaffListItem } from "../types/staff.types";

export async function getStaff(
  params: StaffFiltersState,
): Promise<PagedResult<StaffListItem>> {
  const { data } = await apiClient.get<PagedResult<StaffListItem>>("/catalog", { params });
  return data;
}
