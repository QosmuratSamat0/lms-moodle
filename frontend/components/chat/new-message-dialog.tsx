"use client";

import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Plus, Loader2, Search, User, Users } from "lucide-react";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Checkbox } from "@/components/ui/checkbox";
import { cn } from "@/lib/utils";
import { getInitials } from "@/lib/helpers";
import chatService from "@/services/chat";
import courseService from "@/services/courses";

interface NewMessageDialogProps {
  trigger?: React.ReactNode;
}

export function NewMessageDialog({ trigger }: NewMessageDialogProps) {
  const [open, setOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<"direct" | "group">("direct");
  const [search, setSearch] = useState("");
  const [selectedUsers, setSelectedUsers] = useState<string[]>([]);
  const [groupName, setGroupName] = useState("");
  const router = useRouter();
  const queryClient = useQueryClient();

  // Fetch enrolled courses to get classmates/teachers
  const { data: coursesData, isLoading: coursesLoading } = useQuery({
    queryKey: ["courses", "enrolled"],
    queryFn: () => courseService.list({ limit: 100 }),
    enabled: open,
  });

  // Get unique participants from courses (teachers)
  const availableUsers: {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    role: string;
  }[] = [];
  const seenIds = new Set<string>();

  coursesData?.courses?.forEach((course) => {
    const teacherId = course.teacher_id || course.owner_teacher_id;
    if (teacherId && !seenIds.has(teacherId)) {
      seenIds.add(teacherId);
      availableUsers.push({
        id: teacherId,
        first_name: course.teacher_first_name || "Teacher",
        last_name: course.teacher_last_name || "",
        email: "",
        role: "teacher",
      });
    }
  });

  const filteredUsers = availableUsers.filter((user) => {
    if (!search) return true;
    const fullName = `${user.first_name} ${user.last_name}`.toLowerCase();
    return fullName.includes(search.toLowerCase());
  });

  // Create direct message
  const dmMutation = useMutation({
    mutationFn: (userId: string) => chatService.getOrCreateDM(userId),
    onSuccess: (room) => {
      queryClient.invalidateQueries({ queryKey: ["chat-rooms"] });
      setOpen(false);
      router.push(`/app/chat/${room.id}`);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to start conversation");
    },
  });

  // Create group chat
  const groupMutation = useMutation({
    mutationFn: () =>
      chatService.createRoom({
        name: groupName,
        type: "group",
        participant_ids: selectedUsers,
      }),
    onSuccess: (room) => {
      toast.success("Group created successfully!");
      queryClient.invalidateQueries({ queryKey: ["chat-rooms"] });
      setOpen(false);
      setGroupName("");
      setSelectedUsers([]);
      router.push(`/app/chat/${room.id}`);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create group");
    },
  });

  const handleUserSelect = (userId: string) => {
    if (activeTab === "direct") {
      dmMutation.mutate(userId);
    } else {
      setSelectedUsers((prev) =>
        prev.includes(userId)
          ? prev.filter((id) => id !== userId)
          : [...prev, userId],
      );
    }
  };

  const handleGroupSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (groupName.trim() && selectedUsers.length > 0) {
      groupMutation.mutate();
    }
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(v) => {
        setOpen(v);
        if (!v) {
          setSearch("");
          setSelectedUsers([]);
          setGroupName("");
        }
      }}
    >
      <DialogTrigger asChild>
        {trigger || (
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            New Message
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>New Conversation</DialogTitle>
          <DialogDescription>
            Start a direct message or create a group chat
          </DialogDescription>
        </DialogHeader>

        <Tabs
          value={activeTab}
          onValueChange={(v) => {
            setActiveTab(v as "direct" | "group");
            setSelectedUsers([]);
            setGroupName("");
          }}
        >
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="direct" className="gap-2">
              <User className="h-4 w-4" />
              Direct Message
            </TabsTrigger>
            <TabsTrigger value="group" className="gap-2">
              <Users className="h-4 w-4" />
              Group Chat
            </TabsTrigger>
          </TabsList>

          <div className="mt-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search people..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          <TabsContent value="direct" className="mt-4">
            <ScrollArea className="h-75">
              {coursesLoading ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : filteredUsers.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <p>No users found</p>
                  <p className="text-sm mt-1">
                    {search
                      ? "Try a different search"
                      : "Join a course to chat with teachers"}
                  </p>
                </div>
              ) : (
                <div className="space-y-1">
                  {filteredUsers.map((user) => (
                    <button
                      key={user.id}
                      onClick={() => handleUserSelect(user.id)}
                      disabled={dmMutation.isPending}
                      className={cn(
                        "w-full flex items-center gap-3 p-3 rounded-lg hover:bg-muted transition-colors text-left",
                        dmMutation.isPending && "opacity-50 cursor-not-allowed",
                      )}
                    >
                      <Avatar className="h-10 w-10">
                        <AvatarFallback>
                          {getInitials(user.first_name, user.last_name)}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium truncate">
                          {user.first_name} {user.last_name}
                        </p>
                        <p className="text-sm text-muted-foreground capitalize">
                          {user.role}
                        </p>
                      </div>
                      {dmMutation.isPending && (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      )}
                    </button>
                  ))}
                </div>
              )}
            </ScrollArea>
          </TabsContent>

          <TabsContent value="group" className="mt-4">
            <form onSubmit={handleGroupSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="groupName">Group Name</Label>
                <Input
                  id="groupName"
                  placeholder="Study Group..."
                  value={groupName}
                  onChange={(e) => setGroupName(e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label>Members</Label>
                <p className="text-sm text-muted-foreground">
                  Select people to add to the group
                </p>
                <ScrollArea className="h-50 border rounded-md p-2">
                  {filteredUsers.length === 0 ? (
                    <div className="text-center py-4 text-muted-foreground text-sm">
                      No users found
                    </div>
                  ) : (
                    <div className="space-y-1">
                      {filteredUsers.map((user) => (
                        <label
                          key={user.id}
                          className={cn(
                            "flex items-center gap-3 p-2 rounded-lg hover:bg-muted cursor-pointer transition-colors",
                            selectedUsers.includes(user.id) && "bg-muted",
                          )}
                        >
                          <Checkbox
                            checked={selectedUsers.includes(user.id)}
                            onCheckedChange={() => handleUserSelect(user.id)}
                          />
                          <Avatar className="h-8 w-8">
                            <AvatarFallback className="text-xs">
                              {getInitials(user.first_name, user.last_name)}
                            </AvatarFallback>
                          </Avatar>
                          <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium truncate">
                              {user.first_name} {user.last_name}
                            </p>
                          </div>
                        </label>
                      ))}
                    </div>
                  )}
                </ScrollArea>
              </div>

              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setOpen(false)}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={
                    groupMutation.isPending ||
                    selectedUsers.length === 0 ||
                    !groupName.trim()
                  }
                >
                  {groupMutation.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Create Group
                </Button>
              </div>
            </form>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}
