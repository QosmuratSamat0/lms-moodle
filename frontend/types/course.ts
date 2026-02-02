export interface Course {
  id: string;
  title: string;
  description?: string;
  teacher_id?: string;
  owner_teacher_id?: string;
  teacher_first_name?: string;
  teacher_last_name?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateCourseRequest {
  title: string;
  description?: string;
}

export interface UpdateCourseRequest {
  title?: string;
  description?: string;
  is_active?: boolean;
}

export interface CourseListResponse {
  courses: Course[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface CourseFilters {
  search?: string;
  is_active?: boolean;
  teacher_id?: string;
}

export type EnrollmentStatus = "active" | "pending" | "rejected" | "dropped";

export interface Enrollment {
  id: string;
  course_id: string;
  student_id: string;
  status: EnrollmentStatus;
  enrolled_at: string;
  course_title?: string;
  student_name?: string;
}

export interface EnrollRequest {
  course_id: string;
}

export interface UpdateEnrollmentRequest {
  status: EnrollmentStatus;
}
