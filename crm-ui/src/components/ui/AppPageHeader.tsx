"use client";

import { Typography } from "antd";

interface AppPageHeaderProps {
  title: string;
  description?: string;
  extra?: React.ReactNode;
}

export function AppPageHeader({ title, description, extra }: AppPageHeaderProps) {
  return (
    <div className="app-page-header">
      <div>
        <Typography.Title level={3} style={{ margin: 0 }}>
          {title}
        </Typography.Title>
        {description ? (
          <Typography.Paragraph type="secondary" style={{ marginTop: 8, marginBottom: 0 }}>
            {description}
          </Typography.Paragraph>
        ) : null}
      </div>
      {extra ? <div className="app-page-header__extra">{extra}</div> : null}
    </div>
  );
}
