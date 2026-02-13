export interface Group {
  id: string;
  code: string;
  name?: string;
  description?: string;
  year_of_admission?: number;
  student_count: number;
  created_at: string;
}

export interface GroupListResponse {
  groups: Group[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface TeacherAssignment {
  id: string;
  teacher_id: string;
  course_id: string;
  group_id: string;
  teacher_first_name?: string;
  teacher_last_name?: string;
  course_title: string;
  group_code: string;
  assigned_at: string;
}

export interface TeacherAssignmentListResponse {
  assignments: TeacherAssignment[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface AssignTeacherRequest {
  teacher_id: string;
  course_id: string;
  group_id: string;
}

export interface CreateGroupRequest {
  code: string;
  name?: string;
  description?: string;
  year_of_admission?: number;
}

export interface UpdateGroupRequest {
  code?: string;
  name?: string;
  description?: string;
  year_of_admission?: number;
}
