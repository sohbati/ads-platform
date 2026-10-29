export interface PagedResult<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ListQuery {
  page?: number;
  pageSize?: number;
  search?: string;
  sort?: string;
  order?: "asc" | "desc";
}

export type ApiErrorKind =
  | "validation"
  | "authentication"
  | "authorization"
  | "not_found"
  | "conflict"
  | "server"
  | "network";

export interface ApiError {
  kind: ApiErrorKind;
  code: string;
  message: string;
  details?: unknown;
  requestId?: string;
  status?: number;
}
