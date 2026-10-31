import { describe, expect, it } from "vitest";
import { requiredPermission, targetStaffId } from "./panelAccess";

describe("requiredPermission", () => {
  it.each([
    ["GET", "staff", "staff.read"],
    ["GET", "staff/3/roles", "staff.read"],
    ["POST", "staff", "staff.update"],
    ["PUT", "staff/3/password", "staff.update"],
    ["POST", "staff/3/roles/9", "staff.update"],
    ["POST", "staff-groups/2/staff/3", "staff.update"],
    ["GET", "staff-groups", "staff_group.read"],
    ["DELETE", "staff-groups/2", "staff_group.update"],
    ["PUT", "roles/4/actions", "role.update"],
    ["GET", "actions", "action.read"],
    ["GET", "ads", "advertisement.read"],
    ["GET", "ads/12", "advertisement.read"],
    ["POST", "ads/12/approve", "advertisement.approve"],
    ["POST", "ads/12/reject", "advertisement.reject"],
    ["GET", "settings", "settings.read"],
    ["PUT", "settings/ads.pending_max_days", "settings.update"],
  ])("%s %s needs %s", (method, path, permission) => {
    expect(requiredPermission(method, path)).toBe(permission);
  });

  it.each([
    ["POST", "auth/login"],
    ["GET", "auth/session/1"],
    ["GET", "staff/abc"],
    ["GET", "staff/1/../../auth/session/1"],
    ["GET", "unknown"],
    ["GET", "ads/12/approve"],
    ["PUT", "ads/12"],
    ["DELETE", "ads/12"],
    ["PUT", "settings"],
    ["PUT", "settings/Bad Key"],
  ])("denies %s %s", (method, path) => {
    expect(requiredPermission(method, path)).toBeNull();
  });
});

describe("targetStaffId", () => {
  it("finds the staff member a write targets", () => {
    expect(targetStaffId("PUT", "staff/5/password")).toBe(5);
    expect(targetStaffId("DELETE", "staff/5/roles/2")).toBe(5);
    expect(targetStaffId("POST", "staff-groups/2/staff/5")).toBe(5);
  });

  it("ignores reads and creates", () => {
    expect(targetStaffId("GET", "staff/5")).toBeNull();
    expect(targetStaffId("POST", "staff")).toBeNull();
  });
});
