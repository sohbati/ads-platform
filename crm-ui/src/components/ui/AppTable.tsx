"use client";

import { Table } from "antd";
import type { TableProps } from "antd";
import { AppEmpty } from "./AppEmpty";

interface AppTableProps<T extends object> extends TableProps<T> {
  emptyText?: string;
}

export function AppTable<T extends object>({
  emptyText,
  pagination,
  ...props
}: AppTableProps<T>) {
  return (
    <Table<T>
      size="middle"
      scroll={{ x: "max-content" }}
      locale={{ emptyText: <AppEmpty description={emptyText} /> }}
      pagination={
        pagination === false
          ? false
          : {
              showSizeChanger: true,
              pageSizeOptions: [10, 20, 50],
              showTotal: (total) => String(total),
              ...pagination,
            }
      }
      {...props}
    />
  );
}
