import { describe, expect, it } from "vitest";
import { formatDate, formatDateTime } from "./date";

const TEHRAN = "Asia/Tehran";

describe("formatDate", () => {
  it.each([
    // 12:30 in Tehran on 13 Aban 1405.
    ["2026-11-04T09:00:00Z", "۱۴۰۵/۰۸/۱۳"],
    // Nowruz 1405 starts at midnight Tehran time, 20:30 UTC the day before.
    ["2026-03-20T20:00:00Z", "۱۴۰۴/۱۲/۲۹"],
    ["2026-03-20T21:00:00Z", "۱۴۰۵/۰۱/۰۱"],
    // Last day of a leap year (1403 has 30 days in Esfand).
    ["2025-03-20T12:00:00Z", "۱۴۰۳/۱۲/۳۰"],
  ])("%s is %s in the Jalali calendar in Tehran", (value, expected) => {
    expect(formatDate(value, "fa", TEHRAN)).toBe(expected);
  });

  it("follows the viewer's zone: the same moment is a different day elsewhere", () => {
    expect(formatDate("2026-11-04T21:00:00Z", "en", TEHRAN)).toBe("05/11/2026");
    expect(formatDate("2026-11-04T21:00:00Z", "en", "America/Toronto")).toBe("04/11/2026");
    expect(formatDate("2026-11-04T21:00:00Z", "en", "UTC")).toBe("04/11/2026");
  });

  it.each([null, undefined, "", "not a date"])("shows a dash for %s", (value) => {
    expect(formatDate(value, "fa")).toBe("—");
  });
});

describe("formatDateTime", () => {
  it("adds the time in the viewer's zone", () => {
    expect(formatDateTime("2026-11-04T09:00:00Z", "fa", TEHRAN)).toBe("۱۴۰۵/۰۸/۱۳ ۱۲:۳۰");
    expect(formatDateTime("2026-11-04T09:00:00Z", "en", "Europe/Berlin")).toBe("04/11/2026 10:00");
  });

  it("uses 00 for midnight, not 24", () => {
    expect(formatDateTime("2026-03-20T20:30:00Z", "en", TEHRAN)).toBe("21/03/2026 00:00");
  });
});
