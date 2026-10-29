import type { AdvertisementStatus } from "../constants/advertisement.constants";

export interface AdvertisementListItem {
  id: number;
  title: string;
  customerId: number;
  status: AdvertisementStatus;
  cityName?: string;
  createdAt: string;
}

export interface AdvertisementFiltersState {
  page: number;
  pageSize: number;
  search?: string;
  status?: AdvertisementStatus;
}
