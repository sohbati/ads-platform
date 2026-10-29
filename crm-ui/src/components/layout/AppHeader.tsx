"use client";

import { LogoutOutlined } from "@ant-design/icons";
import { Button, Layout, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { logout } from "@/features/auth/api/login";
import { useSession } from "@/lib/auth/SessionProvider";
import { LocaleSwitcher } from "./LocaleSwitcher";

export function AppHeader() {
  const { t } = useI18n();
  const session = useSession();
  const router = useRouter();
  const queryClient = useQueryClient();

  return (
    <Layout.Header className="app-header">
      <Typography.Text className="app-header__name">
        {session.name}
      </Typography.Text>
      <Space>
        <LocaleSwitcher />
        <Button
          icon={<LogoutOutlined />}
          onClick={async () => {
            await logout();
            queryClient.clear();
            router.replace("/login");
            router.refresh();
          }}
        >
          {t("signOut")}
        </Button>
      </Space>
    </Layout.Header>
  );
}
