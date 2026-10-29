import { ActionsScreen } from "@/features/access/components/ActionsScreen";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function ActionsPage() {
  await requirePermission(Permission.actionRead);
  return <ActionsScreen />;
}
