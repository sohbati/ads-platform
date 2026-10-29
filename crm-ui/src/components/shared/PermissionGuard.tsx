"use client";

import { can } from "@/lib/auth/permissions";
import type { PermissionCode } from "@/lib/permissions/codes";

interface PermissionGuardProps {
  permission: PermissionCode;
  permissions: readonly string[] | undefined;
  children: React.ReactNode;
}

export function PermissionGuard({
  permission,
  permissions,
  children,
}: PermissionGuardProps) {
  if (!can(permissions, permission)) {
    return null;
  }
  return <>{children}</>;
}
