export interface AdminSession {
  staffId: number;
  name: string;
  mobile: string;
  isSuperAdmin: boolean;
  /** Action codes from the staff member's roles, or `["*"]` for a super admin. */
  permissions: string[];
}

/** Response of the backend `/api/v1/panel/auth/*` endpoints. */
export interface BackendStaffSession {
  staff: {
    id: number;
    name: string;
    mobile: string;
    is_super_admin: boolean;
  };
  permissions: string[];
}

/** What the signed cookie carries. Permissions are always reloaded from the backend. */
export interface SessionToken {
  sid: number;
  exp: number;
}
