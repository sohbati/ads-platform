import type { PermissionCode } from "@/lib/permissions/codes";

export function can(
  permissions: readonly string[] | undefined,
  permission: PermissionCode,
): boolean {
  if (!permissions || permissions.length === 0) {
    return false;
  }
  return permissions.includes("*") || permissions.includes(permission);
}

export function canAny(
  permissions: readonly string[] | undefined,
  required: readonly PermissionCode[],
): boolean {
  return required.some((permission) => can(permissions, permission));
}
