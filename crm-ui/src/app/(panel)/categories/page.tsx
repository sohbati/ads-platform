import { Suspense } from "react";
import { CategoryScreen } from "@/features/categories/components/CategoryScreen";
import { AppLoading } from "@/components/ui/AppLoading";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function CategoriesPage() {
  await requirePermission(Permission.categoryRead);
  return (
    <Suspense fallback={<AppLoading />}>
      <CategoryScreen />
    </Suspense>
  );
}
