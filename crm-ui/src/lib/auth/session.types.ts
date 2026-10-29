import type { PermissionCode } from "@/lib/permissions/codes";

export interface AdminSession {
  staffId: number;
  name: string;
  mobile: string;
  permissions: PermissionCode[];
}
