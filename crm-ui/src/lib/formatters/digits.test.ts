import { describe, expect, it } from "vitest";
import { formatDigits, toLatinDigits, toPersianDigits } from "./digits";

describe("digits", () => {
  it("converts latin digits to persian", () => {
    expect(toPersianDigits("0912")).toBe("۰۹۱۲");
  });

  it("converts persian digits back to latin", () => {
    expect(toLatinDigits("۰۹۱۲")).toBe("0912");
  });

  it("keeps latin digits for english locale", () => {
    expect(formatDigits(12, "en")).toBe("12");
  });
});
