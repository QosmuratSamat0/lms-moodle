"use client";

import { useEffect, useRef, useCallback, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { tokenStorage } from "@/lib/api-client";
import { chatService } from "@/services/chat";
import type { ChatMessage } from "@/types/chat";

const WS_BASE_URL =
  process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1/ws";

// Outbound message types (sent to server)
interface WSOutboundMessage {
  type: "message" | "typing";
  content?: string;
  reply_to_id?: string;
}

// Inbound message types (received from server)
interface WSInboundMessage {
  type: "message" | "typing" | "read" | "error";
  payload?: ChatMessage | TypingPayload | { error: string };
  // Direct message format (alternative)
  message?: ChatMessage;
  user_id?: string;
  user_name?: string;
  is_typing?: boolean;
  error?: string;
}

interface TypingPayload {
  user_id: string;
  user_name: string;
  is_typing: boolean;
}

interface UseChatOptions {
  roomId: string;
  onError?: (error: string) => void;
}

export function useChat({ roomId, onError }: UseChatOptions) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [isConnecting, setIsConnecting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [typingUsers, setTypingUsers] = useState<Map<string, string>>(
    new Map(),
  );

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttempts = useRef(0);
  const shouldReconnectRef = useRef(true);
  const typingTimeoutsRef = useRef<Map<string, NodeJS.Timeout>>(new Map());
  const queryClient = useQueryClient();

  const onErrorRef = useRef(onError);
  useEffect(() => {
    onErrorRef.current = onError;
  }, [onError]);

  // Load initial messages via REST
  useEffect(() => {
    async function loadMessages() {
      try {
        setIsLoading(true);
        const response = await chatService.getMessages(roomId, { limit: 50 });
        // Messages come newest first, reverse for display (oldest first)
        setMessages(response.messages?.reverse() || []);
      } catch (err) {
        console.error("Failed to load messages:", err);
        onErrorRef.current?.((err as Error).message);
      } finally {
        setIsLoading(false);
      }
    }

    if (roomId) {
      loadMessages();
    }
  }, [roomId]);

  // Setup WebSocket connection
  useEffect(() => {
    if (!roomId) {
      console.log("No roomId, skipping WebSocket connection");
      return;
    }

    const initialToken = tokenStorage.getAccessToken();
    if (!initialToken) {
      console.log("No token, skipping WebSocket connection");
      return;
    }

    shouldReconnectRef.current = true;

    const connect = () => {
      if (!shouldReconnectRef.current) return;

      // Get a fresh token each time we connect/reconnect
      const currentToken = tokenStorage.getAccessToken();
      if (!currentToken) {
        console.log("No token available for WebSocket");
        return;
      }

      // Close existing connection
      if (wsRef.current) {
        wsRef.current.close();
      }

      setIsConnecting(true);
      const wsUrl = `${WS_BASE_URL}/chat/${roomId}?token=${currentToken}`;
      console.log("Connecting to WebSocket:", wsUrl);

      try {
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
          console.log("WebSocket connected");
          setIsConnected(true);
          setIsConnecting(false);
          reconnectAttempts.current = 0;
        };

        ws.onmessage = (event) => {
          try {
            const data: WSInboundMessage = JSON.parse(event.data);
            console.log("WebSocket message received:", data);

            // Handle different message formats
            if (data.type === "message") {
              // Support both payload format and direct flat format from WS
              const message: ChatMessage | undefined =
                (data.payload as ChatMessage) ||
                data.message ||
                (data as unknown as { id?: string }).id
                  ? ({
                      id: (data as unknown as Record<string, string>).id || "",
                      room_id:
                        (data as unknown as Record<string, string>).room_id ||
                        "",
                      sender_id:
                        (data as unknown as Record<string, string>).sender_id ||
                        "",
                      sender_first_name: (
                        data as unknown as Record<string, string>
                      ).sender_first_name,
                      sender_last_name: (
                        data as unknown as Record<string, string>
                      ).sender_last_name,
                      sender_email:
                        (data as unknown as Record<string, string>)
                          .sender_email || "",
                      content:
                        (data as unknown as Record<string, string>).content ||
                        "",
                      created_at:
                        (data as unknown as Record<string, string>)
                          .created_at || new Date().toISOString(),
                    } as ChatMessage)
                  : undefined;
              if (message) {
                setMessages((prev) => {
                  // Remove any message with the same id, then add the new one
                  const filtered = prev.filter((m) => m.id !== message.id);
                  return [...filtered, message];
                });
                // Invalidate queries to update room list
                queryClient.invalidateQueries({ queryKey: ["chat-rooms"] });
              }
            } else if (data.type === "typing") {
              const rawData = data as unknown as Record<
                string,
                string | boolean
              >;
              const typing = (data.payload as TypingPayload) || {
                user_id: data.user_id || (rawData.sender_id as string) || "",
                user_name:
                  data.user_name ||
                  (rawData.sender_first_name as string) ||
                  "Someone",
                is_typing: data.is_typing ?? true,
              };

              setTypingUsers((prev) => {
                const next = new Map(prev);
                if (typing.is_typing) {
                  next.set(typing.user_id, typing.user_name);

                  // Clear typing after 3 seconds
                  const existingTimeout = typingTimeoutsRef.current.get(
                    typing.user_id,
                  );
                  if (existingTimeout) clearTimeout(existingTimeout);

                  const timeout = setTimeout(() => {
                    setTypingUsers((p) => {
                      const n = new Map(p);
                      n.delete(typing.user_id);
                      return n;
                    });
                  }, 3000);
                  typingTimeoutsRef.current.set(typing.user_id, timeout);
                } else {
                  next.delete(typing.user_id);
                }
                return next;
              });
            } else if (data.type === "error") {
              const errorMsg =
                (data.payload as { error: string })?.error ||
                data.error ||
                "Unknown error";
              onErrorRef.current?.(errorMsg);
            }
          } catch (e) {
            console.error("Failed to parse WebSocket message:", e);
          }
        };

        ws.onerror = (error) => {
          console.error("WebSocket error:", error);
          setIsConnected(false);
          setIsConnecting(false);
        };

        ws.onclose = (event) => {
          console.log("WebSocket closed:", event.code, event.reason);
          setIsConnected(false);
          setIsConnecting(false);
          wsRef.current = null;

          // Attempt to reconnect with exponential backoff
          if (shouldReconnectRef.current && reconnectAttempts.current < 5) {
            const delay = Math.min(
              1000 * Math.pow(2, reconnectAttempts.current),
              30000,
            );
            console.log(
              `Reconnecting in ${delay}ms (attempt ${reconnectAttempts.current + 1})`,
            );

            reconnectTimeoutRef.current = setTimeout(() => {
              reconnectAttempts.current++;
              connect();
            }, delay);
          }
        };
      } catch (error) {
        console.error("Failed to create WebSocket:", error);
        setIsConnecting(false);
      }
    };

    connect();

    // Copy refs for cleanup
    const reconnectTimeout = reconnectTimeoutRef.current;
    const ws = wsRef.current;
    const typingTimeouts = typingTimeoutsRef.current;

    return () => {
      shouldReconnectRef.current = false;
      if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
      }
      if (ws) {
        ws.close();
      }
      // Clear typing timeouts
      typingTimeouts.forEach((timeout) => clearTimeout(timeout));
      typingTimeouts.clear();
    };
  }, [roomId, queryClient]);

  // Send message (prefer WebSocket, fallback to REST)
  const sendMessage = useCallback(
    async (content: string, replyToId?: string) => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        const message: WSOutboundMessage = {
          type: "message",
          content,
          reply_to_id: replyToId,
        };
        try {
          wsRef.current.send(JSON.stringify(message));
          return true;
        } catch (error) {
          console.error("Failed to send WebSocket message:", error);
        }
      }

      // Fallback to REST API ONLY if not connected to WS
      if (!isConnected) {
        try {
          const message = await chatService.sendMessage(roomId, {
            content,
            reply_to_id: replyToId,
          });
          setMessages((prev) => [...prev, message]);
          return true;
        } catch (err) {
          console.error("Failed to send message via REST:", err);
          onErrorRef.current?.((err as Error).message);
          return false;
        }
      }
    },
    [roomId],
  );

  const sendTyping = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      try {
        wsRef.current.send(JSON.stringify({ type: "typing" }));
      } catch (error) {
        console.error("Failed to send typing indicator:", error);
      }
    }
  }, []);

  const disconnect = useCallback(() => {
    shouldReconnectRef.current = false;
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setIsConnected(false);
    reconnectAttempts.current = 0;
  }, []);

  const reconnect = useCallback(() => {
    disconnect();
    shouldReconnectRef.current = true;
    reconnectAttempts.current = 0;
  }, [disconnect]);

  return {
    messages,
    isConnected,
    isConnecting,
    isLoading,
    typingUsers: Array.from(typingUsers.values()),
    sendMessage,
    sendTyping,
    reconnect,
    disconnect,
  };
}

// Keep old export for backward compatibility
export const useChatWebSocket = useChat;
export default useChat;
