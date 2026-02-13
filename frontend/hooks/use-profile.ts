import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api-client";
import { useAuthStore } from "@/store/auth-store";

interface StudentProfile {
  id: string;
  user_id: string;
  first_name: string;
  last_name: string;
  student_code?: string;
  major?: string;
  year?: number;
  gpa?: number;
  enrollment_status?: string;
  group_id?: string;
  created_at: string;
  updated_at?: string;
}

interface TeacherProfile {
  id: string;
  user_id: string;
  first_name: string;
  last_name: string;
  department?: string;
  created_at: string;
  updated_at?: string;
}

export function useStudentProfile() {
  const { user } = useAuthStore();
  const isStudent = user?.role === "student";

  return useQuery<StudentProfile>({
    queryKey: ["student-profile", "me"],
    queryFn: () => api.get<StudentProfile>("/students/me"),
    enabled: isStudent,
    staleTime: 5 * 60 * 1000, // Cache for 5 minutes
  });
}

export function useTeacherProfile() {
  const { user } = useAuthStore();
  const isTeacher = user?.role === "teacher";

  return useQuery<TeacherProfile>({
    queryKey: ["teacher-profile", "me"],
    queryFn: () => api.get<TeacherProfile>("/teachers/me"),
    enabled: isTeacher,
    staleTime: 5 * 60 * 1000,
  });
}
