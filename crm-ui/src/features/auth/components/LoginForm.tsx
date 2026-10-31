"use client";

import { Alert, Button, Card, Form, Input, Typography } from "antd";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { userMessageForError, toApiError } from "@/lib/api/errors";
import { loginSchema, type LoginValues } from "../schemas/login.schema";
import { useLogin } from "../hooks/useLogin";

export function LoginForm() {
  const { t, locale } = useI18n();
  const loginMutation = useLogin();
  const { control, handleSubmit } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { mobile: "", password: "" },
  });

  return (
    <Card className="login-card">
      <Typography.Title level={3} style={{ marginTop: 0 }}>
        {t("appName")}
      </Typography.Title>
      <Typography.Paragraph type="secondary">{t("loginHint")}</Typography.Paragraph>
      {loginMutation.error ? (
        <Alert
          type="error"
          showIcon
          style={{ marginBottom: 16 }}
          title={userMessageForError(toApiError(loginMutation.error, t("error")), locale)}
        />
      ) : null}
      <Form
        layout="vertical"
        onFinish={handleSubmit((values) => loginMutation.mutate(values))}
      >
        <Controller
          name="mobile"
          control={control}
          render={({ field, fieldState }) => (
            <Form.Item
              label={t("mobile")}
              validateStatus={fieldState.error ? "error" : undefined}
              help={fieldState.error?.message}
            >
              <Input {...field} autoComplete="username" inputMode="tel" />
            </Form.Item>
          )}
        />
        <Controller
          name="password"
          control={control}
          render={({ field, fieldState }) => (
            <Form.Item
              label={t("password")}
              validateStatus={fieldState.error ? "error" : undefined}
              help={fieldState.error?.message}
            >
              <Input.Password {...field} autoComplete="current-password" />
            </Form.Item>
          )}
        />
        <Button type="primary" htmlType="submit" block loading={loginMutation.isPending}>
          {t("signIn")}
        </Button>
      </Form>
    </Card>
  );
}
