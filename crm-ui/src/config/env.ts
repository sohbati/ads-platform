const DEFAULT_BACKEND_URL = "http://localhost:8092";
const DEFAULT_CDN_URL = "http://localhost:4000";
const DEV_SESSION_SECRET = "dev-only-change-me";
const MIN_SECRET_LENGTH = 32;

function read(name: string, fallback: string): string {
  const value = process.env[name];
  if (value && value.trim() !== "") {
    return value.trim();
  }
  return fallback;
}

// Read lazily so `next build` does not need runtime secrets.
function sessionSecret(): string {
  const secret = read("SESSION_SECRET", "");
  if (process.env.NODE_ENV === "production") {
    if (secret.length < MIN_SECRET_LENGTH || secret === DEV_SESSION_SECRET) {
      throw new Error(
        `SESSION_SECRET must be set to a random value of at least ${MIN_SECRET_LENGTH} characters in production.`,
      );
    }
    return secret;
  }
  return secret || DEV_SESSION_SECRET;
}

export const env = {
  backendApiBaseUrl: read("BACKEND_API_BASE_URL", DEFAULT_BACKEND_URL).replace(/\/$/, ""),
  cdnBaseUrl: read("CDN_BASE_URL", DEFAULT_CDN_URL).replace(/\/$/, ""),
  get sessionSecret() {
    return sessionSecret();
  },
} as const;
