import { describe, expect, it } from "vitest";
import { Permission } from "@/lib/permissions/codes";
import { can, canAny } from "./permissions";

describe("permissions", () => {
  it("denies when the list is empty", () => {
    expect(can([], Permission.customerRead)).toBe(false);
  });

  it("allows an exact permission", () => {
    expect(can([Permission.customerRead], Permission.customerRead)).toBe(true);
  });

  it("allows wildcard", () => {
    expect(can(["*"], Permission.advertisementApprove)).toBe(true);
  });

  it("checks any of the required permissions", () => {
    expect(
      canAny([Permission.roleRead], [Permission.actionRead, Permission.roleRead]),
    ).toBe(true);
  });
});
