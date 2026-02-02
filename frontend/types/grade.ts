export interface Grade {
  id: string;
  submission_id: string;
  student_id?: string;
  assignment_id?: string;
  score: number;
  max_points: number;
  percentage: number;
  feedback?: string;
  graded_by_teacher_id?: string;
  teacher_first_name?: string;
  teacher_last_name?: string;
  assignment_title?: string;
  course_title?: string;
  student_first_name?: string;
  student_last_name?: string;
  student_email?: string;
  graded_at: string;
}

export interface GradeSubmissionRequest {
  submission_id: string;
  score: number;
  feedback?: string;
}

export interface UpdateGradeRequest {
  score?: number;
  feedback?: string;
}

export interface GradeListResponse {
  grades: Grade[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}
