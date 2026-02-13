import { z } from "zod";

export const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(1, "Password is required"),
  remember: z.boolean().optional(),
});

export type LoginFormData = z.infer<typeof loginSchema>;

export const signupSchema = z
  .object({
    email: z.string().email("Invalid email address"),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
      .regex(/[a-z]/, "Password must contain at least one lowercase letter")
      .regex(/[0-9]/, "Password must contain at least one number"),
    confirmPassword: z.string(),
    role: z.enum(["student", "teacher"]),
    first_name: z.string().min(1, "First name is required"),
    last_name: z.string().min(1, "Last name is required"),
    group_name: z.string().optional(),
    department: z.string().optional(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
  })
  .refine(
    (data) => {
      if (data.role === "student" && !data.group_name) {
        return false;
      }
      return true;
    },
    {
      message: "Group name is required for students",
      path: ["group_name"],
    },
  )
  .refine(
    (data) => {
      if (data.role === "teacher" && !data.department) {
        return false;
      }
      return true;
    },
    {
      message: "Department is required for teachers",
      path: ["department"],
    },
  );

export type SignupFormData = z.infer<typeof signupSchema>;
