"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  Bell,
  BellOff,
  Check,
  CheckCheck,
  Trash2,
  BookOpen,
  FileText,
  GraduationCap,
} from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { formatRelativeTime, cn } from "@/lib/helpers";
import notificationService from "@/services/notifications";
import type { Notification } from "@/types";

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.03 },
  },
};

const itemVariants = {
  hidden: { opacity: 0, x: -20 },
  visible: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 20 },
};

const notificationIcons: Record<string, typeof Bell> = {
  info: Bell,
  warning: BookOpen,
  success: GraduationCap,
  error: FileText,
};

function NotificationSkeleton() {
  return (
    <div className="flex items-start gap-4 p-4 rounded-lg border">
      <Skeleton className="h-10 w-10 rounded-full" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-3 w-1/2" />
      </div>
      <Skeleton className="h-8 w-8" />
    </div>
  );
}

function NotificationCard({
  notification,
  onMarkRead,
  onDelete,
}: {
  notification: Notification;
  onMarkRead: (id: string) => void;
  onDelete: (id: string) => void;
}) {
  const Icon = notificationIcons[notification.type] || Bell;

  const iconColors: Record<string, string> = {
    info: "bg-blue-100 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400",
    warning:
      "bg-orange-100 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400",
    success:
      "bg-green-100 dark:bg-green-900/20 text-green-600 dark:text-green-400",
    error: "bg-red-100 dark:bg-red-900/20 text-red-600 dark:text-red-400",
  };

  return (
    <motion.div variants={itemVariants} layout>
      <Card
        className={cn(
          "transition-colors",
          !notification.is_read && "border-primary/50 bg-primary/5",
        )}
      >
        <CardContent className="flex items-start gap-3 sm:gap-4 p-3 sm:p-4">
          <div
            className={cn(
              "flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded-full shrink-0",
              iconColors[notification.type] || iconColors.info,
            )}
          >
            <Icon className="h-4 w-4 sm:h-5 sm:w-5" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <div className="min-w-0 flex-1">
                <h4 className="font-medium text-sm sm:text-base line-clamp-1">
                  {notification.title}
                </h4>
                <p className="text-xs sm:text-sm text-muted-foreground mt-0.5 sm:mt-1 line-clamp-2">
                  {notification.message}
                </p>
              </div>
              {!notification.is_read && (
                <Badge variant="default" className="shrink-0 text-xs">
                  New
                </Badge>
              )}
            </div>
            <p className="text-xs text-muted-foreground mt-1.5 sm:mt-2">
              {formatRelativeTime(notification.created_at)}
            </p>
          </div>
          <div className="flex flex-col sm:flex-row items-center gap-1 shrink-0">
            {!notification.is_read && (
              <Button
                variant="ghost"
                size="icon"
                className="h-9 w-9 sm:h-8 sm:w-8 touch-manipulation"
                onClick={() => onMarkRead(notification.id)}
              >
                <Check className="h-4 w-4" />
              </Button>
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-9 w-9 sm:h-8 sm:w-8 text-muted-foreground hover:text-destructive touch-manipulation"
              onClick={() => onDelete(notification.id)}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}

export default function NotificationsPage() {
  const [activeTab, setActiveTab] = useState<"all" | "unread">("all");
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => notificationService.list(),
  });

  const markReadMutation = useMutation({
    mutationFn: (id: string) => notificationService.markAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
  });

  const markAllReadMutation = useMutation({
    mutationFn: () => notificationService.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      toast.success("All notifications marked as read");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => notificationService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      toast.success("Notification deleted");
    },
  });

  const notifications = data?.notifications || [];
  const unreadCount = notifications.filter((n) => !n.is_read).length;

  const filteredNotifications =
    activeTab === "unread"
      ? notifications.filter((n) => !n.is_read)
      : notifications;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Notifications"
        description={`You have ${unreadCount} unread notification${unreadCount !== 1 ? "s" : ""}`}
        action={
          unreadCount > 0 && (
            <Button
              variant="outline"
              onClick={() => markAllReadMutation.mutate()}
              disabled={markAllReadMutation.isPending}
            >
              <CheckCheck className="mr-2 h-4 w-4" />
              Mark all as read
            </Button>
          )
        }
      />

      <Tabs
        value={activeTab}
        onValueChange={(v) => setActiveTab(v as "all" | "unread")}
      >
        <TabsList>
          <TabsTrigger value="all">
            All
            <Badge variant="secondary" className="ml-2">
              {notifications.length}
            </Badge>
          </TabsTrigger>
          <TabsTrigger value="unread">
            Unread
            {unreadCount > 0 && (
              <Badge variant="default" className="ml-2">
                {unreadCount}
              </Badge>
            )}
          </TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="mt-6">
          {isLoading ? (
            <div className="space-y-3">
              {[1, 2, 3, 4, 5].map((i) => (
                <NotificationSkeleton key={i} />
              ))}
            </div>
          ) : error ? (
            <EmptyState
              title="Error loading notifications"
              description="There was a problem loading your notifications."
              action={
                <Button onClick={() => window.location.reload()}>Retry</Button>
              }
            />
          ) : filteredNotifications.length === 0 ? (
            <EmptyState
              icon={
                activeTab === "unread" ? (
                  <CheckCheck className="h-8 w-8 text-muted-foreground" />
                ) : (
                  <BellOff className="h-8 w-8 text-muted-foreground" />
                )
              }
              title={
                activeTab === "unread" ? "All caught up!" : "No notifications"
              }
              description={
                activeTab === "unread"
                  ? "You've read all your notifications."
                  : "You don't have any notifications yet."
              }
            />
          ) : (
            <motion.div
              variants={containerVariants}
              initial="hidden"
              animate="visible"
              className="space-y-3"
            >
              {filteredNotifications.map((notification) => (
                <NotificationCard
                  key={notification.id}
                  notification={notification}
                  onMarkRead={(id) => markReadMutation.mutate(id)}
                  onDelete={(id) => deleteMutation.mutate(id)}
                />
              ))}
            </motion.div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
