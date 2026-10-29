import { Suspense } from "react";
import { CustomerScreen } from "@/features/customers/components/CustomerScreen";
import { Permission } from "@/lib/permissions/codes";
import { requirePermission } from "@/lib/auth/guards";
import { AppLoading } from "@/components/ui/AppLoading";

export default async function CustomersPage() {
  await requirePermission(Permission.customerRead);
  return (
    <Suspense fallback={<AppLoading />}>
      <CustomerScreen />
    </Suspense>
  );
}
