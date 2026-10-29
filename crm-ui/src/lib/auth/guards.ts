import { redirect } from "next/navigation";
import { getSession } from "./session";
import { can } from "./permissions";
import type { PermissionCode } from "@/lib/permissions/codes";

export async function requireSession() {
  const session = await getSession();
  if (!session) {
    redirect("/login");
  }
  return session;
}

export async function requirePermission(permission: PermissionCode) {
  const session = await requireSession();
  if (!can(session.permissions, permission)) {
    redirect("/dashboard");
  }
  return session;
}
