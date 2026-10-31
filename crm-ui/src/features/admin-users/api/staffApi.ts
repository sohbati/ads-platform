import { apiClient } from "@/lib/api/client";
import type { PagedResult } from "@/lib/api/types";
import type { IamGroup, IamRole } from "@/features/access/types/access.types";
import type { StaffFiltersState, StaffInput, StaffListItem } from "../types/staff.types";

export async function getStaff(
  params: StaffFiltersState,
): Promise<PagedResult<StaffListItem>> {
  const { data } = await apiClient.get<PagedResult<StaffListItem>>("/panel/staff", {
    params,
  });
  return data;
}

export async function createStaff(input: StaffInput): Promise<StaffListItem> {
  const { data } = await apiClient.post<StaffListItem>("/panel/staff", input);
  return data;
}

export async function updateStaff(id: number, input: StaffInput): Promise<StaffListItem> {
  const { data } = await apiClient.put<StaffListItem>(`/panel/staff/${id}`, input);
  return data;
}

export async function deleteStaff(id: number): Promise<void> {
  await apiClient.delete(`/panel/staff/${id}`);
}

export async function setStaffPassword(id: number, password: string): Promise<void> {
  await apiClient.put(`/panel/staff/${id}/password`, { password });
}

export async function getStaffRoles(staffId: number): Promise<IamRole[]> {
  const { data } = await apiClient.get<IamRole[]>(`/panel/staff/${staffId}/roles`);
  return data;
}

export async function addStaffRole(staffId: number, roleId: number): Promise<void> {
  await apiClient.post(`/panel/staff/${staffId}/roles/${roleId}`);
}

export async function removeStaffRole(staffId: number, roleId: number): Promise<void> {
  await apiClient.delete(`/panel/staff/${staffId}/roles/${roleId}`);
}

export async function getStaffGroups(staffId: number): Promise<IamGroup[]> {
  const { data } = await apiClient.get<IamGroup[]>(`/panel/staff/${staffId}/groups`);
  return data;
}

export async function addStaffGroup(staffId: number, groupId: number): Promise<void> {
  await apiClient.post(`/panel/staff-groups/${groupId}/staff/${staffId}`);
}

export async function removeStaffGroup(staffId: number, groupId: number): Promise<void> {
  await apiClient.delete(`/panel/staff-groups/${groupId}/staff/${staffId}`);
}
