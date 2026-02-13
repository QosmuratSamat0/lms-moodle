import api from "@/lib/api-client";
import type {
  ScheduleEvent,
  ScheduleListResponse,
  CreateEventRequest,
  UpdateEventRequest,
  CalendarQuery,
} from "@/types/schedule";

type QueryParams = Record<string, string | number | boolean | undefined>;

export const scheduleService = {
  getMySchedule: async (): Promise<ScheduleListResponse> => {
    return api.get<ScheduleListResponse>("/schedule");
  },

  getCalendar: async (params: CalendarQuery): Promise<ScheduleListResponse> => {
    return api.get<ScheduleListResponse>(
      "/schedule/calendar",
      params as unknown as QueryParams,
    );
  },

  get: async (id: string): Promise<ScheduleEvent> => {
    return api.get<ScheduleEvent>(`/schedule/events/${id}`);
  },

  getCourseSchedule: async (
    courseId: string,
  ): Promise<ScheduleListResponse> => {
    return api.get<ScheduleListResponse>(`/courses/${courseId}/schedule`);
  },

  create: async (data: CreateEventRequest): Promise<ScheduleEvent> => {
    return api.post<ScheduleEvent>("/schedule/events", data);
  },

  update: async (
    id: string,
    data: UpdateEventRequest,
  ): Promise<ScheduleEvent> => {
    return api.put<ScheduleEvent>(`/schedule/events/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    return api.delete(`/schedule/events/${id}`);
  },
};

// Mock data
const mockEvents: ScheduleEvent[] = [
  {
    id: "1",
    course_id: "1",
    course_title: "Introduction to Computer Science",
    title: "CS101 Lecture",
    description: "Introduction to programming concepts",
    event_type: "class",
    start_time: "2026-01-27T09:00:00Z",
    end_time: "2026-01-27T10:30:00Z",
    location: "Room 101",
    recurrence: "weekly",
    created_at: "2025-01-01T00:00:00Z",
  },
  {
    id: "2",
    course_id: "2",
    course_title: "Advanced Mathematics",
    title: "Math Lecture",
    description: "Calculus review",
    event_type: "class",
    start_time: "2026-01-27T11:00:00Z",
    end_time: "2026-01-27T12:30:00Z",
    location: "Room 205",
    recurrence: "weekly",
    created_at: "2025-01-01T00:00:00Z",
  },
  {
    id: "3",
    course_id: "1",
    course_title: "Introduction to Computer Science",
    title: "CS101 Lab",
    description: "Hands-on programming practice",
    event_type: "lab",
    start_time: "2026-01-28T14:00:00Z",
    end_time: "2026-01-28T16:00:00Z",
    location: "Computer Lab B",
    recurrence: "weekly",
    created_at: "2025-01-01T00:00:00Z",
  },
  {
    id: "4",
    course_id: "1",
    course_title: "Introduction to Computer Science",
    title: "CS101 Midterm Exam",
    description: "Midterm examination",
    event_type: "exam",
    start_time: "2026-02-15T09:00:00Z",
    end_time: "2026-02-15T11:00:00Z",
    location: "Examination Hall A",
    recurrence: "none",
    created_at: "2025-01-01T00:00:00Z",
  },
  {
    id: "5",
    title: "Office Hours - Dr. Johnson",
    description: "Open office hours for questions",
    event_type: "office_hours",
    start_time: "2026-01-29T15:00:00Z",
    end_time: "2026-01-29T17:00:00Z",
    location: "Office 302",
    recurrence: "weekly",
    created_at: "2025-01-01T00:00:00Z",
  },
];

export const scheduleServiceMock = {
  getMySchedule: async (): Promise<ScheduleListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    return {
      events: mockEvents,
      total: mockEvents.length,
    };
  },

  getCalendar: async (params: CalendarQuery): Promise<ScheduleListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));

    const startDate = new Date(params.start_date);
    const endDate = new Date(params.end_date);

    let filtered = mockEvents.filter((e) => {
      const eventDate = new Date(e.start_time);
      return eventDate >= startDate && eventDate <= endDate;
    });

    if (params.course_id) {
      filtered = filtered.filter((e) => e.course_id === params.course_id);
    }

    return {
      events: filtered,
      total: filtered.length,
    };
  },

  get: async (id: string): Promise<ScheduleEvent> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    const event = mockEvents.find((e) => e.id === id);
    if (!event) throw new Error("Event not found");
    return event;
  },

  getCourseSchedule: async (
    courseId: string,
  ): Promise<ScheduleListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    const filtered = mockEvents.filter((e) => e.course_id === courseId);
    return {
      events: filtered,
      total: filtered.length,
    };
  },

  create: async (data: CreateEventRequest): Promise<ScheduleEvent> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: String(Date.now()),
      ...data,
      recurrence: data.recurrence || "none",
      created_at: new Date().toISOString(),
    };
  },

  update: async (
    id: string,
    data: UpdateEventRequest,
  ): Promise<ScheduleEvent> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    const event = mockEvents.find((e) => e.id === id);
    if (!event) throw new Error("Event not found");
    return { ...event, ...data };
  },

  delete: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },
};

const useMock = process.env.NEXT_PUBLIC_USE_MOCK === "true";
export default useMock ? scheduleServiceMock : scheduleService;
