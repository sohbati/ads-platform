import { apiClient } from "@/lib/api/client";
import type { Category } from "../types/category.types";

export async function getCategories(): Promise<Category[]> {
  const { data } = await apiClient.get<Category[]>("/categories");
  return data;
}
