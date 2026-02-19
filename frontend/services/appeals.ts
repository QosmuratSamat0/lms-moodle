import api from "@/lib/api-client";
import type { PaginationParams } from "@/types/common";

// Types
export type AppealStatus = "pending" | "approved" | "rejected";

export interface GradeAppeal {
  id: string;
  grade_id: string;
  student_id: string;
  reason: string;
  evidence?: string;
  status: AppealStatus;
  teacher_id?: string;
  response?: string;
  new_score?: number;
  resolved_at?: string;
  created_at: string;
  updated_at: string;
}

export interface AppealWithDetails extends GradeAppeal {
  student_name: string;
  student_email: string;
  course_name: string;
  assignment_title: string;
  original_score: number;
  max_points: number;
}

// Input types
export interface CreateAppealInput {
  grade_id: string;
  reason: string;
  evidence?: string;
}

export interface ResolveAppealInput {
  status: "approved" | "rejected";
  response: string;
  new_score?: number;
}

// Response types
export interface AppealsListResponse {
  data: GradeAppeal[];
  total: number;
}

export interface AppealsWithDetailsListResponse {
  data: AppealWithDetails[];
  total: number;
}

type QueryParams = Record<string, string | number | boolean | undefined>;

export const appealService = {
  // Get student's appeals
  getMyAppeals: async (
    params?: PaginationParams & { status?: AppealStatus }
  ): Promise<AppealsListResponse> => {
    return api.get<AppealsListResponse>("/appeals", params as QueryParams);
  },

  // Get appeal by ID
  getById: async (id: string): Promise<GradeAppeal> => {
    return api.get<GradeAppeal>(`/appeals/${id}`);
  },

  // Create appeal (student)
  create: async (data: CreateAppealInput): Promise<GradeAppeal> => {
    return api.post<GradeAppeal>("/appeals", data);
  },

  // Delete appeal (student/admin)
  delete: async (id: string): Promise<void> => {
    return api.delete(`/appeals/${id}`);
  },

  // Get teacher's appeals
  getTeacherAppeals: async (
    params?: PaginationParams & { status?: AppealStatus }
  ): Promise<AppealsWithDetailsListResponse> => {
    return api.get<AppealsWithDetailsListResponse>("/teacher/appeals", params as QueryParams);
  },

  // Resolve appeal (teacher/admin)
  resolve: async (id: string, data: ResolveAppealInput): Promise<GradeAppeal> => {
    return api.put<GradeAppeal>(`/appeals/${id}/resolve`, data);
  },
};

export default appealService;
