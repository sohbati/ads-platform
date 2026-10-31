import { env } from "@/config/env";
import { AppApiError, kindFromStatus } from "./errors";

interface BackendRequest {
  path: string;
  method?: string;
  body?: unknown;
  search?: URLSearchParams;
  headers?: Record<string, string>;
}

export async function backendFetch<T>(request: BackendRequest): Promise<T> {
  const url = new URL(`${env.backendApiBaseUrl}${request.path}`);
  if (request.search) {
    url.search = request.search.toString();
  }
  const response = await fetch(url, {
    method: request.method ?? "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      ...request.headers,
    },
    body: request.body === undefined ? undefined : JSON.stringify(request.body),
    cache: "no-store",
  });

  if (!response.ok) {
    let code = `HTTP_${response.status}`;
    try {
      const payload = (await response.json()) as { error?: string };
      if (payload.error) {
        code = payload.error;
      }
    } catch {
      // keep status code
    }
    throw new AppApiError({
      kind: kindFromStatus(response.status),
      code,
      message: "Backend request failed.",
      status: response.status,
    });
  }

  if (response.status === 204) {
    return undefined as T;
  }
  return (await response.json()) as T;
}

export function emptyPage<T>(page = 1, pageSize = 20) {
  return { items: [] as T[], total: 0, page, pageSize };
}
