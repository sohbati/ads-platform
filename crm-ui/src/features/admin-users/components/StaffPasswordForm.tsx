"use client";

import { useEffect } from "react";
import { Button, Form, Input, Modal } from "antd";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { passwordSchema } from "@/lib/validation/common";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { MessageKey } from "@/lib/i18n/messages";
import type { StaffListItem } from "../types/staff.types";

const schema = z.object({ password: passwordSchema });

interface StaffPasswordFormProps {
  staff: StaffListItem | null;
  loading?: boolean;
  onCancel: () => void;
  onSubmit: (password: string) => void;
}

export function StaffPasswordForm({
  staff,
  loading,
  onCancel,
  onSubmit,
}: StaffPasswordFormProps) {
  const { t } = useI18n();
  const { control, handleSubmit, reset } = useForm<{ password: string }>({
    resolver: zodResolver(schema),
    defaultValues: { password: "" },
  });

  useEffect(() => {
    if (staff) {
      reset({ password: "" });
    }
  }, [staff, reset]);

  return (
    <Modal
      open={staff !== null}
      title={staff ? `${t("setPassword")}: ${staff.name}` : undefined}
      footer={null}
      destroyOnHidden
      onCancel={onCancel}
    >
      <Form
        layout="vertical"
        onFinish={handleSubmit((values) => onSubmit(values.password))}
      >
        <Controller
          name="password"
          control={control}
          render={({ field, fieldState }) => (
            <Form.Item
              label={t("newPassword")}
              required
              validateStatus={fieldState.error ? "error" : undefined}
              help={
                fieldState.error ? t(fieldState.error.message as MessageKey) : undefined
              }
            >
              <Input.Password {...field} autoComplete="new-password" />
            </Form.Item>
          )}
        />
        <div className="app-form-actions">
          <Button onClick={onCancel}>{t("cancel")}</Button>
          <Button type="primary" htmlType="submit" loading={loading}>
            {t("save")}
          </Button>
        </div>
      </Form>
    </Modal>
  );
}
