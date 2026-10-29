"use client";

import { Layout } from "antd";
import { AppHeader } from "./AppHeader";
import { AppSider } from "./AppSider";

export function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <Layout className="app-shell">
      <AppSider />
      <Layout>
        <AppHeader />
        <Layout.Content className="app-content">{children}</Layout.Content>
      </Layout>
    </Layout>
  );
}
