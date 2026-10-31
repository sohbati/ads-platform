import { describe, expect, it } from "vitest";
import { createRateLimiter } from "./rateLimit";

describe("createRateLimiter", () => {
  it("blocks after the limit until the window ends", () => {
    const limiter = createRateLimiter(2, 60_000);
    limiter.fail("k", 0);
    expect(limiter.retryAfter("k", 0)).toBe(0);
    limiter.fail("k", 0);
    expect(limiter.retryAfter("k", 1_000)).toBe(59);
    expect(limiter.retryAfter("k", 60_000)).toBe(0);
  });

  it("reset clears a key", () => {
    const limiter = createRateLimiter(1, 60_000);
    limiter.fail("k", 0);
    limiter.reset("k");
    expect(limiter.retryAfter("k", 0)).toBe(0);
  });
});
