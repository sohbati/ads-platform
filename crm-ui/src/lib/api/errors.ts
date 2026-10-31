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

const CODE_MESSAGES: Record<string, { fa: string; en: string }> = {
  AUTH_INVALID_CREDENTIALS: {
    fa: "شماره موبایل یا رمز عبور نادرست است.",
    en: "Invalid mobile number or password.",
  },
  AUTH_TOO_MANY_ATTEMPTS: {
    fa: "تلاش‌های ناموفق زیادی انجام شد. چند دقیقه دیگر دوباره تلاش کنید.",
    en: "Too many failed attempts. Try again in a few minutes.",
  },
  STAFF_INVALID_PASSWORD: {
    fa: "رمز عبور باید دست‌کم ۸ نویسه باشد.",
    en: "Password must be at least 8 characters.",
  },
  STAFF_SUPER_ADMIN_PROTECTED: {
    fa: "مدیر کل را نمی‌توان حذف یا غیرفعال کرد.",
    en: "A super admin cannot be deleted or deactivated.",
  },
  ROLE_DUPLICATE_CODE: {
    fa: "نقشی با این کد وجود دارد.",
    en: "A role with this code already exists.",
  },
  STAFF_GROUP_DUPLICATE_CODE: {
    fa: "گروهی با این کد وجود دارد.",
    en: "A group with this code already exists.",
  },
  STAFF_DUPLICATE_MOBILE: {
    fa: "کارمندی با این شماره موبایل وجود دارد.",
    en: "A staff member with this mobile number already exists.",
  },
  STAFF_DUPLICATE_NATIONAL_ID: {
    fa: "کارمندی با این کد ملی وجود دارد.",
    en: "A staff member with this national ID already exists.",
  },
  STAFF_INVALID_MOBILE: {
    fa: "شماره موبایل معتبر نیست.",
    en: "Mobile number is not valid.",
  },
  STAFF_INVALID_NATIONAL_ID: {
    fa: "کد ملی باید ۱۰ رقم باشد.",
    en: "National ID must be 10 digits.",
  },
};

export function userMessageForError(error: ApiError, locale: "fa" | "en"): string {
  const known = CODE_MESSAGES[error.code];
  if (known) {
    return known[locale];
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
