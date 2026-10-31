import { z } from "zod";
import { toLatinDigits } from "@/lib/formatters/digits";

export const namedEntitySchema = z.object({
  code: z
    .string()
    .trim()
    .min(2)
    .max(80)
    .regex(/^[a-z][a-z0-9._-]*$/i, "code must be alphanumeric"),
  name: z.string().trim().min(2).max(120),
});

export const renameEntitySchema = z.object({
  name: z.string().trim().min(2).max(120),
});

export const reasonSchema = z.object({
  reason: z.string().trim().min(3).max(500),
});

const latinDigits = z.string().transform((value) => toLatinDigits(value).trim());

// Messages are i18n keys; forms translate them.
export const staffSchema = z.object({
  name: z.string().trim().min(2, "invalidName").max(100, "invalidName"),
  mobile: latinDigits.pipe(
    z.string().regex(/^(\+98|0098|98|0)?9\d{9}$/, "invalidMobile"),
  ),
  national_id: latinDigits.pipe(z.string().regex(/^\d{10}$/, "invalidNationalId")),
});

export const passwordSchema = z
  .string()
  .min(8, "invalidPassword")
  .refine((value) => new TextEncoder().encode(value).length <= 72, "invalidPassword");

export const staffCreateSchema = staffSchema.extend({
  password: passwordSchema,
  is_active: z.boolean().optional(),
});

// Passwords change through their own endpoint, so updates ignore the field.
export const staffUpdateSchema = staffSchema.extend({
  password: z.string().optional(),
  is_active: z.boolean().optional(),
});
