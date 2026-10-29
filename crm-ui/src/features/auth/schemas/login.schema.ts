import { z } from "zod";

export const loginSchema = z.object({
  mobile: z.string().trim().min(10).max(15),
  password: z.string().min(1),
});

export type LoginValues = z.infer<typeof loginSchema>;
