export type SubmissionStatus = "pending" | "submitted" | "graded" | "late";

export interface Submission {
  id: string;
  assignment_id: string;
  student_id: string;
  content_text?: string;
  file_url?: string;
  status: SubmissionStatus;
  submitted_at: string;
  assignment_title?: string;
  course_title?: string;
  student_first_name?: string;
  student_last_name?: string;
  student_email?: string;
  max_points?: number;
  due_at?: string;
  grade_score?: number;
  grade_feedback?: string;
}

export interface CreateSubmissionRequest {
  assignment_id: string;
  content_text?: string;
  file_url?: string;
}

export interface UpdateSubmissionRequest {
  content_text?: string;
  file_url?: string;
}

export interface SubmissionListResponse {
  submissions: Submission[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}
