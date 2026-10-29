"use client";

import { Alert, Button } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { userMessageForError } from "@/lib/api/errors";
import type { ApiError } from "@/lib/api/types";
import { toApiError } from "@/lib/api/errors";

interface AppErrorProps {
  error: unknown;
  onRetry?: () => void;
}

export function AppError({ error, onRetry }: AppErrorProps) {
  const { t, locale } = useI18n();
  const normalized = toApiError(error, t("error"));
  const message = userMessageForError(normalized as ApiError, locale);

  return (
    <Alert
      type="error"
      showIcon
      title={t("error")}
      description={
        <div>
          <p>{message}</p>
          {normalized.requestId ? (
            <p style={{ marginBottom: 0, opacity: 0.7 }}>{normalized.requestId}</p>
          ) : null}
        </div>
      }
      action={
        onRetry ? (
          <Button size="small" onClick={onRetry}>
            {t("retry")}
          </Button>
        ) : undefined
      }
    />
  );
}
