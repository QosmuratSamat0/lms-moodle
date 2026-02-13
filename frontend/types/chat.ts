export type RoomType = "direct" | "group" | "course";
export type MemberRole = "admin" | "member";
export type MessageType = "text" | "image" | "file" | "system";

export interface ChatParticipant {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  avatar_url?: string;
}

export interface LastMessage {
  id: string;
  content: string;
  sender_id: string;
  sender_first_name?: string;
  created_at: string;
}

export interface ChatRoom {
  id: string;
  name?: string;
  type: RoomType;
  course_id?: string;
  course_title?: string;
  member_count: number;
  participants?: ChatParticipant[];
  last_message?: LastMessage;
  last_message_at?: string;
  unread_count: number;
  created_at: string;
}

export interface MessageSender {
  id: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
}

export interface ChatMessage {
  id: string;
  room_id: string;
  sender_id: string;
  sender_first_name?: string;
  sender_last_name?: string;
  sender_email: string;
  sender?: MessageSender;
  content: string;
  message_type?: MessageType;
  file_url?: string;
  reply_to_id?: string;
  reply_to_content?: string;
  edited_at?: string;
  is_read?: boolean;
  created_at: string;
}

export interface ChatMember {
  id: string;
  room_id: string;
  user_id: string;
  role: MemberRole;
  first_name?: string;
  last_name?: string;
  email: string;
  joined_at: string;
}

export interface CreateRoomRequest {
  name?: string;
  type: RoomType;
  course_id?: string;
  members?: string[];
  participant_ids?: string[];
}

export interface SendMessageRequest {
  content: string;
  message_type?: MessageType;
  reply_to_id?: string;
  file_url?: string;
}

export interface UpdateMessageRequest {
  content: string;
}

export interface RoomListResponse {
  rooms: ChatRoom[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface MessageListResponse {
  messages: ChatMessage[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

// WebSocket message types
export type WSMessageType = "message" | "typing" | "read" | "join" | "leave";

export interface WSMessage {
  type: WSMessageType;
  payload: unknown;
}

export interface WSNewMessage {
  room_id: string;
  message: ChatMessage;
}

export interface WSTyping {
  room_id: string;
  user_id: string;
  user_name: string;
  is_typing: boolean;
}
