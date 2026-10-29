"use client";

import { Form, Input, Modal } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";

interface ConfirmDialogProps {
  open: boolean;
  title: string;
  loading?: boolean;
  requireReason?: boolean;
  onCancel: () => void;
  onConfirm: (reason?: string) => void;
}

export function ConfirmDialog({
  open,
  title,
  loading,
  requireReason,
  onCancel,
  onConfirm,
}: ConfirmDialogProps) {
  const { t } = useI18n();
  const [form] = Form.useForm<{ reason?: string }>();

  return (
    <Modal
      open={open}
      title={title}
      okText={t("confirm")}
      cancelText={t("cancel")}
      confirmLoading={loading}
      onCancel={() => {
        form.resetFields();
        onCancel();
      }}
      onOk={async () => {
        if (requireReason) {
          const values = await form.validateFields();
          onConfirm(values.reason);
          return;
        }
        onConfirm();
      }}
    >
      {requireReason ? (
        <Form form={form} layout="vertical">
          <Form.Item
            name="reason"
            label={t("reason")}
            rules={[{ required: true, min: 3 }]}
          >
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      ) : null}
    </Modal>
  );
}
