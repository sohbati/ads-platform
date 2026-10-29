import type { AxiosInstance, AxiosError } from "axios";
import { AppApiError, kindFromStatus } from "./errors";

interface BackendErrorBody {
  error?: string;
  message?: string;
  statusCode?: number;
  requestId?: string;
}

function newRequestId(): string {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return `req-${Date.now()}`;
}

export function attachInterceptors(client: AxiosInstance): void {
  client.interceptors.request.use((config) => {
    const requestId = newRequestId();
    config.headers.set("X-Request-Id", requestId);
    return config;
  });

  client.interceptors.response.use(
    (response) => response,
    (error: AxiosError<BackendErrorBody>) => {
      if (!error.response) {
        return Promise.reject(
          new AppApiError({
            kind: "network",
            code: "NETWORK_ERROR",
            message: "Could not reach the server.",
          }),
        );
      }
      const status = error.response.status;
      const body = error.response.data;
      const requestId =
        error.response.headers["x-request-id"] ??
        error.config?.headers.get("X-Request-Id") ??
        undefined;
      return Promise.reject(
        new AppApiError({
          kind: kindFromStatus(status),
          code: body?.error ?? `HTTP_${status}`,
          message: body?.message ?? "Request failed.",
          details: body,
          requestId: typeof requestId === "string" ? requestId : undefined,
          status,
        }),
      );
    },
  );
}
