import { SettingsScreen } from "@/features/settings/components/SettingsScreen";
import { requirePermission } from "@/lib/auth/guards";
import { Permission } from "@/lib/permissions/codes";

export default async function SettingsPage() {
  await requirePermission(Permission.settingsRead);
  return <SettingsScreen />;
}
