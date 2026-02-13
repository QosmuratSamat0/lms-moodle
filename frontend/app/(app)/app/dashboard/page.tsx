"use client";

import { useMemo } from "react";
import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import {
  BookOpen,
  Calendar,
  Bell,
  ArrowRight,
  Clock,
  FileText,
  MessageCircle,
  Send,
  ExternalLink,
} from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/store/auth-store";
import {
  formatRelativeTime,
  getFullName,
  isOverdue,
  getDaysUntil,
} from "@/lib/helpers";
import courseService from "@/services/courses";
import assignmentService from "@/services/assignments";
import notificationService from "@/services/notifications";
import { useTeacherProfile } from "@/hooks/use-profile";

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.1 },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
};

export default function DashboardPage() {
  const { user } = useAuthStore();
  const isTeacher = user?.role === "teacher";

  // Fetch teacher profile to get the teacher's profile ID
  const { data: teacherProfile } = useTeacherProfile();

  // Fetch student enrollments to filter courses
  const { data: enrollmentsData, isLoading: enrollmentsLoading } = useQuery({
    queryKey: ["enrollments", "student", user?.id],
    queryFn: () => courseService.getMyEnrollments(user!.id),
    enabled: !isTeacher && !!user?.id,
  });

  // Get set of enrolled course IDs for students
  const enrolledCourseIds = useMemo(() => {
    const enrollments = Array.isArray(enrollmentsData)
      ? enrollmentsData
      : enrollmentsData?.enrollments || [];
    return new Set(enrollments.map((e: { course_id: string }) => e.course_id));
  }, [enrollmentsData]);

  const { data: coursesData, isLoading: coursesLoading } = useQuery({
    queryKey: ["courses", "recent"],
    queryFn: () => courseService.list({ limit: 20 }),
  });

  const { data: notificationsData, isLoading: notificationsLoading } = useQuery(
    {
      queryKey: ["notifications", "recent"],
      queryFn: () => notificationService.list({ limit: 5 }),
    },
  );

  // Filter courses based on role
  const filteredCourses = useMemo(() => {
    const courses = Array.isArray(coursesData)
      ? coursesData
      : (coursesData as any)?.courses || [];
    if (isTeacher) {
      // Filter by teacher's profile ID (owner_teacher_id)
      if (teacherProfile?.id) {
        return courses.filter(
          (course: any) => course.owner_teacher_id === teacherProfile.id,
        );
      }
      return [];
    }
    // For students, filter by enrolled courses
    if (enrolledCourseIds.size > 0) {
      return courses.filter((course: any) => enrolledCourseIds.has(course.id));
    }
    return [];
  }, [coursesData, isTeacher, teacherProfile, enrolledCourseIds]);

  // Get course IDs for fetching assignments
  const dashboardCourseIds = useMemo(() => {
    return filteredCourses.map((c: any) => c.id);
  }, [filteredCourses]);

  // Fetch assignments for each relevant course
  const { data: assignmentsData, isLoading: assignmentsLoading } = useQuery({
    queryKey: ["dashboard-assignments", dashboardCourseIds],
    queryFn: async () => {
      if (dashboardCourseIds.length === 0) return { assignments: [], total: 0 };
      const results = await Promise.all(
        dashboardCourseIds.map((courseId: string) =>
          assignmentService.getByCourse(courseId, { limit: 100 }),
        ),
      );
      const allAssignments = results.flatMap((r) => r.assignments || []);
      return { assignments: allAssignments, total: allAssignments.length };
    },
    enabled: dashboardCourseIds.length > 0,
  });

  // Filter assignments based on role
  const filteredAssignments = useMemo(() => {
    return assignmentsData?.assignments || [];
  }, [assignmentsData?.assignments]);

  const loading =
    coursesLoading || (!isTeacher && enrollmentsLoading) || assignmentsLoading;

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="space-y-6"
    >
      <PageHeader
        title={`Welcome back, ${getFullName(user?.first_name, user?.last_name)}`}
        description="Here's what's happening with your courses today."
      />

      {/* Stats Grid */}
      <motion.div
        variants={itemVariants}
        className="grid gap-3 sm:gap-4 grid-cols-2 lg:grid-cols-4"
      >
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-3 sm:p-6 sm:pb-2">
            <CardTitle className="text-xs sm:text-sm font-medium">
              {isTeacher ? "Assigned Courses" : "Enrolled Courses"}
            </CardTitle>
            <BookOpen className="h-4 w-4 text-muted-foreground hidden sm:block" />
          </CardHeader>
          <CardContent className="p-3 pt-0 sm:p-6 sm:pt-0">
            <div className="text-xl sm:text-2xl font-bold">
              {loading ? (
                <Skeleton className="h-6 sm:h-8 w-12 sm:w-16" />
              ) : (
                filteredCourses.length
              )}
            </div>
            <p className="text-xs text-muted-foreground">Active courses</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-3 sm:p-6 sm:pb-2">
            <CardTitle className="text-xs sm:text-sm font-medium">
              Upcoming Deadlines
            </CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground hidden sm:block" />
          </CardHeader>
          <CardContent className="p-3 pt-0 sm:p-6 sm:pt-0">
            <div className="text-xl sm:text-2xl font-bold">
              {loading || assignmentsLoading ? (
                <Skeleton className="h-6 sm:h-8 w-12 sm:w-16" />
              ) : (
                filteredAssignments.filter(
                  (a) => a.due_at && !isOverdue(a.due_at),
                ).length || 0
              )}
            </div>
            <p className="text-xs text-muted-foreground">This week</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-3 sm:p-6 sm:pb-2">
            <CardTitle className="text-xs sm:text-sm font-medium">
              {isTeacher ? "Total Assignments" : "Pending Tasks"}
            </CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground hidden sm:block" />
          </CardHeader>
          <CardContent className="p-3 pt-0 sm:p-6 sm:pt-0">
            <div className="text-xl sm:text-2xl font-bold">
              {loading || assignmentsLoading ? (
                <Skeleton className="h-6 sm:h-8 w-12 sm:w-16" />
              ) : (
                filteredAssignments.length
              )}
            </div>
            <p className="text-xs text-muted-foreground">To complete</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-3 sm:p-6 sm:pb-2">
            <CardTitle className="text-xs sm:text-sm font-medium">
              Notifications
            </CardTitle>
            <Bell className="h-4 w-4 text-muted-foreground hidden sm:block" />
          </CardHeader>
          <CardContent className="p-3 pt-0 sm:p-6 sm:pt-0">
            <div className="text-xl sm:text-2xl font-bold">
              {notificationsLoading ? (
                <Skeleton className="h-6 sm:h-8 w-12 sm:w-16" />
              ) : (
                (notificationsData?.notifications ?? []).filter(
                  (n) => !n.is_read,
                ).length || 0
              )}
            </div>
            <p className="text-xs text-muted-foreground">Unread</p>
          </CardContent>
        </Card>
      </motion.div>

      <div className="grid gap-4 sm:gap-6 lg:grid-cols-2">
        {/* Recent Courses */}
        <motion.div variants={itemVariants}>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between p-4 sm:p-6">
              <div>
                <CardTitle className="text-base sm:text-lg">
                  {isTeacher ? "Your Courses" : "Recent Courses"}
                </CardTitle>
                <CardDescription className="text-xs sm:text-sm">
                  {isTeacher
                    ? "Courses you are assigned to"
                    : "Your enrolled courses"}
                </CardDescription>
              </div>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="touch-manipulation"
              >
                <Link href="/app/courses">
                  <span className="hidden sm:inline">View all</span>
                  <ArrowRight className="sm:ml-2 h-4 w-4" />
                </Link>
              </Button>
            </CardHeader>
            <CardContent className="p-4 pt-0 sm:p-6 sm:pt-0">
              {loading ? (
                <div className="space-y-3 sm:space-y-4">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="flex items-center gap-3 sm:gap-4">
                      <Skeleton className="h-10 w-10 sm:h-12 sm:w-12 rounded" />
                      <div className="space-y-1.5 sm:space-y-2 flex-1">
                        <Skeleton className="h-3.5 sm:h-4 w-3/4" />
                        <Skeleton className="h-2.5 sm:h-3 w-1/2" />
                      </div>
                    </div>
                  ))}
                </div>
              ) : filteredCourses.length > 0 ? (
                <div className="space-y-2 sm:space-y-4">
                  {filteredCourses.slice(0, 4).map((course) => (
                    <Link
                      key={course.id}
                      href={`/app/courses/${course.id}`}
                      className="flex items-center gap-3 sm:gap-4 rounded-lg p-2 transition-colors hover:bg-muted touch-manipulation"
                    >
                      <div className="flex h-10 w-10 sm:h-12 sm:w-12 items-center justify-center rounded bg-primary/10 shrink-0">
                        <BookOpen className="h-5 w-5 sm:h-6 sm:w-6 text-primary" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium truncate text-sm sm:text-base">
                          {course.title}
                        </p>
                        <p className="text-xs sm:text-sm text-muted-foreground truncate">
                          {getFullName(
                            course.teacher_first_name,
                            course.teacher_last_name,
                          )}
                        </p>
                      </div>
                    </Link>
                  ))}
                </div>
              ) : (
                <p className="text-center text-muted-foreground py-6 sm:py-8 text-sm">
                  No courses yet
                </p>
              )}
            </CardContent>
          </Card>
        </motion.div>

        {/* Upcoming Deadlines */}
        <motion.div variants={itemVariants}>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between p-4 sm:p-6">
              <div>
                <CardTitle className="text-base sm:text-lg">
                  {isTeacher ? "Recent Assignments" : "Upcoming Deadlines"}
                </CardTitle>
                <CardDescription className="text-xs sm:text-sm">
                  {isTeacher
                    ? "Assignments in your courses"
                    : "Assignments due soon"}
                </CardDescription>
              </div>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="touch-manipulation"
              >
                <Link href="/app/assignments">
                  <span className="hidden sm:inline">View all</span>
                  <ArrowRight className="sm:ml-2 h-4 w-4" />
                </Link>
              </Button>
            </CardHeader>
            <CardContent className="p-4 pt-0 sm:p-6 sm:pt-0">
              {loading || assignmentsLoading ? (
                <div className="space-y-3 sm:space-y-4">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="flex items-center gap-3 sm:gap-4">
                      <Skeleton className="h-9 w-9 sm:h-10 sm:w-10 rounded" />
                      <div className="space-y-1.5 sm:space-y-2 flex-1">
                        <Skeleton className="h-4 w-3/4" />
                        <Skeleton className="h-3 w-1/2" />
                      </div>
                    </div>
                  ))}
                </div>
              ) : filteredAssignments.length > 0 ? (
                <div className="space-y-3">
                  {filteredAssignments.slice(0, 5).map((assignment) => {
                    const daysUntil = assignment.due_at
                      ? getDaysUntil(assignment.due_at)
                      : null;
                    const overdue =
                      assignment.due_at && isOverdue(assignment.due_at);

                    return (
                      <Link
                        key={assignment.id}
                        href={`/app/courses/${assignment.course_id}/assignments/${assignment.id}`}
                        className="flex items-center gap-3 sm:gap-4 rounded-lg p-2 transition-colors hover:bg-muted touch-manipulation"
                      >
                        <div className="flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded bg-orange-100 dark:bg-orange-900/20 shrink-0">
                          <Clock className="h-4 w-4 sm:h-5 sm:w-5 text-orange-600 dark:text-orange-400" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="font-medium truncate text-sm sm:text-base">
                            {assignment.title}
                          </p>
                          <p className="text-xs sm:text-sm text-muted-foreground truncate">
                            {assignment.course_title}
                          </p>
                        </div>
                        <Badge
                          variant={
                            overdue
                              ? "destructive"
                              : daysUntil && daysUntil <= 3
                                ? "secondary"
                                : "outline"
                          }
                          className="shrink-0 text-xs"
                        >
                          {overdue
                            ? "Overdue"
                            : daysUntil !== null
                              ? daysUntil === 0
                                ? "Today"
                                : daysUntil === 1
                                  ? "Tomorrow"
                                  : `${daysUntil}d`
                              : "No date"}
                        </Badge>
                      </Link>
                    );
                  })}
                </div>
              ) : (
                <p className="text-center text-muted-foreground py-6 sm:py-8 text-sm">
                  {isTeacher ? "No assignments yet" : "No upcoming deadlines"}
                </p>
              )}
            </CardContent>
          </Card>
        </motion.div>
      </div>

      {/* Recent Notifications */}
      <motion.div variants={itemVariants}>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between p-4 sm:p-6">
            <div>
              <CardTitle className="text-base sm:text-lg">
                Recent Notifications
              </CardTitle>
              <CardDescription className="text-xs sm:text-sm">
                Stay updated with your activities
              </CardDescription>
            </div>
            <Button
              variant="ghost"
              size="sm"
              asChild
              className="touch-manipulation"
            >
              <Link href="/app/notifications">
                <span className="hidden sm:inline">View all</span>
                <ArrowRight className="sm:ml-2 h-4 w-4" />
              </Link>
            </Button>
          </CardHeader>
          <CardContent className="p-4 pt-0 sm:p-6 sm:pt-0">
            {notificationsLoading ? (
              <div className="space-y-3 sm:space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="flex items-start gap-3 sm:gap-4">
                    <Skeleton className="h-7 w-7 sm:h-8 sm:w-8 rounded-full" />
                    <div className="space-y-1.5 sm:space-y-2 flex-1">
                      <Skeleton className="h-3.5 sm:h-4 w-3/4" />
                      <Skeleton className="h-2.5 sm:h-3 w-1/2" />
                    </div>
                  </div>
                ))}
              </div>
            ) : notificationsData?.notifications &&
              notificationsData.notifications.length > 0 ? (
              <div className="space-y-2 sm:space-y-4">
                {notificationsData.notifications
                  .slice(0, 5)
                  .map((notification) => (
                    <div
                      key={notification.id}
                      className="flex items-start gap-3 sm:gap-4 rounded-lg p-2 transition-colors hover:bg-muted"
                    >
                      <div
                        className={`flex h-7 w-7 sm:h-8 sm:w-8 items-center justify-center rounded-full shrink-0 ${
                          notification.type === "success"
                            ? "bg-green-100 dark:bg-green-900/20"
                            : notification.type === "warning"
                              ? "bg-yellow-100 dark:bg-yellow-900/20"
                              : notification.type === "error"
                                ? "bg-red-100 dark:bg-red-900/20"
                                : "bg-blue-100 dark:bg-blue-900/20"
                        }`}
                      >
                        <Bell
                          className={`h-3.5 w-3.5 sm:h-4 sm:w-4 ${
                            notification.type === "success"
                              ? "text-green-600 dark:text-green-400"
                              : notification.type === "warning"
                                ? "text-yellow-600 dark:text-yellow-400"
                                : notification.type === "error"
                                  ? "text-red-600 dark:text-red-400"
                                  : "text-blue-600 dark:text-blue-400"
                          }`}
                        />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <p className="font-medium text-sm sm:text-base line-clamp-1">
                            {notification.title}
                          </p>
                          {!notification.is_read && (
                            <div className="h-2 w-2 rounded-full bg-primary shrink-0" />
                          )}
                        </div>
                        <p className="text-xs sm:text-sm text-muted-foreground line-clamp-1">
                          {notification.message}
                        </p>
                        <p className="text-xs text-muted-foreground mt-0.5 sm:mt-1">
                          {formatRelativeTime(notification.created_at)}
                        </p>
                      </div>
                    </div>
                  ))}
              </div>
            ) : (
              <p className="text-center text-muted-foreground py-6 sm:py-8 text-sm">
                No notifications
              </p>
            )}
          </CardContent>
        </Card>
      </motion.div>

      {/* Quick Actions: Chat & Telegram Bot */}
      <motion.div
        variants={itemVariants}
        className="grid gap-4 sm:gap-6 lg:grid-cols-2"
      >
        {/* Chat Section */}
        <Card>
          <CardHeader className="p-4 sm:p-6">
            <CardTitle className="flex items-center gap-2 text-base sm:text-lg">
              <MessageCircle className="h-5 w-5 text-primary" />
              Course Chat
            </CardTitle>
            <CardDescription className="text-xs sm:text-sm">
              Connect with teachers and classmates
            </CardDescription>
          </CardHeader>
          <CardContent className="p-4 pt-0 sm:p-6 sm:pt-0">
            <p className="text-xs sm:text-sm text-muted-foreground mb-4">
              Have questions about assignments or course materials? Start a
              conversation with your teachers or fellow students.
            </p>
            <Button asChild className="w-full h-11 sm:h-10 touch-manipulation">
              <Link href="/app/chat">
                <MessageCircle className="mr-2 h-4 w-4" />
                Open Chat
              </Link>
            </Button>
          </CardContent>
        </Card>

        {/* Telegram Bot */}
        <Card className="bg-linear-to-br from-blue-50 to-cyan-50 dark:from-blue-950/30 dark:to-cyan-950/30 border-blue-200 dark:border-blue-800">
          <CardHeader className="p-4 sm:p-6">
            <CardTitle className="flex items-center gap-2 text-base sm:text-lg">
              <Send className="h-5 w-5 text-blue-500" />
              Telegram Bot
            </CardTitle>
            <CardDescription className="text-xs sm:text-sm">
              Get instant notifications and updates
            </CardDescription>
          </CardHeader>
          <CardContent className="p-4 pt-0 sm:p-6 sm:pt-0">
            <p className="text-xs sm:text-sm text-muted-foreground mb-3 sm:mb-4">
              Stay on top of your studies with our Telegram bot! Get notified
              about:
            </p>
            <ul className="text-xs sm:text-sm text-muted-foreground space-y-1 mb-4">
              <li className="flex items-center gap-2">
                <Clock className="h-3 w-3" />
                Upcoming deadlines
              </li>
              <li className="flex items-center gap-2">
                <FileText className="h-3 w-3" />
                New assignments
              </li>
              <li className="flex items-center gap-2">
                <Calendar className="h-3 w-3" />
                Schedule reminders
              </li>
              <li className="flex items-center gap-2">
                <Bell className="h-3 w-3" />
                Points and grades
              </li>
            </ul>
            <Button
              asChild
              variant="default"
              className="w-full h-11 sm:h-10 bg-blue-500 hover:bg-blue-600 touch-manipulation"
            >
              <a
                href="https://t.me/a_lms_bot"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Send className="mr-2 h-4 w-4" />
                Open @a_lms_bot
                <ExternalLink className="ml-2 h-3 w-3" />
              </a>
            </Button>
          </CardContent>
        </Card>
      </motion.div>
    </motion.div>
  );
}
