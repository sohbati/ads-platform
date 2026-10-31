"use client";

import { useState } from "react";
import { Layout, Menu } from "antd";
import Image from "next/image";
import { usePathname, useRouter } from "next/navigation";
import { navigation, type NavItem } from "@/config/navigation";
import { can } from "@/lib/auth/permissions";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useSession } from "@/lib/auth/SessionProvider";

function visibleItems(
  items: NavItem[],
  permissions: readonly string[] | undefined,
): NavItem[] {
  return items
    .filter((item) => !item.permission || can(permissions, item.permission))
    .map((item) => ({
      ...item,
      children: item.children ? visibleItems(item.children, permissions) : undefined,
    }));
}

export function AppSider() {
  const { t } = useI18n();
  const pathname = usePathname();
  const router = useRouter();
  const session = useSession();
  const [collapsed, setCollapsed] = useState(false);
  const items = visibleItems(navigation, session.permissions);

  const selected = pathname.startsWith("/advertisements/review")
    ? "advertisements-review"
    : pathname.startsWith("/settings/general")
      ? "settings-general"
      : pathname.startsWith("/settings/staff")
        ? "settings-staff"
        : pathname.startsWith("/settings/roles")
          ? "settings-roles"
          : pathname.startsWith("/settings/actions")
            ? "settings-actions"
            : pathname.startsWith("/settings/groups")
              ? "settings-groups"
              : (items.find(
                  (item) =>
                    pathname === item.href || pathname.startsWith(`${item.href}/`),
                )?.key ?? "dashboard");

  return (
    <Layout.Sider
      breakpoint="lg"
      collapsedWidth={64}
      width={248}
      className="app-sider"
      onCollapse={setCollapsed}
    >
      <div className="app-brand">
        <Image src="/logo-mark.png" alt={t("appName")} width={32} height={32} />
        {collapsed ? null : <span>{t("appName")}</span>}
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={[selected]}
        defaultOpenKeys={
          pathname.startsWith("/settings")
            ? ["settings"]
            : pathname.startsWith("/advertisements")
              ? ["advertisements"]
              : []
        }
        items={items.map((item) => ({
          key: item.key,
          icon: <item.icon />,
          label: t(item.labelKey),
          children: item.children?.map((child) => ({
            key: child.key,
            icon: <child.icon />,
            label: t(child.labelKey),
            onClick: () => router.push(child.href),
          })),
          onClick: item.children ? undefined : () => router.push(item.href),
        }))}
      />
    </Layout.Sider>
  );
}
