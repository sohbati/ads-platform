"use client";

import { Input } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";

interface SearchInputProps {
  value?: string;
  onChange: (value: string) => void;
}

export function SearchInput({ value, onChange }: SearchInputProps) {
  const { t } = useI18n();
  return (
    <Input.Search
      allowClear
      value={value}
      placeholder={t("search")}
      onChange={(event) => onChange(event.target.value)}
      onSearch={onChange}
      style={{ maxWidth: 320 }}
    />
  );
}
