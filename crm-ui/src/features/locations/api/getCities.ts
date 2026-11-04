import { apiClient } from "@/lib/api/client";
import type { City } from "../types/city.types";

export async function getCities(): Promise<City[]> {
  const { data } = await apiClient.get<City[]>("/cities");
  return data;
}
