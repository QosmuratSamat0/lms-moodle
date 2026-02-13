export type EventType = "class" | "lab" | "exam" | "office_hours" | "event";
export type RecurrenceType = "none" | "daily" | "weekly" | "monthly";

export interface ScheduleEvent {
  id: string;
  course_id?: string;
  course_title?: string;
  title: string;
  description?: string;
  event_type: EventType;
  start_time: string;
  end_time: string;
  location?: string;
  is_online?: boolean;
  meeting_url?: string;
  recurrence: RecurrenceType;
  recur_until?: string;
  creator_name?: string;
  creator_email?: string;
  created_by_id?: string;
  created_at: string;
}

export interface CreateEventRequest {
  course_id?: string;
  title: string;
  description?: string;
  event_type: EventType;
  start_time: string;
  end_time: string;
  location?: string;
  is_online?: boolean;
  meeting_url?: string;
  recurrence?: RecurrenceType;
  recur_until?: string;
}

export interface UpdateEventRequest {
  title?: string;
  description?: string;
  event_type?: EventType;
  start_time?: string;
  end_time?: string;
  location?: string;
  is_online?: boolean;
  meeting_url?: string;
  recurrence?: RecurrenceType;
  recur_until?: string;
}

export interface ScheduleListResponse {
  events: ScheduleEvent[];
  total: number;
  page?: number;
  limit?: number;
  total_pages?: number;
}

export interface CalendarQuery {
  start_date: string;
  end_date: string;
  course_id?: string;
}
