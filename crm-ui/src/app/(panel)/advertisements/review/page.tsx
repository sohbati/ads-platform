import { Suspense } from "react";
import { AdvertisementScreen } from "@/features/advertisements/components/AdvertisementScreen";
import { AppLoading } from "@/components/ui/AppLoading";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";
import { AdvertisementStatus } from "@/features/advertisements/constants/advertisement.constants";

export default async function AdvertisementReviewPage() {
  await requirePermission(Permission.advertisementRead);
  return (
    <Suspense fallback={<AppLoading />}>
      <AdvertisementScreen titleKey="review" presetStatus={AdvertisementStatus.pending} />
    </Suspense>
  );
}
