import api from "@/lib/api-client";
import type {
  Group,
  GroupListResponse,
  TeacherAssignment,
  TeacherAssignmentListResponse,
  AssignTeacherRequest,
  CreateGroupRequest,
  UpdateGroupRequest,
} from "@/types/group";
import type { PaginationParams } from "@/types/common";

type QueryParams = Record<string, string | number | boolean | undefined>;

export const groupService = {
  // List all groups
  list: async (params?: PaginationParams): Promise<GroupListResponse> => {
    return api.get<GroupListResponse>("/groups", params as QueryParams);
  },

  // Get group by ID
  get: async (id: string): Promise<Group> => {
    return api.get<Group>(`/groups/${id}`);
  },

  // Get group by code (e.g., "SE-2430")
  getByCode: async (code: string): Promise<Group> => {
    return api.get<Group>(`/groups/code/${code}`);
  },

  // Create a new group (admin only)
  create: async (data: CreateGroupRequest): Promise<Group> => {
    return api.post<Group>("/groups", data);
  },

  // Update a group (admin only)
  update: async (id: string, data: UpdateGroupRequest): Promise<Group> => {
    return api.put<Group>(`/groups/${id}`, data);
  },

  // Delete a group (admin only)
  delete: async (id: string): Promise<void> => {
    return api.delete(`/groups/${id}`);
  },

  // Get current teacher's group assignments
  getMyAssignments: async (
    params?: PaginationParams,
  ): Promise<TeacherAssignmentListResponse> => {
    return api.get<TeacherAssignmentListResponse>(
      "/teachers/me/group-assignments",
      params as QueryParams,
    );
  },

  // Get a specific teacher's group assignments (admin only)
  getTeacherAssignments: async (
    teacherId: string,
    params?: PaginationParams,
  ): Promise<TeacherAssignmentListResponse> => {
    return api.get<TeacherAssignmentListResponse>(
      `/teachers/${teacherId}/group-assignments`,
      params as QueryParams,
    );
  },

  // Assign a teacher to a course-group (admin only)
  assignTeacher: async (
    data: AssignTeacherRequest,
  ): Promise<TeacherAssignment> => {
    return api.post<TeacherAssignment>("/groups/assignments", data);
  },

  // Remove a teacher assignment (admin only)
  unassignTeacher: async (assignmentId: string): Promise<void> => {
    return api.delete(`/groups/assignments/${assignmentId}`);
  },

  // Get group assignments for a specific group
  getGroupAssignments: async (
    groupId: string,
  ): Promise<TeacherAssignmentListResponse> => {
    return api.get<TeacherAssignmentListResponse>(
      `/groups/${groupId}/assignments`,
    );
  },

  // Get members of a group (with student names)
  getMembers: async (groupId: string): Promise<any[]> => {
    return api.get<any[]>(`/groups/${groupId}/members`);
  },
};

export default groupService;
