import { describe, expect, it } from "vitest";
import { userMessageForError } from "./errors";

describe("userMessageForError", () => {
  it("does not call a login 401 a session expiry", () => {
    expect(
      userMessageForError(
        {
          kind: "authentication",
          code: "AUTH_INVALID_CREDENTIALS",
          message: "Request failed.",
          status: 401,
        },
        "fa",
      ),
    ).toBe("شماره موبایل یا رمز عبور نادرست است.");
  });
});
