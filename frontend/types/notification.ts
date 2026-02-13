export type NotificationType = "info" | "warning" | "success" | "error";
export type NotificationChannel = "in_app" | "email" | "sms" | "push";

export interface Notification {
  id: string;
  user_id: string;
  title: string;
  message: string;
  type: NotificationType;
  channel: NotificationChannel;
  is_read: boolean;
  read_at?: string;
  created_at: string;
}

export interface NotificationListResponse {
  notifications: Notification[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface UnreadCountResponse {
  count: number;
}

export interface NotificationFilters {
  type?: NotificationType;
  is_read?: boolean;
}
