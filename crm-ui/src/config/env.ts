const DEFAULT_BACKEND_URL = "http://localhost:8092";
const DEFAULT_SESSION_SECRET = "dev-only-change-me";

function read(name: string, fallback: string): string {
  const value = process.env[name];
  if (value && value.trim() !== "") {
    return value.trim();
  }
  return fallback;
}

export const env = {
  backendApiBaseUrl: read("BACKEND_API_BASE_URL", DEFAULT_BACKEND_URL).replace(
    /\/$/,
    "",
  ),
  sessionSecret: read("SESSION_SECRET", DEFAULT_SESSION_SECRET),
} as const;
