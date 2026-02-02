import api from "@/lib/api-client";
import type {
  ChatRoom,
  ChatMessage,
  ChatMember,
  RoomListResponse,
  MessageListResponse,
  CreateRoomRequest,
  SendMessageRequest,
  UpdateMessageRequest,
} from "@/types/chat";
import type { PaginationParams } from "@/types/common";

const WS_BASE_URL =
  process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1/ws";

type QueryParams = Record<string, string | number | boolean | undefined>;

export const chatService = {
  // Rooms
  getRooms: async (params?: PaginationParams): Promise<RoomListResponse> => {
    return api.get<RoomListResponse>("/chat/rooms", params as QueryParams);
  },

  getRoom: async (id: string): Promise<ChatRoom> => {
    return api.get<ChatRoom>(`/chat/rooms/${id}`);
  },

  createRoom: async (data: CreateRoomRequest): Promise<ChatRoom> => {
    return api.post<ChatRoom>("/chat/rooms", data);
  },

  updateRoom: async (
    id: string,
    data: Partial<CreateRoomRequest>,
  ): Promise<ChatRoom> => {
    return api.put<ChatRoom>(`/chat/rooms/${id}`, data);
  },

  deleteRoom: async (id: string): Promise<void> => {
    return api.delete(`/chat/rooms/${id}`);
  },

  getOrCreateDM: async (userId: string): Promise<ChatRoom> => {
    return api.post<ChatRoom>("/chat/rooms/direct", { user_id: userId });
  },

  // Messages
  getMessages: async (
    roomId: string,
    params?: PaginationParams,
  ): Promise<MessageListResponse> => {
    return api.get<MessageListResponse>(
      `/chat/rooms/${roomId}/messages`,
      params as QueryParams,
    );
  },

  sendMessage: async (
    roomId: string,
    data: SendMessageRequest,
  ): Promise<ChatMessage> => {
    return api.post<ChatMessage>(`/chat/rooms/${roomId}/messages`, data);
  },

  updateMessage: async (
    roomId: string,
    messageId: string,
    data: UpdateMessageRequest,
  ): Promise<ChatMessage> => {
    return api.put<ChatMessage>(
      `/chat/rooms/${roomId}/messages/${messageId}`,
      data,
    );
  },

  deleteMessage: async (roomId: string, messageId: string): Promise<void> => {
    return api.delete(`/chat/rooms/${roomId}/messages/${messageId}`);
  },

  // Members
  getMembers: async (roomId: string): Promise<{ members: ChatMember[] }> => {
    return api.get(`/chat/rooms/${roomId}/members`);
  },

  addMember: async (roomId: string, userId: string): Promise<ChatMember> => {
    return api.post<ChatMember>(`/chat/rooms/${roomId}/members`, {
      user_id: userId,
    });
  },

  removeMember: async (roomId: string, userId: string): Promise<void> => {
    return api.delete(`/chat/rooms/${roomId}/members/${userId}`);
  },

  leaveRoom: async (roomId: string): Promise<void> => {
    return api.post(`/chat/rooms/${roomId}/leave`);
  },

  // WebSocket URL builder - room ID is now a path param, token goes in header when connecting
  getWebSocketUrl: (roomId: string, token: string): string => {
    return `${WS_BASE_URL}/chat/${roomId}?token=${token}`;
  },
};

// Mock data
const mockRooms: ChatRoom[] = [
  {
    id: "1",
    name: "CS101 General Discussion",
    type: "course",
    course_id: "1",
    course_title: "Introduction to Computer Science",
    member_count: 25,
    participants: [
      {
        id: "2",
        first_name: "Alice",
        last_name: "Smith",
        email: "alice@example.com",
      },
    ],
    last_message: {
      id: "m1",
      content: "Has anyone started the assignment?",
      sender_id: "2",
      sender_first_name: "Alice",
      created_at: "2025-01-25T14:30:00Z",
    },
    last_message_at: "2025-01-25T14:30:00Z",
    unread_count: 3,
    created_at: "2025-01-01T00:00:00Z",
  },
  {
    id: "2",
    name: "Study Group Alpha",
    type: "group",
    member_count: 5,
    participants: [
      {
        id: "3",
        first_name: "Bob",
        last_name: "Jones",
        email: "bob@example.com",
      },
    ],
    last_message: {
      id: "m2",
      content: "Meeting tomorrow at 3pm",
      sender_id: "3",
      sender_first_name: "Bob",
      created_at: "2025-01-25T12:00:00Z",
    },
    last_message_at: "2025-01-25T12:00:00Z",
    unread_count: 0,
    created_at: "2025-01-10T00:00:00Z",
  },
  {
    id: "3",
    type: "direct",
    member_count: 2,
    participants: [
      {
        id: "4",
        first_name: "Dr. Sarah",
        last_name: "Johnson",
        email: "sarah@example.com",
      },
    ],
    last_message: {
      id: "m3",
      content: "Thanks for the help!",
      sender_id: "4",
      sender_first_name: "Dr. Sarah",
      created_at: "2025-01-24T18:45:00Z",
    },
    last_message_at: "2025-01-24T18:45:00Z",
    unread_count: 1,
    created_at: "2025-01-15T00:00:00Z",
  },
];

const mockMessages: Record<string, ChatMessage[]> = {
  "1": [
    {
      id: "1",
      room_id: "1",
      sender_id: "2",
      sender_first_name: "Alice",
      sender_last_name: "Smith",
      sender_email: "alice@example.com",
      content: "Hey everyone! Welcome to CS101.",
      created_at: "2025-01-25T10:00:00Z",
    },
    {
      id: "2",
      room_id: "1",
      sender_id: "3",
      sender_first_name: "Bob",
      sender_last_name: "Johnson",
      sender_email: "bob@example.com",
      content: "Thanks! Excited to be here.",
      created_at: "2025-01-25T10:05:00Z",
    },
    {
      id: "3",
      room_id: "1",
      sender_id: "1",
      sender_first_name: "John",
      sender_last_name: "Doe",
      sender_email: "john@example.com",
      content: "Has anyone started the assignment?",
      created_at: "2025-01-25T14:30:00Z",
    },
  ],
  "2": [
    {
      id: "4",
      room_id: "2",
      sender_id: "4",
      sender_first_name: "Carol",
      sender_last_name: "White",
      sender_email: "carol@example.com",
      content: "Should we meet this week?",
      created_at: "2025-01-25T11:00:00Z",
    },
    {
      id: "5",
      room_id: "2",
      sender_id: "1",
      sender_first_name: "John",
      sender_last_name: "Doe",
      sender_email: "john@example.com",
      content: "Meeting tomorrow at 3pm",
      created_at: "2025-01-25T12:00:00Z",
    },
  ],
  "3": [
    {
      id: "6",
      room_id: "3",
      sender_id: "5",
      sender_first_name: "Dr. Sarah",
      sender_last_name: "Johnson",
      sender_email: "sarah@example.com",
      content: "Let me know if you have any questions about the lecture.",
      created_at: "2025-01-24T18:30:00Z",
    },
    {
      id: "7",
      room_id: "3",
      sender_id: "1",
      sender_first_name: "John",
      sender_last_name: "Doe",
      sender_email: "john@example.com",
      content: "Thanks for the help!",
      created_at: "2025-01-24T18:45:00Z",
    },
  ],
};

export const chatServiceMock = {
  getRooms: async (params?: PaginationParams): Promise<RoomListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    return {
      rooms: mockRooms,
      total: mockRooms.length,
      page: params?.page || 1,
      limit: params?.limit || 10,
      total_pages: 1,
    };
  },

  getRoom: async (id: string): Promise<ChatRoom> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    const room = mockRooms.find((r) => r.id === id);
    if (!room) throw new Error("Room not found");
    return room;
  },

  createRoom: async (data: CreateRoomRequest): Promise<ChatRoom> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: String(Date.now()),
      name: data.name,
      type: data.type,
      course_id: data.course_id,
      member_count: (data.members?.length || 0) + 1,
      unread_count: 0,
      created_at: new Date().toISOString(),
    };
  },

  updateRoom: async (
    id: string,
    data: Partial<CreateRoomRequest>,
  ): Promise<ChatRoom> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    const room = mockRooms.find((r) => r.id === id);
    if (!room) throw new Error("Room not found");
    return { ...room, ...data };
  },

  deleteRoom: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  getOrCreateDM: async (): Promise<ChatRoom> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    return mockRooms[2];
  },

  getMessages: async (
    roomId: string,
    params?: PaginationParams,
  ): Promise<MessageListResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    const messages = mockMessages[roomId] || [];
    return {
      messages,
      total: messages.length,
      page: params?.page || 1,
      limit: params?.limit || 50,
      total_pages: 1,
    };
  },

  sendMessage: async (
    roomId: string,
    data: SendMessageRequest,
  ): Promise<ChatMessage> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return {
      id: String(Date.now()),
      room_id: roomId,
      sender_id: "1",
      sender_first_name: "John",
      sender_last_name: "Doe",
      sender_email: "john@example.com",
      content: data.content,
      reply_to_id: data.reply_to_id,
      created_at: new Date().toISOString(),
    };
  },

  updateMessage: async (
    roomId: string,
    messageId: string,
    data: UpdateMessageRequest,
  ): Promise<ChatMessage> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return {
      id: messageId,
      room_id: roomId,
      sender_id: "1",
      sender_first_name: "John",
      sender_last_name: "Doe",
      sender_email: "john@example.com",
      content: data.content,
      edited_at: new Date().toISOString(),
      created_at: new Date().toISOString(),
    };
  },

  deleteMessage: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  getMembers: async (): Promise<{ members: ChatMember[] }> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return { members: [] };
  },

  addMember: async (roomId: string, userId: string): Promise<ChatMember> => {
    await new Promise((resolve) => setTimeout(resolve, 400));
    return {
      id: String(Date.now()),
      room_id: roomId,
      user_id: userId,
      role: "member",
      email: "member@example.com",
      joined_at: new Date().toISOString(),
    };
  },

  removeMember: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  leaveRoom: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
  },

  getWebSocketUrl: (roomId: string, token: string): string => {
    return `${WS_BASE_URL}/chat?room_id=${roomId}&token=${token}`;
  },
};

const useMock = process.env.NEXT_PUBLIC_USE_MOCK === "true";
export default useMock ? chatServiceMock : chatService;
