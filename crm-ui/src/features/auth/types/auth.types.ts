import type { PermissionCode } from "@/lib/permissions/codes";

export interface SessionStaff {
  staffId: number;
  name: string;
  mobile: string;
  permissions: PermissionCode[];
}

export interface LoginInput {
  mobile: string;
  password: string;
}
