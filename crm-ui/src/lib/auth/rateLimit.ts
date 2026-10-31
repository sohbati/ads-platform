interface Bucket {
  count: number;
  resetAt: number;
}

/**
 * Fixed-window attempt counter kept in process memory. It protects a single server instance;
 * running several instances needs a shared store (e.g. Redis) instead.
 */
export function createRateLimiter(limit: number, windowMs: number) {
  const buckets = new Map<string, Bucket>();

  return {
    /** Returns seconds until the key may try again, or 0 if it is allowed now. */
    retryAfter(key: string, now = Date.now()): number {
      const bucket = buckets.get(key);
      if (!bucket || bucket.resetAt <= now) {
        return 0;
      }
      return bucket.count >= limit ? Math.ceil((bucket.resetAt - now) / 1000) : 0;
    },
    fail(key: string, now = Date.now()): void {
      const bucket = buckets.get(key);
      if (!bucket || bucket.resetAt <= now) {
        buckets.set(key, { count: 1, resetAt: now + windowMs });
        return;
      }
      bucket.count += 1;
    },
    reset(key: string): void {
      buckets.delete(key);
    },
  };
}
