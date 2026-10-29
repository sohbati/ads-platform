import { z } from "zod";

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
