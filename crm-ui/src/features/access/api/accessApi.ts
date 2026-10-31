import { apiClient } from "@/lib/api/client";
import type {
  IamAction,
  IamGroup,
  IamRole,
  NamedEntityInput,
} from "../types/access.types";

export async function getActions(): Promise<IamAction[]> {
  const { data } = await apiClient.get<IamAction[]>("/panel/actions");
  return data;
}

export async function createAction(input: NamedEntityInput): Promise<IamAction> {
  const { data } = await apiClient.post<IamAction>("/panel/actions", input);
  return data;
}

export async function deleteAction(id: number): Promise<void> {
  await apiClient.delete(`/panel/actions/${id}`);
}

export async function getRoles(): Promise<IamRole[]> {
  const { data } = await apiClient.get<IamRole[]>("/panel/roles");
  return data;
}

export async function createRole(input: NamedEntityInput): Promise<IamRole> {
  const { data } = await apiClient.post<IamRole>("/panel/roles", input);
  return data;
}

export async function deleteRole(id: number): Promise<void> {
  await apiClient.delete(`/panel/roles/${id}`);
}

export async function getGroups(): Promise<IamGroup[]> {
  const { data } = await apiClient.get<IamGroup[]>("/panel/staff-groups");
  return data;
}

export async function createGroup(input: NamedEntityInput): Promise<IamGroup> {
  const { data } = await apiClient.post<IamGroup>("/panel/staff-groups", input);
  return data;
}

export async function deleteGroup(id: number): Promise<void> {
  await apiClient.delete(`/panel/staff-groups/${id}`);
}

// Action codes are matched by permission checks, so only the name can change.
export async function updateAction(
  id: number,
  input: NamedEntityInput,
): Promise<IamAction> {
  const { data } = await apiClient.put<IamAction>(`/panel/actions/${id}`, {
    name: input.name,
  });
  return data;
}

export async function updateRole(id: number, input: NamedEntityInput): Promise<IamRole> {
  const { data } = await apiClient.put<IamRole>(`/panel/roles/${id}`, input);
  return data;
}

export async function updateGroup(
  id: number,
  input: NamedEntityInput,
): Promise<IamGroup> {
  const { data } = await apiClient.put<IamGroup>(`/panel/staff-groups/${id}`, input);
  return data;
}

export async function getRoleActions(roleId: number): Promise<IamAction[]> {
  const { data } = await apiClient.get<IamAction[]>(`/panel/roles/${roleId}/actions`);
  return data;
}

export async function setRoleActions(
  roleId: number,
  actionIds: number[],
): Promise<IamAction[]> {
  const { data } = await apiClient.put<IamAction[]>(`/panel/roles/${roleId}/actions`, {
    action_ids: actionIds,
  });
  return data;
}
