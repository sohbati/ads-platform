import { createHmac, timingSafeEqual } from "node:crypto";
import type { SessionToken } from "./session.types";

export const SESSION_TTL_SECONDS = 60 * 60 * 12;

function sign(payload: string, secret: string): string {
  return createHmac("sha256", secret).update(payload).digest("base64url");
}

export function encodeToken(token: SessionToken, secret: string): string {
  const payload = Buffer.from(JSON.stringify(token), "utf8").toString("base64url");
  return `${payload}.${sign(payload, secret)}`;
}

export function decodeToken(
  value: string | undefined,
  secret: string,
  now = Date.now(),
): SessionToken | null {
  if (!value) {
    return null;
  }
  const [payload, signature] = value.split(".");
  if (!payload || !signature) {
    return null;
  }
  const left = Buffer.from(signature);
  const right = Buffer.from(sign(payload, secret));
  if (left.length !== right.length || !timingSafeEqual(left, right)) {
    return null;
  }
  try {
    const parsed = JSON.parse(
      Buffer.from(payload, "base64url").toString("utf8"),
    ) as SessionToken;
    if (
      !Number.isInteger(parsed.sid) ||
      parsed.sid <= 0 ||
      typeof parsed.exp !== "number"
    ) {
      return null;
    }
    return parsed.exp * 1000 > now ? parsed : null;
  } catch {
    return null;
  }
}
