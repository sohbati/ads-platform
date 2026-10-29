"use client";

import { Button, Form, Input, Modal } from "antd";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { namedEntitySchema } from "@/lib/validation/common";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { NamedEntityInput } from "../types/access.types";

interface NamedEntityFormProps {
  open: boolean;
  loading?: boolean;
  onCancel: () => void;
  onSubmit: (values: NamedEntityInput) => void;
}

export function NamedEntityForm({
  open,
  loading,
  onCancel,
  onSubmit,
}: NamedEntityFormProps) {
  const { t } = useI18n();
  const { control, handleSubmit, reset } = useForm<NamedEntityInput>({
    resolver: zodResolver(namedEntitySchema),
    defaultValues: { code: "", name: "" },
  });

  return (
    <Modal
      open={open}
      title={t("create")}
      footer={null}
      onCancel={() => {
        reset();
        onCancel();
      }}
    >
      <Form
        layout="vertical"
        onFinish={handleSubmit((values) => {
          onSubmit(values);
          reset();
        })}
      >
        <Controller
          name="code"
          control={control}
          render={({ field, fieldState }) => (
            <Form.Item
              label={t("code")}
              validateStatus={fieldState.error ? "error" : undefined}
              help={fieldState.error?.message}
            >
              <Input {...field} />
            </Form.Item>
          )}
        />
        <Controller
          name="name"
          control={control}
          render={({ field, fieldState }) => (
            <Form.Item
              label={t("name")}
              validateStatus={fieldState.error ? "error" : undefined}
              help={fieldState.error?.message}
            >
              <Input {...field} />
            </Form.Item>
          )}
        />
        <Button type="primary" htmlType="submit" loading={loading}>
          {t("save")}
        </Button>
      </Form>
    </Modal>
  );
}
