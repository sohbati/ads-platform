export interface StaffListItem {
  id: number;
  name: string;
  mobile: string;
  nationalId: string;
  createdAt: string;
}

export interface StaffFiltersState {
  page: number;
  pageSize: number;
  search?: string;
}
