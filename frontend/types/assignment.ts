export interface Assignment {
  id: string;
  course_id: string;
  title: string;
  description?: string;
  due_at?: string;
  max_points: number;
  created_by_teacher_id?: string;
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
}

export interface UpdateAssignmentRequest {
  title?: string;
  description?: string;
  due_at?: string;
  max_points?: number;
}

export interface AssignmentListResponse {
  assignments: Assignment[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface AssignmentFilters {
  course_id?: string;
  search?: string;
}
