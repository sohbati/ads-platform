// Values safe to ship to the browser. Next inlines NEXT_PUBLIC_* at build time.
export const publicEnv = {
  /** Base URL for ad photos, which the backend stores as paths like /ads-media/... */
  mediaBaseUrl: (
    process.env.NEXT_PUBLIC_MEDIA_BASE_URL ?? "http://localhost:8098"
  ).replace(/\/$/, ""),
} as const;

export function mediaUrl(path: string | undefined): string | undefined {
  if (!path) {
    return undefined;
  }
  if (/^https?:\/\//.test(path)) {
    return path;
  }
  return `${publicEnv.mediaBaseUrl}${path.startsWith("/") ? "" : "/"}${path}`;
}
