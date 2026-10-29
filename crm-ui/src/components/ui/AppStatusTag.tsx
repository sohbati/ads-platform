"use client";

import { Tag } from "antd";

type StatusTone = "success" | "processing" | "warning" | "error" | "default";

interface AppStatusTagProps {
  label: string;
  tone?: StatusTone;
}

export function AppStatusTag({ label, tone = "default" }: AppStatusTagProps) {
  return <Tag color={tone}>{label}</Tag>;
}
