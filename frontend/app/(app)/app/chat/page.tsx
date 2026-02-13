"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { MessageSquare, Plus, Users, MoreVertical } from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import { SearchInput } from "@/components/common/search-input";
import { EmptyState } from "@/components/common/empty-state";
import { NewMessageDialog } from "@/components/chat";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { formatRelativeTime, getInitials } from "@/lib/helpers";
import chatService from "@/services/chat";
import type { ChatRoom } from "@/types";

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.05 },
  },
};

const itemVariants = {
  hidden: { opacity: 0, x: -20 },
  visible: { opacity: 1, x: 0 },
};

function RoomSkeleton() {
  return (
    <div className="flex items-center gap-4 p-4 rounded-lg border">
      <Skeleton className="h-12 w-12 rounded-full" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-32" />
        <Skeleton className="h-3 w-48" />
      </div>
      <Skeleton className="h-4 w-16" />
    </div>
  );
}

function RoomCard({ room }: { room: ChatRoom }) {
  const isGroup = room.type === "group";

  return (
    <motion.div variants={itemVariants}>
      <Link href={`/app/chat/${room.id}`}>
        <Card className="hover:border-primary/50 transition-colors">
          <CardContent className="flex items-center gap-4 p-4">
            <Avatar className="h-12 w-12">
              {isGroup ? (
                <AvatarFallback className="bg-primary/10 text-primary">
                  <Users className="h-5 w-5" />
                </AvatarFallback>
              ) : (
                <>
                  <AvatarImage src={room.participants?.[0]?.avatar_url} />
                  <AvatarFallback>
                    {getInitials(
                      room.participants?.[0]?.first_name || "",
                      room.participants?.[0]?.last_name || "",
                    )}
                  </AvatarFallback>
                </>
              )}
            </Avatar>

            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <h3 className="font-medium truncate">
                  {room.name ||
                    (room.participants?.[0]
                      ? `${room.participants[0].first_name} ${room.participants[0].last_name}`
                      : "Chat")}
                </h3>
                {room.unread_count > 0 && (
                  <Badge variant="default" className="h-5 min-w-[20px] px-1.5">
                    {room.unread_count}
                  </Badge>
                )}
              </div>
              {room.last_message && (
                <p className="text-sm text-muted-foreground truncate">
                  {room.last_message.content}
                </p>
              )}
            </div>

            <div className="flex flex-col items-end gap-2">
              {room.last_message && (
                <span className="text-xs text-muted-foreground">
                  {formatRelativeTime(room.last_message.created_at)}
                </span>
              )}
              <DropdownMenu>
                <DropdownMenuTrigger
                  asChild
                  onClick={(e) => e.preventDefault()}
                >
                  <Button variant="ghost" size="icon" className="h-8 w-8">
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem>Mark as read</DropdownMenuItem>
                  <DropdownMenuItem>Mute notifications</DropdownMenuItem>
                  {isGroup && <DropdownMenuItem>Leave group</DropdownMenuItem>}
                  <DropdownMenuItem className="text-destructive">
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </CardContent>
        </Card>
      </Link>
    </motion.div>
  );
}

export default function ChatPage() {
  const [search, setSearch] = useState("");
  const [activeFilter, setActiveFilter] = useState<"all" | "direct" | "group">(
    "all",
  );

  const { data, isLoading, error } = useQuery({
    queryKey: ["chat-rooms"],
    queryFn: () => chatService.getRooms(),
  });

  const rooms = data?.rooms || [];

  const filteredRooms = rooms.filter((room) => {
    // Filter by type
    if (activeFilter === "direct" && room.type !== "direct") return false;
    if (activeFilter === "group" && room.type !== "group") return false;

    // Filter by search
    if (search) {
      const searchLower = search.toLowerCase();
      const roomName = room.name?.toLowerCase() || "";
      const participantNames =
        room.participants
          ?.map((p: { first_name: string; last_name: string }) =>
            `${p.first_name} ${p.last_name}`.toLowerCase(),
          )
          .join(" ") || "";
      return (
        roomName.includes(searchLower) || participantNames.includes(searchLower)
      );
    }

    return true;
  });

  // Sort by last message time
  const sortedRooms = [...filteredRooms].sort((a, b) => {
    const aTime = a.last_message?.created_at || a.created_at;
    const bTime = b.last_message?.created_at || b.created_at;
    return new Date(bTime).getTime() - new Date(aTime).getTime();
  });

  return (
    <div className="space-y-6">
      <PageHeader
        title="Messages"
        description="Chat with classmates and instructors"
        action={<NewMessageDialog />}
      />

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
        <SearchInput
          placeholder="Search conversations..."
          value={search}
          onChange={setSearch}
          className="sm:max-w-sm"
        />
        <div className="flex gap-2">
          <Button
            variant={activeFilter === "all" ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveFilter("all")}
          >
            All
          </Button>
          <Button
            variant={activeFilter === "direct" ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveFilter("direct")}
          >
            Direct
          </Button>
          <Button
            variant={activeFilter === "group" ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveFilter("group")}
          >
            Groups
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3, 4, 5].map((i) => (
            <RoomSkeleton key={i} />
          ))}
        </div>
      ) : error ? (
        <EmptyState
          title="Error loading messages"
          description="There was a problem loading your conversations."
          action={
            <Button onClick={() => window.location.reload()}>Retry</Button>
          }
        />
      ) : sortedRooms.length === 0 ? (
        <EmptyState
          icon={<MessageSquare className="h-8 w-8 text-muted-foreground" />}
          title="No conversations"
          description={
            search
              ? "No conversations match your search."
              : "Start a new conversation to connect with others."
          }
          action={<NewMessageDialog />}
        />
      ) : (
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="space-y-3"
        >
          {sortedRooms.map((room) => (
            <RoomCard key={room.id} room={room} />
          ))}
        </motion.div>
      )}
    </div>
  );
}
