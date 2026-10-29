import { Suspense } from "react";
import { AdvertisementScreen } from "@/features/advertisements/components/AdvertisementScreen";
import { AppLoading } from "@/components/ui/AppLoading";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function AdvertisementsPage() {
  await requirePermission(Permission.advertisementRead);
  return (
    <Suspense fallback={<AppLoading />}>
      <AdvertisementScreen />
    </Suspense>
  );
}
