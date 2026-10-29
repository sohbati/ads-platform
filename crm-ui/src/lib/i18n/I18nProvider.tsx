"use client";

import { createContext, useContext, useMemo } from "react";
import type { Locale } from "@/types/locale.types";
import { directionForLocale } from "@/types/locale.types";
import { messages, type MessageKey } from "./messages";

interface I18nContextValue {
  locale: Locale;
  dir: "rtl" | "ltr";
  t: (key: MessageKey) => string;
}

const I18nContext = createContext<I18nContextValue | null>(null);

export function I18nProvider({
  locale,
  children,
}: {
  locale: Locale;
  children: React.ReactNode;
}) {
  const value = useMemo<I18nContextValue>(
    () => ({
      locale,
      dir: directionForLocale(locale),
      t: (key) => messages[locale][key],
    }),
    [locale],
  );

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n(): I18nContextValue {
  const ctx = useContext(I18nContext);
  if (!ctx) {
    throw new Error("useI18n must be used within I18nProvider");
  }
  return ctx;
}
