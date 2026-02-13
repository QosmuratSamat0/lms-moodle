import api from "@/lib/api-client";
import type {
  Notification,
  NotificationListResponse,
  NotificationFilters,
  UnreadCountResponse,
} from "@/types/notification";
import type { PaginationParams } from "@/types/common";

type QueryParams = Record<string, string | number | boolean | undefined>;

export const notificationService = {
  list: async (
    params?: PaginationParams & NotificationFilters,
  ): Promise<NotificationListResponse> => {
    // Backend uses skip/take, not page/limit
    const backendParams: Record<string, string | number | boolean | undefined> =
      {};
    if (params?.limit) backendParams.take = params.limit;
    if (params?.page && params?.limit)
      backendParams.skip = (params.page - 1) * params.limit;
    const notifications = await api.get<any>("/notifications", backendParams);
    const list = Array.isArray(notifications)
      ? notifications
      : notifications?.notifications || [];
    return {
      notifications: list,
      total: list.length,
      unread_count: list.filter((n: any) => !n.is_read && !n.read).length,
    };
  },

  get: async (id: string): Promise<Notification> => {
    return api.get<Notification>(`/notifications/${id}`);
  },

  getUnreadCount: async (): Promise<UnreadCountResponse> => {
    return api.get<UnreadCountResponse>("/notifications/unread-count");
  },

  markAsRead: async (id: string): Promise<void> => {
    return api.patch(`/notifications/${id}/read`);
  },

  markAllAsRead: async (): Promise<void> => {
    return api.post("/notifications/read-all");
  },

  delete: async (id: string): Promise<void> => {
    return api.delete(`/notifications/${id}`);
  },

  deleteAll: async (): Promise<void> => {
    return api.delete("/notifications");
  },
};

// Mock data
const mockNotifications: Notification[] = [
  {
    id: "1",
    user_id: "1",
    title: "Assignment Due Soon",
    message: "Programming Assignment 1 is due in 2 days.",
    type: "warning",
    channel: "in_app",
    is_read: false,
    created_at: "2025-01-25T10:00:00Z",
  },
  {
    id: "2",
    user_id: "1",
    title: "Grade Posted",
    message: "Your grade for Calculus Problem Set 1 has been posted.",
    type: "success",
    channel: "in_app",
    is_read: false,
    created_at: "2025-01-24T15:30:00Z",
  },
  {
    id: "3",
    user_id: "1",
    title: "New Course Material",
    message: "New lecture slides have been uploaded for CS101.",
    type: "info",
    channel: "in_app",
    is_read: true,
    read_at: "2025-01-24T12:00:00Z",
    created_at: "2025-01-24T10:00:00Z",
  },
  {
    id: "4",
    user_id: "1",
    title: "Enrollment Approved",
    message: "You have been enrolled in Data Structures and Algorithms.",
    type: "success",
    channel: "in_app",
    is_read: true,
    read_at: "2025-01-23T09:00:00Z",
    created_at: "2025-01-23T08:00:00Z",
  },
  {
    id: "5",
    user_id: "1",
    title: "System Maintenance",
    message: "The system will be undergoing maintenance on Sunday.",
    type: "info",
    channel: "in_app",
    is_read: true,
    read_at: "2025-01-22T14:00:00Z",
    created_at: "2025-01-22T12:00:00Z",
  },
];

export const notificationServiceMock = {
  list: async (
    params?: PaginationParams & NotificationFilters,
  ): Promise<NotificationListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));

    let filtered = [...mockNotifications];

    if (params?.type) {
      filtered = filtered.filter((n) => n.type === params.type);
    }

    if (params?.is_read !== undefined) {
      filtered = filtered.filter((n) => n.is_read === params.is_read);
    }

    return {
      notifications: filtered,
      total: filtered.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: Math.ceil(filtered.length / (params?.limit || 10)),
    };
  },

  get: async (id: string): Promise<Notification> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    const notification = mockNotifications.find((n) => n.id === id);
    if (!notification) throw new Error("Notification not found");
    return notification;
  },

  getUnreadCount: async (): Promise<UnreadCountResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 200));
    const unread = mockNotifications.filter((n) => !n.is_read).length;
    return { count: unread };
  },

  markAsRead: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 200));
  },

  markAllAsRead: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  delete: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 200));
  },

  deleteAll: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },
};

const useMock = process.env.NEXT_PUBLIC_USE_MOCK === "true";
export default useMock ? notificationServiceMock : notificationService;
