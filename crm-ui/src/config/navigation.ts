import {
  ApartmentOutlined,
  AuditOutlined,
  BarChartOutlined,
  BellOutlined,
  DashboardOutlined,
  EnvironmentOutlined,
  FileTextOutlined,
  SafetyCertificateOutlined,
  SettingOutlined,
  TeamOutlined,
  UserOutlined,
} from "@ant-design/icons";
import type { PermissionCode } from "@/lib/permissions/codes";
import { Permission } from "@/lib/permissions/codes";
import type { MessageKey } from "@/lib/i18n/messages";

export interface NavItem {
  key: string;
  href: string;
  labelKey: MessageKey;
  icon: typeof DashboardOutlined;
  permission?: PermissionCode;
  children?: NavItem[];
}

export const navigation: NavItem[] = [
  {
    key: "dashboard",
    href: "/dashboard",
    labelKey: "dashboard",
    icon: DashboardOutlined,
  },
  {
    key: "advertisements",
    href: "/advertisements",
    labelKey: "advertisements",
    icon: FileTextOutlined,
    permission: Permission.advertisementRead,
    children: [
      {
        key: "advertisements-all",
        href: "/advertisements",
        labelKey: "advertisements",
        icon: FileTextOutlined,
        permission: Permission.advertisementRead,
      },
      {
        key: "advertisements-review",
        href: "/advertisements/review",
        labelKey: "review",
        icon: AuditOutlined,
        permission: Permission.advertisementRead,
      },
    ],
  },
  {
    key: "customers",
    href: "/customers",
    labelKey: "customers",
    icon: UserOutlined,
    permission: Permission.customerRead,
  },
  {
    key: "reports",
    href: "/reports",
    labelKey: "reports",
    icon: BarChartOutlined,
    permission: Permission.reportRead,
  },
  {
    key: "categories",
    href: "/categories",
    labelKey: "categories",
    icon: ApartmentOutlined,
    permission: Permission.categoryRead,
  },
  {
    key: "locations",
    href: "/locations",
    labelKey: "locations",
    icon: EnvironmentOutlined,
    permission: Permission.locationRead,
  },
  {
    key: "content",
    href: "/content",
    labelKey: "content",
    icon: FileTextOutlined,
    permission: Permission.contentRead,
  },
  {
    key: "notifications",
    href: "/notifications",
    labelKey: "notifications",
    icon: BellOutlined,
    permission: Permission.notificationRead,
  },
  {
    key: "settings",
    href: "/settings",
    labelKey: "settings",
    icon: SettingOutlined,
    permission: Permission.settingsRead,
    children: [
      {
        key: "settings-staff",
        href: "/settings/staff",
        labelKey: "adminUsers",
        icon: TeamOutlined,
        permission: Permission.staffRead,
      },
      {
        key: "settings-roles",
        href: "/settings/roles",
        labelKey: "roles",
        icon: SafetyCertificateOutlined,
        permission: Permission.roleRead,
      },
      {
        key: "settings-actions",
        href: "/settings/actions",
        labelKey: "actions",
        icon: SafetyCertificateOutlined,
        permission: Permission.actionRead,
      },
      {
        key: "settings-groups",
        href: "/settings/groups",
        labelKey: "groups",
        icon: TeamOutlined,
        permission: Permission.staffGroupRead,
      },
    ],
  },
];
