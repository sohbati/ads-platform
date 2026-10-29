import { Suspense } from "react";
import { CatalogScreen } from "@/features/catalog/components/CatalogScreen";
import { AppLoading } from "@/components/ui/AppLoading";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";
import { queryKeys } from "@/lib/query/keys";

export default async function CategoriesPage() {
  await requirePermission(Permission.categoryRead);
  return (
    <Suspense fallback={<AppLoading />}>
      <CatalogScreen titleKey="categories" queryKey={queryKeys.categories({})} />
    </Suspense>
  );
}
