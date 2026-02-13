// Grading categories for syllabus-based grading
export type GradingCategory =
  | "register_midterm"
  | "register_endterm"
  | "final"
  | "bonus";

export interface Assignment {
  id: string;
  course_id: string;
  title: string;
  description?: string;
  due_at?: string;
  max_points: number;
  weight_percentage?: number;
  grading_category?: GradingCategory;
  created_by_teacher_id?: string;
  file_url?: string;
  course_title?: string;
  teacher_first_name?: string;
  teacher_last_name?: string;
  submission_count?: number;
  created_at: string;
}

export interface CreateAssignmentRequest {
  course_id: string;
  title: string;
  description?: string;
  due_at?: string;
  max_points: number;
  weight_percentage?: number;
  grading_category?: GradingCategory;
  file_url?: string;
}

export interface UpdateAssignmentRequest {
  title?: string;
  description?: string;
  due_at?: string;
  max_points?: number;
  weight_percentage?: number;
  grading_category?: GradingCategory;
  file_url?: string;
}

export interface AssignmentListResponse {
  assignments: Assignment[];
  total: number;
  page?: number;
  limit?: number;
  total_pages?: number;
}

export interface AssignmentFilters {
  course_id?: string;
  search?: string;
}
