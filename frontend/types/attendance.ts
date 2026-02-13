export type AttendanceStatus = "present" | "absent" | "late" | "excused";

export interface AttendanceSession {
  id: string;
  course_id: string;
  starts_at: string;
  ends_at?: string;
  created_by_teacher_id?: string;
  created_at: string;
  total_students: number;
  present_count: number;
  absent_count: number;
  late_count: number;
  excused_count: number;
}

export interface AttendanceMark {
  id: string;
  session_id: string;
  student_id: string;
  status: AttendanceStatus;
  marked_at: string;
  student_first_name?: string;
  student_last_name?: string;
  student_email?: string;
}

export interface AttendanceSummary {
  course_id: string;
  total_classes: number;
  present: number;
  absent: number;
  late: number;
  excused: number;
  percentage: number;
}

export interface SessionListResponse {
  sessions: AttendanceSession[];
  total: number;
}

export interface MarksResponse {
  marks: AttendanceMark[];
  total: number;
}

export interface StudentAttendanceResponse {
  marks: AttendanceMark[];
  summary: AttendanceSummary;
}

export interface CreateSessionRequest {
  course_id: string;
  date: string; // YYYY-MM-DD
  start_time?: string; // HH:mm
  end_time?: string; // HH:mm
}

export interface BulkMarkRequest {
  session_id: string;
  marks: { student_id: string; status: AttendanceStatus }[];
}
