export type UserRole = "student" | "teacher" | "manager" | "admin";

export interface User {
  id: string;
  email: string;
  role: UserRole;
  is_active: boolean;
  first_name?: string;
  last_name?: string;
  avatar_url?: string;
  group_name?: string; // for students
  department?: string; // for teachers
  created_at: string;
  updated_at?: string;
}

export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  group_name?: string;
  department?: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface StudentProfile {
  user_id: string;
  first_name: string;
  last_name: string;
  group_name: string;
  email: string;
}

export interface TeacherProfile {
  user_id: string;
  first_name: string;
  last_name: string;
  department: string;
  email: string;
}
