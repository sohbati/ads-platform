export type { AdminSession as SessionStaff } from "@/lib/auth/session.types";

export interface LoginInput {
  mobile: string;
  password: string;
}
