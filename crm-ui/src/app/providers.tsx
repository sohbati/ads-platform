"use client";

import { ConfigProvider } from "antd";

export function Providers({ children }: { children: React.ReactNode }) {
  return <ConfigProvider>{children}</ConfigProvider>;
}
