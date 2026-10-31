"use client";

import { useEffect } from "react";
import { Button, Form, Input, Modal } from "antd";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { namedEntitySchema, renameEntitySchema } from "@/lib/validation/common";
import { useI18n } from "@/lib/i18n/I18nProvider";
import type { NamedEntity, NamedEntityInput } from "../types/access.types";

interface NamedEntityFormProps {
  open: boolean;
  /** When set, the form edits this item. */
  item?: NamedEntity | null;
  /** Whether the code may change when editing. */
  codeEditable?: boolean;
  loading?: boolean;
  onCancel: () => void;
  onSubmit: (values: NamedEntityInput) => void;
}

const EMPTY: NamedEntityInput = { code: "", name: "" };

// Used when the code is locked: only the name is validated.
const editSchema = renameEntitySchema.extend({ code: z.string() });

export function NamedEntityForm({
  open,
  item,
  codeEditable = false,
  loading,
  onCancel,
  onSubmit,
}: NamedEntityFormProps) {
  const { t } = useI18n();
  const editing = Boolean(item);
  const codeLocked = editing && !codeEditable;
  const { control, handleSubmit, reset } = useForm<NamedEntityInput>({
    resolver: zodResolver(codeLocked ? editSchema : namedEntitySchema),
    defaultValues: EMPTY,
  });

  useEffect(() => {
    if (open) {
      reset(item ? { code: item.code, name: item.name } : EMPTY);
    }
  }, [open, item, reset]);

  return (
    <Modal
      open={open}
      title={editing ? t("edit") : t("create")}
      footer={null}
      destroyOnHidden
      onCancel={onCancel}
    >
      <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
        <Controller
          name="code"
          control={control}
          render={({ field, fieldState }) => (
            <Form.Item
              label={t("code")}
              validateStatus={fieldState.error ? "error" : undefined}
              help={
                fieldState.error?.message ?? (codeLocked ? t("codeImmutable") : undefined)
              }
            >
              <Input {...field} disabled={codeLocked} dir="ltr" />
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
