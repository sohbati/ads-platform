import type { CustomerStatus } from "../constants/customer.constants";

export interface CustomerListItem {
  id: number;
  name: string;
  mobile: string;
  nationalId: string;
  status: CustomerStatus;
  createdAt: string;
}

export interface CustomerFiltersState {
  page: number;
  pageSize: number;
  search?: string;
  status?: CustomerStatus;
}
