import type { ApiError, ApiErrorKind } from "./types";

export class AppApiError extends Error implements ApiError {
  readonly kind: ApiErrorKind;
  readonly code: string;
  readonly details?: unknown;
  readonly requestId?: string;
  readonly status?: number;

  constructor(error: ApiError) {
    super(error.message);
    this.name = "AppApiError";
    this.kind = error.kind;
    this.code = error.code;
    this.details = error.details;
    this.requestId = error.requestId;
    this.status = error.status;
  }
}

export function kindFromStatus(status: number): ApiErrorKind {
  if (status === 400) {
    return "validation";
  }
  if (status === 401) {
    return "authentication";
  }
  if (status === 403) {
    return "authorization";
  }
  if (status === 404) {
    return "not_found";
  }
  if (status === 409) {
    return "conflict";
  }
  return "server";
}

export function toApiError(error: unknown, fallbackMessage: string): ApiError {
  if (error instanceof AppApiError) {
    return error;
  }
  if (error instanceof Error && error.message) {
    return {
      kind: "network",
      code: "NETWORK_ERROR",
      message: fallbackMessage,
      details: error.message,
    };
  }
  return {
    kind: "server",
    code: "UNKNOWN_ERROR",
    message: fallbackMessage,
  };
}

export function userMessageForError(error: ApiError, locale: "fa" | "en"): string {
  if (error.code === "AUTH_INVALID_CREDENTIALS") {
    return locale === "fa"
      ? "شماره موبایل یا رمز عبور نادرست است."
      : "Invalid mobile number or password.";
  }

  const fa: Record<ApiErrorKind, string> = {
    validation: "اطلاعات واردشده معتبر نیست.",
    authentication: "نشست شما منقضی شده است. دوباره وارد شوید.",
    authorization: "اجازه این کار را ندارید.",
    not_found: "موردی پیدا نشد.",
    conflict: "این مورد با داده موجود تداخل دارد.",
    server: "خطای سرور. دوباره تلاش کنید.",
    network: "ارتباط با سرور برقرار نشد.",
  };
  const en: Record<ApiErrorKind, string> = {
    validation: "The submitted data is not valid.",
    authentication: "Your session expired. Please sign in again.",
    authorization: "You do not have permission for this action.",
    not_found: "The requested item was not found.",
    conflict: "This change conflicts with existing data.",
    server: "A server error occurred. Try again.",
    network: "Could not reach the server.",
  };
  return locale === "fa" ? fa[error.kind] : en[error.kind];
}
