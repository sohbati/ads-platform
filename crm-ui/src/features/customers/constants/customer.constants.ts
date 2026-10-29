export const CustomerStatus = {
  ACTIVE: "ACTIVE",
  BLOCKED: "BLOCKED",
  SUSPENDED: "SUSPENDED",
} as const;

export type CustomerStatus = (typeof CustomerStatus)[keyof typeof CustomerStatus];
