import { z } from "zod";

export const createAssignmentSchema = z.object({
  title: z.string().min(3, "Title must be at least 3 characters"),
  description: z.string().optional(),
  due_at: z.string().optional(),
  max_points: z
    .number()
    .min(0)
    .max(1000, "Max points must be between 0 and 1000"),
});

export type CreateAssignmentFormData = z.infer<typeof createAssignmentSchema>;

export const updateAssignmentSchema = z.object({
  title: z.string().min(3, "Title must be at least 3 characters").optional(),
  description: z.string().optional(),
  due_at: z.string().optional(),
  max_points: z.number().min(0).max(1000).optional(),
});

export type UpdateAssignmentFormData = z.infer<typeof updateAssignmentSchema>;

export const submitAssignmentSchema = z
  .object({
    content_text: z.string().max(50000, "Content is too long").optional(),
    file_url: z.string().url().optional(),
  })
  .refine((data) => data.content_text || data.file_url, {
    message: "Please provide either text content or a file",
    path: ["content_text"],
  });

export type SubmitAssignmentFormData = z.infer<typeof submitAssignmentSchema>;

export const gradeSubmissionSchema = z.object({
  score: z.number().min(0, "Score must be at least 0"),
  feedback: z.string().optional(),
});

export type GradeSubmissionFormData = z.infer<typeof gradeSubmissionSchema>;
