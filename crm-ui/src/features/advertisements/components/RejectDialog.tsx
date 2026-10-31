"use client";

import { useEffect } from "react";
import { Button, Form, Input, Modal, Typography } from "antd";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { MessageKey } from "@/lib/i18n/messages";
import { REASON_MAX_LENGTH } from "../constants/advertisement.constants";
import type { RejectInput } from "../types/advertisement.types";

const schema = z.object({
  customer_reason: z
    .string()
    .trim()
    .min(1, "customerReasonRequired")
    .max(REASON_MAX_LENGTH, "reasonTooLong"),
  internal_reason: z.string().trim().max(REASON_MAX_LENGTH, "reasonTooLong"),
});

type Values = z.infer<typeof schema>;

interface RejectDialogProps {
  open: boolean;
  loading?: boolean;
  onCancel: () => void;
  onSubmit: (input: RejectInput) => void;
}

export function RejectDialog({ open, loading, onCancel, onSubmit }: RejectDialogProps) {
  const { t } = useI18n();
  const { control, handleSubmit, reset } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: { customer_reason: "", internal_reason: "" },
  });

  useEffect(() => {
    if (open) {
      reset({ customer_reason: "", internal_reason: "" });
    }
  }, [open, reset]);

  const fields: {
    name: keyof Values;
    label: MessageKey;
    hint: MessageKey;
    required: boolean;
  }[] = [
    {
      name: "customer_reason",
      label: "customerReason",
      hint: "customerReasonHint",
      required: true,
    },
    {
      name: "internal_reason",
      label: "internalReason",
      hint: "internalReasonHint",
      required: false,
    },
  ];

  return (
    <Modal
      open={open}
      title={t("rejectAd")}
      footer={null}
      destroyOnHidden
      onCancel={onCancel}
    >
      <Typography.Paragraph type="secondary">{t("rejectAdHint")}</Typography.Paragraph>
      <Form
        layout="vertical"
        onFinish={handleSubmit((values) =>
          onSubmit({
            customer_reason: values.customer_reason,
            internal_reason: values.internal_reason || undefined,
          }),
        )}
      >
        {fields.map((field) => (
          <Controller
            key={field.name}
            name={field.name}
            control={control}
            render={({ field: input, fieldState }) => (
              <Form.Item
                label={t(field.label)}
                required={field.required}
                validateStatus={fieldState.error ? "error" : undefined}
                help={
                  fieldState.error
                    ? t(fieldState.error.message as MessageKey)
                    : t(field.hint)
                }
              >
                <Input.TextArea
                  {...input}
                  rows={3}
                  maxLength={REASON_MAX_LENGTH}
                  showCount
                />
              </Form.Item>
            )}
          />
        ))}
        <div className="app-form-actions">
          <Button onClick={onCancel}>{t("cancel")}</Button>
          <Button danger type="primary" htmlType="submit" loading={loading}>
            {t("reject")}
          </Button>
        </div>
      </Form>
    </Modal>
  );
}
