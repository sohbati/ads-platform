import { apiClient } from "@/lib/api/client";
import type { LoginInput, SessionStaff } from "../types/auth.types";

export async function login(input: LoginInput): Promise<void> {
  await apiClient.post("/auth/login", input);
}

export async function logout(): Promise<void> {
  await apiClient.post("/auth/logout");
}

export async function getCurrentStaff(): Promise<SessionStaff> {
  const { data } = await apiClient.get<SessionStaff>("/auth/me");
  return data;
}
