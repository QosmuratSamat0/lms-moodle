import api from "@/lib/api-client";
import type {
  SessionListResponse,
  MarksResponse,
  StudentAttendanceResponse,
  CreateSessionRequest,
  BulkMarkRequest,
  AttendanceSession,
} from "@/types/attendance";

export const attendanceService = {
  // Sessions
  createSession: async (
    data: CreateSessionRequest,
  ): Promise<AttendanceSession> => {
    return api.post<AttendanceSession>("/attendance/sessions", data);
  },

  getSession: async (id: string): Promise<AttendanceSession> => {
    return api.get<AttendanceSession>(`/attendance/sessions/${id}`);
  },

  listSessions: async (courseId: string): Promise<SessionListResponse> => {
    return api.get<SessionListResponse>(
      `/attendance/course/${courseId}/sessions`,
    );
  },

  deleteSession: async (id: string): Promise<void> => {
    return api.delete(`/attendance/sessions/${id}`);
  },

  // Marks
  getMarksBySession: async (sessionId: string): Promise<MarksResponse> => {
    return api.get<MarksResponse>(`/attendance/marks/session/${sessionId}`);
  },

  bulkMark: async (data: BulkMarkRequest): Promise<void> => {
    return api.post("/attendance/marks/bulk", data);
  },

  // Student view
  getMyAttendance: async (
    courseId: string,
  ): Promise<StudentAttendanceResponse> => {
    return api.get<StudentAttendanceResponse>(
      `/attendance/course/${courseId}/me`,
    );
  },
};
