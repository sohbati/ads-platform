export interface StaffListItem {
  id: number;
  name: string;
  mobile: string;
  national_id: string;
  is_super_admin: boolean;
  is_active: boolean;
  /** One of their roles lets them take ads from the review queue. */
  is_reviewer: boolean;
  created_at: string;
  updated_at: string;
}

export interface StaffFiltersState {
  page: number;
  pageSize: number;
  search?: string;
}

export interface StaffInput {
  name: string;
  mobile: string;
  national_id: string;
  /** Required when creating; ignored on update (use the password endpoint). */
  password?: string;
  is_active?: boolean;
}
