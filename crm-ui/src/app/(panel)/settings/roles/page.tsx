import { RolesScreen } from "@/features/access/components/RolesScreen";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function RolesPage() {
  await requirePermission(Permission.roleRead);
  return <RolesScreen />;
}
