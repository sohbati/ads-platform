import { Suspense } from "react";
import { CatalogScreen } from "@/features/catalog/components/CatalogScreen";
import { AppLoading } from "@/components/ui/AppLoading";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";
import { queryKeys } from "@/lib/query/keys";

export default async function NotificationsPage() {
  await requirePermission(Permission.notificationRead);
  return (
    <Suspense fallback={<AppLoading />}>
      <CatalogScreen titleKey="notifications" queryKey={queryKeys.notifications({})} />
    </Suspense>
  );
}
