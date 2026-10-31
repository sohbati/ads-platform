import { GeneralSettingsScreen } from "@/features/settings/components/GeneralSettingsScreen";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function GeneralSettingsPage() {
  await requirePermission(Permission.settingsRead);
  return <GeneralSettingsScreen />;
}
