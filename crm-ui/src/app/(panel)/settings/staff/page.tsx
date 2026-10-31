import { Suspense } from "react";
import { StaffScreen } from "@/features/admin-users/components/StaffScreen";
import { AppLoading } from "@/components/ui/AppLoading";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function AdminUsersPage() {
  await requirePermission(Permission.staffRead);
  return (
    <Suspense fallback={<AppLoading />}>
      <StaffScreen />
    </Suspense>
  );
}
