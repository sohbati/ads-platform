export interface AppSetting {
  key: string;
  value: string;
  type: "int";
  min?: number;
  max?: number;
  default: string;
  updated_by_staff_id: number | null;
  updated_at: string | null;
}

/** Keys the backend knows; keep in sync with ads-platform-back setting registry. */
export const SettingKey = {
  adsPendingMaxDays: "ads.pending_max_days",
} as const;
