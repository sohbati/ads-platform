import { describe, expect, it } from "vitest";
import { isDemoCredentials } from "./session";

describe("isDemoCredentials", () => {
  it("accepts latin demo digits", () => {
    expect(isDemoCredentials("09120000000", "admin")).toBe(true);
  });

  it("accepts persian demo digits shown on the login page", () => {
    expect(isDemoCredentials("۰۹۱۲۰۰۰۰۰۰۰", "admin")).toBe(true);
  });

  it("accepts arabic-indic demo digits", () => {
    expect(isDemoCredentials("٠٩١٢٠٠٠٠٠٠٠", "admin")).toBe(true);
  });

  it("rejects a wrong password", () => {
    expect(isDemoCredentials("۰۹۱۲۰۰۰۰۰۰۰", "wrong")).toBe(false);
  });
});
