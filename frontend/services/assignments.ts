import api from "@/lib/api-client";
import type {
  Assignment,
  AssignmentListResponse,
  AssignmentFilters,
  CreateAssignmentRequest,
  UpdateAssignmentRequest,
} from "@/types/assignment";
import type { PaginationParams } from "@/types/common";

type QueryParams = Record<string, string | number | boolean | undefined>;

export const assignmentService = {
  list: async (
    params?: PaginationParams & AssignmentFilters,
  ): Promise<AssignmentListResponse> => {
    return api.get<AssignmentListResponse>(
      "/assignments",
      params as QueryParams,
    );
  },

  getByCourse: async (
    courseId: string,
    params?: PaginationParams,
  ): Promise<AssignmentListResponse> => {
    return api.get<AssignmentListResponse>(
      `/courses/${courseId}/assignments`,
      params as QueryParams,
    );
  },

  get: async (id: string): Promise<Assignment> => {
    return api.get<Assignment>(`/assignments/${id}`);
  },

  create: async (data: CreateAssignmentRequest): Promise<Assignment> => {
    return api.post<Assignment>("/assignments", data);
  },

  update: async (
    id: string,
    data: UpdateAssignmentRequest,
  ): Promise<Assignment> => {
    return api.put<Assignment>(`/assignments/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    return api.delete(`/assignments/${id}`);
  },
};

// Mock data for development
const mockAssignments: Assignment[] = [
  {
    id: "1",
    course_id: "1",
    title: "Programming Assignment 1: Hello World",
    description:
      'Write your first program that prints "Hello, World!" to the console.',
    due_at: "2026-02-01T23:59:59Z",
    max_points: 100,
    course_title: "Introduction to Computer Science",
    teacher_first_name: "Dr. Sarah",
    teacher_last_name: "Johnson",
    submission_count: 15,
    created_at: "2025-01-20T00:00:00Z",
  },
  {
    id: "2",
    course_id: "1",
    title: "Programming Assignment 2: Variables and Types",
    description: "Practice working with different data types and variables.",
    due_at: "2026-02-15T23:59:59Z",
    max_points: 100,
    course_title: "Introduction to Computer Science",
    teacher_first_name: "Dr. Sarah",
    teacher_last_name: "Johnson",
    submission_count: 8,
    created_at: "2025-01-25T00:00:00Z",
  },
  {
    id: "3",
    course_id: "2",
    title: "Calculus Problem Set 1",
    description: "Solve the following differential equations.",
    due_at: "2026-01-30T23:59:59Z",
    max_points: 50,
    course_title: "Advanced Mathematics",
    teacher_first_name: "Prof. Michael",
    teacher_last_name: "Chen",
    submission_count: 20,
    created_at: "2025-01-18T00:00:00Z",
  },
  {
    id: "4",
    course_id: "3",
    title: "Implement a Binary Search Tree",
    description: "Create a BST with insert, delete, and search operations.",
    due_at: "2026-02-10T23:59:59Z",
    max_points: 150,
    course_title: "Data Structures and Algorithms",
    teacher_first_name: "Dr. Sarah",
    teacher_last_name: "Johnson",
    submission_count: 5,
    created_at: "2025-01-22T00:00:00Z",
  },
];

export const assignmentServiceMock = {
  list: async (
    params?: PaginationParams & AssignmentFilters,
  ): Promise<AssignmentListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 500));

    let filtered = [...mockAssignments];

    if (params?.course_id) {
      filtered = filtered.filter((a) => a.course_id === params.course_id);
    }

    if (params?.search) {
      const search = params.search.toLowerCase();
      filtered = filtered.filter(
        (a) =>
          a.title.toLowerCase().includes(search) ||
          a.description?.toLowerCase().includes(search),
      );
    }

    return {
      assignments: filtered,
      total: filtered.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: Math.ceil(filtered.length / (params?.limit || 10)),
    };
  },

  getByCourse: async (
    courseId: string,
    params?: PaginationParams,
  ): Promise<AssignmentListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));

    const filtered = mockAssignments.filter((a) => a.course_id === courseId);

    return {
      assignments: filtered,
      total: filtered.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: Math.ceil(filtered.length / (params?.limit || 10)),
    };
  },

  get: async (id: string): Promise<Assignment> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    const assignment = mockAssignments.find((a) => a.id === id);
    if (!assignment) throw new Error("Assignment not found");
    return assignment;
  },

  create: async (data: CreateAssignmentRequest): Promise<Assignment> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: String(Date.now()),
      ...data,
      submission_count: 0,
      created_at: new Date().toISOString(),
    };
  },

  update: async (
    id: string,
    data: UpdateAssignmentRequest,
  ): Promise<Assignment> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    const assignment = mockAssignments.find((a) => a.id === id);
    if (!assignment) throw new Error("Assignment not found");
    return { ...assignment, ...data };
  },

  delete: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },
};

const useMock = process.env.NEXT_PUBLIC_USE_MOCK === "true";
export default useMock ? assignmentServiceMock : assignmentService;
