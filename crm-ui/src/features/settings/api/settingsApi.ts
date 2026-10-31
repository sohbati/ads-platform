import { apiClient } from "@/lib/api/client";
import type { AppSetting } from "../types/settings.types";

export async function getSettings(): Promise<AppSetting[]> {
  const { data } = await apiClient.get<AppSetting[]>("/panel/settings");
  return data;
}

export async function updateSetting(key: string, value: string): Promise<AppSetting> {
  const { data } = await apiClient.put<AppSetting>(`/panel/settings/${key}`, { value });
  return data;
}
