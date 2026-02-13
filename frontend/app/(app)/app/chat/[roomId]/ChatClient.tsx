"use client";

import { useState, useRef, useEffect, useMemo } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { toast } from "sonner";
import {
  ArrowLeft,
  Send,
  MoreVertical,
  Info,
  Paperclip,
  Smile,
  Image as ImageIcon,
  Users,
  Check,
  CheckCheck,
  Wifi,
  WifiOff,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { useAuthStore } from "@/store/auth-store";
import { formatTime, getInitials, cn } from "@/lib/helpers";
import chatService from "@/services/chat";
import { useChat } from "@/hooks/use-chat-websocket";
import type { ChatMessage, ChatParticipant } from "@/types";

const messageVariants = {
  hidden: { opacity: 0, y: 10 },
  visible: { opacity: 1, y: 0 },
};

function MessageBubble({
  message,
  isOwn,
  showAvatar,
}: {
  message: ChatMessage;
  isOwn: boolean;
  showAvatar: boolean;
}) {
  return (
    <motion.div
      variants={messageVariants}
      initial="hidden"
      animate="visible"
      className={cn("flex gap-2", isOwn ? "flex-row-reverse" : "flex-row")}
    >
      {showAvatar ? (
        <Avatar className="h-8 w-8 mt-1">
          <AvatarImage src={message.sender?.avatar_url} />
          <AvatarFallback className="text-xs">
            {getInitials(
              message.sender_first_name || message.sender?.first_name || "",
              message.sender_last_name || message.sender?.last_name || "",
            )}
          </AvatarFallback>
        </Avatar>
      ) : (
        <div className="w-8" />
      )}
      <div className={cn("flex flex-col max-w-[70%]", isOwn && "items-end")}>
        {showAvatar && !isOwn && (
          <span className="text-xs text-muted-foreground mb-1">
            {message.sender_first_name || message.sender?.first_name}
          </span>
        )}
        <div
          className={cn(
            "px-4 py-2 rounded-2xl",
            isOwn
              ? "bg-primary text-primary-foreground rounded-br-sm"
              : "bg-muted rounded-bl-sm",
          )}
        >
          {message.message_type === "image" && message.file_url && (
            /* eslint-disable-next-line @next/next/no-img-element */
            <img
              src={message.file_url}
              alt="Shared image"
              className="rounded-lg max-w-full mb-2"
            />
          )}
          {message.message_type === "file" && message.file_url && (
            <a
              href={message.file_url}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-sm underline"
            >
              <Paperclip className="h-4 w-4" />
              Attached file
            </a>
          )}
          {message.content && <p className="text-sm">{message.content}</p>}
        </div>
        <div className="flex items-center gap-1 mt-1">
          <span className="text-xs text-muted-foreground">
            {formatTime(message.created_at)}
          </span>
          {isOwn && (
            <span className="text-xs text-muted-foreground">
              {message.is_read ? (
                <CheckCheck className="h-3 w-3 text-primary" />
              ) : (
                <Check className="h-3 w-3" />
              )}
            </span>
          )}
        </div>
      </div>
    </motion.div>
  );
}

function ChatHeader({
  room,
  roomId,
  isLoading,
}: {
  room?: { name?: string; type: string; participants?: ChatParticipant[] };
  roomId: string;
  isLoading: boolean;
}) {
  if (isLoading) {
    return (
      <div className="flex items-center gap-4 p-4 border-b">
        <Skeleton className="h-10 w-10 rounded-full" />
        <div className="space-y-2">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-20" />
        </div>
      </div>
    );
  }

  const isGroup = room?.type === "group";
  const participant = room?.participants?.[0];

  return (
    <div className="flex items-center justify-between p-4 border-b bg-card">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" asChild className="lg:hidden">
          <Link href="/app/chat">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <Avatar className="h-10 w-10">
          {isGroup ? (
            <AvatarFallback className="bg-primary/10 text-primary">
              <Users className="h-5 w-5" />
            </AvatarFallback>
          ) : (
            <>
              <AvatarImage src={participant?.avatar_url} />
              <AvatarFallback>
                {getInitials(
                  participant?.first_name || "",
                  participant?.last_name || "",
                )}
              </AvatarFallback>
            </>
          )}
        </Avatar>
        <div>
          <h2 className="font-semibold">
            {room?.name ||
              (participant
                ? `${participant.first_name} ${participant.last_name}`
                : "Chat")}
          </h2>
          <p className="text-xs text-muted-foreground">
            {isGroup ? `${room?.participants?.length || 0} members` : "Online"}
          </p>
        </div>
      </div>
      <div className="flex items-center gap-1">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon">
              <Info className="h-4 w-4" />
            </Button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>Chat Info</SheetTitle>
              <SheetDescription>
                {isGroup ? "Group details" : "User details"}
              </SheetDescription>
            </SheetHeader>
            <div className="mt-6 space-y-6">
              <div className="flex flex-col items-center gap-3">
                <Avatar className="h-20 w-20">
                  {isGroup ? (
                    <AvatarFallback className="bg-primary/10 text-primary text-2xl">
                      <Users className="h-10 w-10" />
                    </AvatarFallback>
                  ) : (
                    <>
                      <AvatarImage src={participant?.avatar_url} />
                      <AvatarFallback className="text-2xl">
                        {getInitials(
                          participant?.first_name || "",
                          participant?.last_name || "",
                        )}
                      </AvatarFallback>
                    </>
                  )}
                </Avatar>
                <div className="text-center">
                  <h3 className="font-semibold text-lg">
                    {room?.name ||
                      (participant
                        ? `${participant.first_name} ${participant.last_name}`
                        : "Chat")}
                  </h3>
                  {participant?.email && (
                    <p className="text-sm text-muted-foreground">
                      {participant.email}
                    </p>
                  )}
                </div>
              </div>
              {isGroup && room?.participants && (
                <div>
                  <h4 className="font-medium mb-3">Members</h4>
                  <div className="space-y-2">
                    {room.participants.map((p) => (
                      <div
                        key={p.id}
                        className="flex items-center gap-3 p-2 rounded hover:bg-muted"
                      >
                        <Avatar className="h-8 w-8">
                          <AvatarImage src={p.avatar_url} />
                          <AvatarFallback className="text-xs">
                            {getInitials(p.first_name, p.last_name)}
                          </AvatarFallback>
                        </Avatar>
                        <span className="text-sm">
                          {p.first_name} {p.last_name}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </SheetContent>
        </Sheet>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem>Search in conversation</DropdownMenuItem>
            <DropdownMenuItem>Mute notifications</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="text-destructive"
              onClick={async () => {
                if (
                  confirm("Are you sure you want to delete this conversation?")
                ) {
                  try {
                    await chatService.deleteRoom(roomId);
                    window.location.href = "/app/chat";
                  } catch (err) {
                    toast.error("Failed to delete conversation");
                  }
                }
              }}
            >
              Delete conversation
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}

export default function ChatClient({ roomId }: { roomId: string }) {
  const { user } = useAuthStore();

  const [messageInput, setMessageInput] = useState("");
  const [isSending, setIsSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const typingTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Use the chat hook (handles WebSocket + REST fallback internally)
  const {
    messages,
    isConnected,
    isConnecting,
    isLoading: messagesLoading,
    typingUsers,
    sendMessage,
    sendTyping,
  } = useChat({
    roomId,
    onError: (error) => toast.error(error),
  });

  const { data: room, isLoading: roomLoading } = useQuery({
    queryKey: ["chat-room", roomId],
    queryFn: () => chatService.getRoom(roomId),
  });

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Send typing indicator with debounce
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMessageInput(e.target.value);

    if (isConnected && e.target.value) {
      if (typingTimeoutRef.current) {
        clearTimeout(typingTimeoutRef.current);
      }
      sendTyping();
      typingTimeoutRef.current = setTimeout(() => {
        typingTimeoutRef.current = null;
      }, 2000);
    }
  };

  const handleSend = async () => {
    if (!messageInput.trim() || isSending) return;

    setIsSending(true);
    const success = await sendMessage(messageInput);
    setIsSending(false);

    if (success) {
      setMessageInput("");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  // Group messages by sender for avatar display
  const groupedMessages = useMemo(() => {
    // Deduplicate by id
    const seen = new Set();
    const deduped = messages.filter((msg) => {
      if (seen.has(msg.id)) return false;
      seen.add(msg.id);
      return true;
    });
    return deduped.map((msg, i) => {
      const prevMsg = deduped[i - 1];
      const showAvatar =
        !prevMsg ||
        prevMsg.sender_id !== msg.sender_id ||
        new Date(msg.created_at).getTime() -
          new Date(prevMsg.created_at).getTime() >
          5 * 60 * 1000; // 5 minutes
      return { ...msg, showAvatar };
    });
  }, [messages]);

  return (
    <div className="flex flex-col h-[calc(100vh-8rem)] -m-6 border rounded-lg overflow-hidden bg-background">
      <ChatHeader room={room || undefined} roomId={roomId} isLoading={roomLoading} />

      {/* Messages Area */}
      <ScrollArea className="flex-1 p-4">
        {messagesLoading ? (
          <div className="space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div
                key={i}
                className={cn(
                  "flex gap-2",
                  i % 2 === 0 ? "flex-row-reverse" : "",
                )}
              >
                <Skeleton className="h-8 w-8 rounded-full" />
                <Skeleton className="h-12 w-48 rounded-2xl" />
              </div>
            ))}
          </div>
        ) : messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <div className="h-16 w-16 rounded-full bg-muted flex items-center justify-center mb-4">
              <Send className="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 className="font-medium">No messages yet</h3>
            <p className="text-sm text-muted-foreground mt-1">
              Start the conversation by sending a message below.
            </p>
          </div>
        ) : (
          <div className="space-y-3">
            <AnimatePresence mode="popLayout">
              {groupedMessages.map((msg) => (
                <MessageBubble
                  key={msg.id}
                  message={msg}
                  isOwn={msg.sender_id === user?.id}
                  showAvatar={msg.showAvatar}
                />
              ))}
            </AnimatePresence>
            <div ref={messagesEndRef} />
          </div>
        )}
      </ScrollArea>

      {/* Input Area */}
      <div className="p-4 border-t bg-card">
        {/* Connection status & typing indicator */}
        <div className="flex items-center justify-between mb-2 min-h-5">
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            {isConnecting ? (
              <>
                <div className="h-2 w-2 rounded-full bg-yellow-500 animate-pulse" />
                <span>Connecting...</span>
              </>
            ) : isConnected ? (
              <>
                <Wifi className="h-3 w-3 text-green-500" />
                <span className="text-green-600">Live</span>
              </>
            ) : (
              <>
                <WifiOff className="h-3 w-3 text-red-500" />
                <span className="text-red-600">Offline (polling)</span>
              </>
            )}
          </div>
          {typingUsers.length > 0 && (
            <span className="text-xs text-muted-foreground italic">
              Someone is typing...
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <Paperclip className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start">
              <DropdownMenuItem>
                <ImageIcon className="h-4 w-4 mr-2" />
                Photo or Video
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Paperclip className="h-4 w-4 mr-2" />
                File
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          {/* Emoji icon removed */}
          <Input
            ref={inputRef}
            placeholder="Type a message..."
            value={messageInput}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            disabled={isSending}
            className="flex-1"
          />
          <Button
            onClick={handleSend}
            disabled={!messageInput.trim() || isSending}
            size="icon"
          >
            <Send className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
