import { describe, expect, it } from "vitest";
import { decodeToken, encodeToken } from "./token";

const SECRET = "test-secret-test-secret-test-secret";
const NOW = 1_800_000_000_000;

describe("session token", () => {
  const token = { sid: 7, exp: NOW / 1000 + 60 };

  it("round-trips a valid token", () => {
    expect(decodeToken(encodeToken(token, SECRET), SECRET, NOW)).toEqual(token);
  });

  it("rejects a token signed with another secret", () => {
    expect(decodeToken(encodeToken(token, "other-secret"), SECRET, NOW)).toBeNull();
  });

  it("rejects a tampered payload", () => {
    const [, signature] = encodeToken(token, SECRET).split(".");
    const forged = Buffer.from(JSON.stringify({ sid: 1, exp: token.exp })).toString(
      "base64url",
    );
    expect(decodeToken(`${forged}.${signature}`, SECRET, NOW)).toBeNull();
  });

  it("rejects an expired token", () => {
    expect(decodeToken(encodeToken(token, SECRET), SECRET, NOW + 61_000)).toBeNull();
  });

  it("rejects garbage", () => {
    expect(decodeToken("not-a-token", SECRET, NOW)).toBeNull();
    expect(decodeToken(undefined, SECRET, NOW)).toBeNull();
  });
});
