import type { Metadata } from "next";
import { cookies } from "next/headers";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import { Providers } from "./providers";
import { localeMeta, resolveLocale, LOCALE_COOKIE } from "@/lib/i18n/locale";
import "./globals.css";

export const metadata: Metadata = {
  title: "Ruab Admin",
  description: "Ruab classifieds administration",
};

export default async function RootLayout({ children }: LayoutProps<"/">) {
  const store = await cookies();
  const locale = resolveLocale(store.get(LOCALE_COOKIE)?.value);
  const { lang, dir } = localeMeta(locale);

  return (
    <html lang={lang} dir={dir}>
      <body
        style={{
          fontFamily:
            locale === "fa"
              ? '"Vazirmatn Variable", sans-serif'
              : '"Inter Variable", sans-serif',
        }}
      >
        <AntdRegistry>
          <Providers locale={locale}>{children}</Providers>
        </AntdRegistry>
      </body>
    </html>
  );
}
