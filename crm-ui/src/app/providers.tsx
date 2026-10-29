"use client";

import { useState } from "react";
import { App, ConfigProvider, theme } from "antd";
import faIR from "antd/locale/fa_IR";
import enUS from "antd/locale/en_US";
import { QueryClientProvider } from "@tanstack/react-query";
import { I18nProvider } from "@/lib/i18n/I18nProvider";
import { createQueryClient } from "@/lib/query/client";
import type { Locale } from "@/types/locale.types";
import { directionForLocale } from "@/types/locale.types";

export function Providers({
  locale,
  children,
}: {
  locale: Locale;
  children: React.ReactNode;
}) {
  const [queryClient] = useState(createQueryClient);
  const direction = directionForLocale(locale);

  return (
    <QueryClientProvider client={queryClient}>
      <I18nProvider locale={locale}>
        <ConfigProvider
          direction={direction}
          locale={locale === "fa" ? faIR : enUS}
          theme={{
            algorithm: theme.defaultAlgorithm,
            token: {
              colorPrimary: "#0f766e",
              borderRadius: 8,
              fontFamily: "inherit",
            },
          }}
        >
          <App>{children}</App>
        </ConfigProvider>
      </I18nProvider>
    </QueryClientProvider>
  );
}
