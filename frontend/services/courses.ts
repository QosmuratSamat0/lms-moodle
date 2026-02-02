import api from "@/lib/api-client";
import type {
  Course,
  CourseListResponse,
  CourseFilters,
  CreateCourseRequest,
  UpdateCourseRequest,
  Enrollment,
  EnrollRequest,
  UpdateEnrollmentRequest,
} from "@/types/course";
import type { PaginationParams } from "@/types/common";

type QueryParams = Record<string, string | number | boolean | undefined>;

export const courseService = {
  list: async (
    params?: PaginationParams & CourseFilters,
  ): Promise<CourseListResponse> => {
    return api.get<CourseListResponse>("/courses", params as QueryParams);
  },

  get: async (id: string): Promise<Course> => {
    return api.get<Course>(`/courses/${id}`);
  },

  create: async (data: CreateCourseRequest): Promise<Course> => {
    return api.post<Course>("/courses", data);
  },

  update: async (id: string, data: UpdateCourseRequest): Promise<Course> => {
    return api.put<Course>(`/courses/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    return api.delete(`/courses/${id}`);
  },

  // Enrollments
  enroll: async (data: EnrollRequest): Promise<Enrollment> => {
    return api.post<Enrollment>("/enrollments", data);
  },

  dropEnrollment: async (courseId: string): Promise<void> => {
    return api.delete(`/enrollments/${courseId}`);
  },

  getMyEnrollments: async (): Promise<{ enrollments: Enrollment[] }> => {
    return api.get("/enrollments/my");
  },

  getCourseEnrollments: async (
    courseId: string,
  ): Promise<{ enrollments: Enrollment[] }> => {
    return api.get(`/courses/${courseId}/enrollments`);
  },

  updateEnrollmentStatus: async (
    enrollmentId: string,
    data: UpdateEnrollmentRequest,
  ): Promise<Enrollment> => {
    return api.put<Enrollment>(`/enrollments/${enrollmentId}`, data);
  },
};

// Mock data for development
const mockCourses: Course[] = [
  {
    id: "1",
    title: "Introduction to Computer Science",
    description:
      "Learn the fundamentals of programming and computer science concepts.",
    owner_teacher_id: "t1",
    teacher_first_name: "Dr. Sarah",
    teacher_last_name: "Johnson",
    is_active: true,
    created_at: "2025-01-01T00:00:00Z",
    updated_at: "2025-01-15T00:00:00Z",
  },
  {
    id: "2",
    title: "Advanced Mathematics",
    description: "Calculus, Linear Algebra, and Differential Equations.",
    owner_teacher_id: "t2",
    teacher_first_name: "Prof. Michael",
    teacher_last_name: "Chen",
    is_active: true,
    created_at: "2025-01-05T00:00:00Z",
    updated_at: "2025-01-20T00:00:00Z",
  },
  {
    id: "3",
    title: "Data Structures and Algorithms",
    description: "Master essential data structures and algorithmic techniques.",
    owner_teacher_id: "t1",
    teacher_first_name: "Dr. Sarah",
    teacher_last_name: "Johnson",
    is_active: true,
    created_at: "2025-01-10T00:00:00Z",
    updated_at: "2025-01-22T00:00:00Z",
  },
  {
    id: "4",
    title: "Web Development Fundamentals",
    description: "HTML, CSS, JavaScript, and modern web frameworks.",
    owner_teacher_id: "t3",
    teacher_first_name: "Emily",
    teacher_last_name: "Williams",
    is_active: true,
    created_at: "2025-01-12T00:00:00Z",
    updated_at: "2025-01-23T00:00:00Z",
  },
];

export const courseServiceMock = {
  list: async (
    params?: PaginationParams & CourseFilters,
  ): Promise<CourseListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 500));

    let filtered = [...mockCourses];

    if (params?.search) {
      const search = params.search.toLowerCase();
      filtered = filtered.filter(
        (c) =>
          c.title.toLowerCase().includes(search) ||
          c.description?.toLowerCase().includes(search),
      );
    }

    return {
      courses: filtered,
      total: filtered.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: Math.ceil(filtered.length / (params?.limit || 10)),
    };
  },

  get: async (id: string): Promise<Course> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    const course = mockCourses.find((c) => c.id === id);
    if (!course) throw new Error("Course not found");
    return course;
  },

  create: async (data: CreateCourseRequest): Promise<Course> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: String(Date.now()),
      ...data,
      is_active: true,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
  },

  update: async (id: string, data: UpdateCourseRequest): Promise<Course> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    const course = mockCourses.find((c) => c.id === id);
    if (!course) throw new Error("Course not found");
    return { ...course, ...data, updated_at: new Date().toISOString() };
  },

  delete: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  enroll: async (data: EnrollRequest): Promise<Enrollment> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: String(Date.now()),
      course_id: data.course_id,
      student_id: "1",
      status: "active",
      enrolled_at: new Date().toISOString(),
    };
  },

  getMyEnrollments: async (): Promise<{ enrollments: Enrollment[] }> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return {
      enrollments: [
        {
          id: "1",
          course_id: "1",
          student_id: "1",
          status: "active",
          enrolled_at: "2025-01-15T00:00:00Z",
          course_title: "Introduction to Computer Science",
        },
        {
          id: "2",
          course_id: "2",
          student_id: "1",
          status: "active",
          enrolled_at: "2025-01-16T00:00:00Z",
          course_title: "Advanced Mathematics",
        },
      ],
    };
  },

  dropEnrollment: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  getCourseEnrollments: async (): Promise<{ enrollments: Enrollment[] }> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return { enrollments: [] };
  },

  updateEnrollmentStatus: async (
    enrollmentId: string,
    data: UpdateEnrollmentRequest,
  ): Promise<Enrollment> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return {
      id: enrollmentId,
      course_id: "1",
      student_id: "1",
      status: data.status,
      enrolled_at: new Date().toISOString(),
    };
  },
};

const useMock = process.env.NEXT_PUBLIC_USE_MOCK === "true";
export default useMock ? courseServiceMock : courseService;
