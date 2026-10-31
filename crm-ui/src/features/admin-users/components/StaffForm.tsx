"use client";

import { useEffect } from "react";
import { Button, Form, Input, Modal, Switch } from "antd";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import type { z } from "zod";
import { staffCreateSchema, staffUpdateSchema } from "@/lib/validation/common";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { MessageKey } from "@/lib/i18n/messages";
import type { StaffInput, StaffListItem } from "../types/staff.types";

interface StaffFormProps {
  open: boolean;
  staff?: StaffListItem | null;
  loading?: boolean;
  onCancel: () => void;
  onSubmit: (values: StaffInput) => void;
}

type TextField = "name" | "mobile" | "national_id" | "password";

const EMPTY: StaffInput = {
  name: "",
  mobile: "",
  national_id: "",
  password: "",
  is_active: true,
};

// Stored mobiles are E.164 (+98912…); show the local form when editing.
function localMobile(mobile: string): string {
  return mobile.startsWith("+98") ? `0${mobile.slice(3)}` : mobile;
}

export function StaffForm({ open, staff, loading, onCancel, onSubmit }: StaffFormProps) {
  const { t } = useI18n();
  const editing = Boolean(staff);
  const schema: z.ZodType<StaffInput, StaffInput> = editing
    ? staffUpdateSchema
    : staffCreateSchema;
  const { control, handleSubmit, reset } = useForm<StaffInput>({
    resolver: zodResolver(schema),
    defaultValues: EMPTY,
  });

  useEffect(() => {
    if (open) {
      reset(
        staff
          ? {
              name: staff.name,
              mobile: localMobile(staff.mobile),
              national_id: staff.national_id,
              is_active: staff.is_active,
            }
          : EMPTY,
      );
    }
  }, [open, staff, reset]);

  const fields: { name: TextField; label: MessageKey; inputMode?: "tel" | "numeric" }[] =
    [
      { name: "name", label: "name" },
      { name: "mobile", label: "mobile", inputMode: "tel" },
      { name: "national_id", label: "nationalId", inputMode: "numeric" },
      ...(editing ? [] : [{ name: "password" as const, label: "password" as const }]),
    ];

  return (
    <Modal
      open={open}
      title={editing ? t("editStaff") : t("addStaff")}
      footer={null}
      destroyOnHidden
      onCancel={onCancel}
    >
      <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
        {fields.map((item) => (
          <Controller
            key={item.name}
            name={item.name}
            control={control}
            render={({ field, fieldState }) => (
              <Form.Item
                label={t(item.label)}
                required
                validateStatus={fieldState.error ? "error" : undefined}
                help={
                  fieldState.error ? t(fieldState.error.message as MessageKey) : undefined
                }
              >
                {item.name === "password" ? (
                  <Input.Password {...field} autoComplete="new-password" />
                ) : (
                  <Input
                    {...field}
                    inputMode={item.inputMode}
                    dir={item.inputMode ? "ltr" : undefined}
                  />
                )}
              </Form.Item>
            )}
          />
        ))}
        {editing ? (
          <Controller
            name="is_active"
            control={control}
            render={({ field }) => (
              <Form.Item label={t("status")}>
                <Switch
                  checked={field.value}
                  onChange={field.onChange}
                  disabled={staff?.is_super_admin}
                  checkedChildren={t("active")}
                  unCheckedChildren={t("inactive")}
                />
              </Form.Item>
            )}
          />
        ) : null}
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
