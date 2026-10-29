export const Permission = {
  customerRead: "customer.read",
  customerUpdate: "customer.update",
  customerBlock: "customer.block",
  advertisementRead: "advertisement.read",
  advertisementApprove: "advertisement.approve",
  advertisementReject: "advertisement.reject",
  advertisementDelete: "advertisement.delete",
  categoryRead: "category.read",
  categoryCreate: "category.create",
  categoryUpdate: "category.update",
  categoryDelete: "category.delete",
  locationRead: "location.read",
  reportRead: "report.read",
  contentRead: "content.read",
  contentUpdate: "content.update",
  notificationRead: "notification.read",
  notificationSend: "notification.send",
  staffRead: "staff.read",
  staffUpdate: "staff.update",
  roleRead: "role.read",
  roleUpdate: "role.update",
  actionRead: "action.read",
  actionUpdate: "action.update",
  staffGroupRead: "staff_group.read",
  staffGroupUpdate: "staff_group.update",
  settingsRead: "settings.read",
} as const;

export type PermissionCode = (typeof Permission)[keyof typeof Permission];

export const ALL_PERMISSIONS: PermissionCode[] = Object.values(Permission);
