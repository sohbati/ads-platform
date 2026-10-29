import { GroupsScreen } from "@/features/access/components/GroupsScreen";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function GroupsPage() {
  await requirePermission(Permission.staffGroupRead);
  return <GroupsScreen />;
}
