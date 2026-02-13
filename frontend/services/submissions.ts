import api from "@/lib/api-client";
import type {
  Submission,
  SubmissionListResponse,
  CreateSubmissionRequest,
  UpdateSubmissionRequest,
} from "@/types/submission";
import type {
  Grade,
  GradeSubmissionRequest,
  GradeListResponse,
} from "@/types/grade";
import type { PaginationParams } from "@/types/common";

type QueryParams = Record<string, string | number | boolean | undefined>;

export const submissionService = {
  // Submissions
  get: async (id: string): Promise<Submission> => {
    return api.get<Submission>(`/submissions/${id}`);
  },

  create: async (data: CreateSubmissionRequest): Promise<Submission> => {
    return api.post<Submission>("/submissions", data);
  },

  update: async (
    id: string,
    data: UpdateSubmissionRequest,
  ): Promise<Submission> => {
    return api.put<Submission>(`/submissions/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    return api.delete(`/submissions/${id}`);
  },

  getMySubmissions: async (
    studentId: string,
    params?: PaginationParams,
  ): Promise<SubmissionListResponse> => {
    const submissions = await api.get<any>(
      `/submissions/student/${studentId}`,
      params as QueryParams,
    );
    const list = Array.isArray(submissions)
      ? submissions
      : submissions?.submissions || [];
    return { submissions: list, total: list.length };
  },

  getByAssignment: async (
    assignmentId: string,
    params?: PaginationParams & { group_id?: string },
  ): Promise<SubmissionListResponse> => {
    const submissions = await api.get<any>(
      `/submissions/assignment/${assignmentId}`,
      params as QueryParams,
    );
    const list = Array.isArray(submissions)
      ? submissions
      : submissions?.submissions || [];
    return { submissions: list, total: list.length };
  },

  // Grades
  gradeSubmission: async (data: GradeSubmissionRequest): Promise<Grade> => {
    return api.post<Grade>("/grades", data);
  },

  getGrade: async (id: string): Promise<Grade> => {
    return api.get<Grade>(`/grades/${id}`);
  },

  getMyGrades: async (
    params?: PaginationParams,
  ): Promise<GradeListResponse> => {
    return api.get<GradeListResponse>("/grades/me", params as QueryParams);
  },

  getCourseGrades: async (
    courseId: string,
    params?: PaginationParams,
  ): Promise<GradeListResponse> => {
    return api.get<GradeListResponse>(
      `/courses/${courseId}/grades`,
      params as QueryParams,
    );
  },
};

// Mock data
const mockSubmissions: Submission[] = [
  {
    id: "1",
    assignment_id: "1",
    student_id: "1",
    content_text: 'print("Hello, World!")',
    status: "graded",
    submitted_at: "2025-01-28T10:30:00Z",
    assignment_title: "Programming Assignment 1: Hello World",
    course_title: "Introduction to Computer Science",
    student_first_name: "John",
    student_last_name: "Doe",
    student_email: "john.doe@example.com",
    max_points: 100,
    due_at: "2026-02-01T23:59:59Z",
    grade_score: 95,
    grade_feedback: "Great work! Clean code.",
  },
  {
    id: "2",
    assignment_id: "2",
    student_id: "1",
    content_text: 'x = 10\ny = "hello"\nz = 3.14',
    status: "submitted",
    submitted_at: "2025-01-29T14:20:00Z",
    assignment_title: "Programming Assignment 2: Variables and Types",
    course_title: "Introduction to Computer Science",
    student_first_name: "John",
    student_last_name: "Doe",
    student_email: "john.doe@example.com",
    max_points: 100,
    due_at: "2026-02-15T23:59:59Z",
  },
];

const mockGrades: Grade[] = [
  {
    id: "1",
    submission_id: "1",
    student_id: "1",
    assignment_id: "1",
    score: 95,
    max_points: 100,
    percentage: 95,
    feedback: "Great work! Clean code.",
    graded_by_teacher_id: "t1",
    teacher_first_name: "Dr. Sarah",
    teacher_last_name: "Johnson",
    assignment_title: "Programming Assignment 1: Hello World",
    course_title: "Introduction to Computer Science",
    student_first_name: "John",
    student_last_name: "Doe",
    student_email: "john.doe@example.com",
    graded_at: "2025-01-29T09:00:00Z",
  },
];

export const submissionServiceMock = {
  get: async (id: string): Promise<Submission> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    const submission = mockSubmissions.find((s) => s.id === id);
    if (!submission) throw new Error("Submission not found");
    return submission;
  },

  create: async (data: CreateSubmissionRequest): Promise<Submission> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: String(Date.now()),
      assignment_id: data.assignment_id,
      student_id: "1",
      content_text: data.content_text,
      file_url: data.file_url,
      status: "submitted",
      submitted_at: new Date().toISOString(),
    };
  },

  update: async (
    id: string,
    data: UpdateSubmissionRequest,
  ): Promise<Submission> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    const submission = mockSubmissions.find((s) => s.id === id);
    if (!submission) throw new Error("Submission not found");
    return { ...submission, ...data };
  },

  delete: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  getMySubmissions: async (
    params?: PaginationParams,
  ): Promise<SubmissionListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    return {
      submissions: mockSubmissions,
      total: mockSubmissions.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: 1,
    };
  },

  getByAssignment: async (
    assignmentId: string,
    params?: PaginationParams,
  ): Promise<SubmissionListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    const filtered = mockSubmissions.filter(
      (s) => s.assignment_id === assignmentId,
    );
    return {
      submissions: filtered,
      total: filtered.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: 1,
    };
  },

  gradeSubmission: async (data: GradeSubmissionRequest): Promise<Grade> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: String(Date.now()),
      submission_id: data.submission_id,
      score: data.score,
      max_points: 100,
      percentage: data.score,
      feedback: data.feedback,
      graded_at: new Date().toISOString(),
    };
  },

  getGrade: async (id: string): Promise<Grade> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    const grade = mockGrades.find((g) => g.id === id);
    if (!grade) throw new Error("Grade not found");
    return grade;
  },

  getMyGrades: async (
    params?: PaginationParams,
  ): Promise<GradeListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    return {
      grades: mockGrades,
      total: mockGrades.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: 1,
    };
  },

  getCourseGrades: async (
    _courseId: string,
    params?: PaginationParams,
  ): Promise<GradeListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    return {
      grades: mockGrades,
      total: mockGrades.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: 1,
    };
  },
};

const useMock = process.env.NEXT_PUBLIC_USE_MOCK === "true";
export default (useMock ? submissionServiceMock : submissionService) as typeof submissionService;
