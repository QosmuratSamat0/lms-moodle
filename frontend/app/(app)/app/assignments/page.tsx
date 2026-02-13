"use client";

import { useState, useMemo } from "react";
import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import {
  FileText,
  Clock,
  CheckCircle2,
  AlertTriangle,
  Calendar,
  XCircle,
  Award,
} from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { SearchInput } from "@/components/common/search-input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  formatDateTime,
  isOverdue,
  getDaysUntil,
  getTimeRemaining,
  getSubmissionTimeInfo,
  cn,
} from "@/lib/helpers";
import { useAuthStore } from "@/store/auth-store";
import { useTeacherProfile } from "@/hooks/use-profile";
import assignmentService from "@/services/assignments";
import courseService from "@/services/courses";
import submissionService from "@/services/submissions";
import type { Assignment } from "@/types";
import type { Submission } from "@/types/submission";

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.05 },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
};

type FilterStatus = "all" | "pending" | "submitted" | "graded" | "overdue";

interface AssignmentWithSubmission extends Assignment {
  submission?: Submission;
}

function AssignmentSkeleton() {
  return (
    <div className="flex items-center gap-4 p-4 rounded-lg border">
      <Skeleton className="h-12 w-12 rounded-lg" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-5 w-3/4" />
        <Skeleton className="h-4 w-1/2" />
      </div>
      <Skeleton className="h-6 w-20" />
    </div>
  );
}

function getSubmissionStatus(assignment: AssignmentWithSubmission) {
  if (assignment.submission?.grade_score !== undefined) {
    return "graded";
  }
  if (assignment.submission) {
    return "submitted";
  }
  if (assignment.due_at && isOverdue(assignment.due_at)) {
    return "overdue";
  }
  return "pending";
}

function StatusBadge({ assignment }: { assignment: AssignmentWithSubmission }) {
  const status = getSubmissionStatus(assignment);
  const submission = assignment.submission;

  switch (status) {
    case "graded":
      return (
        <div className="flex flex-col items-end gap-1">
          <Badge className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 gap-1">
            <Award className="h-3 w-3" />
            Graded
          </Badge>
          <span className="text-sm font-semibold text-green-600 dark:text-green-400">
            {submission?.grade_score?.toFixed(0)}%
          </span>
        </div>
      );
    case "submitted":
      return (
        <Badge className="bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400 gap-1">
          <CheckCircle2 className="h-3 w-3" />
          Submitted
        </Badge>
      );
    case "overdue":
      return (
        <Badge variant="destructive" className="gap-1">
          <XCircle className="h-3 w-3" />
          Not submitted
        </Badge>
      );
    default:
      if (!assignment.due_at) {
        return (
          <Badge variant="secondary" className="gap-1">
            <Calendar className="h-3 w-3" />
            No deadline
          </Badge>
        );
      }
      const daysLeft = getDaysUntil(assignment.due_at);
      if (daysLeft <= 1) {
        return (
          <Badge variant="destructive" className="gap-1">
            <AlertTriangle className="h-3 w-3" />
            Due soon
          </Badge>
        );
      }
      if (daysLeft <= 3) {
        return (
          <Badge className="bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400 gap-1">
            <Clock className="h-3 w-3" />
            {daysLeft} days left
          </Badge>
        );
      }
      return (
        <Badge variant="outline" className="gap-1">
          <Clock className="h-3 w-3" />
          Not submitted
        </Badge>
      );
  }
}

export default function AssignmentsPage() {
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<FilterStatus>("all");
  const [courseFilter, setCourseFilter] = useState<string>("all");
  const { isStudent, isTeacher, user } = useAuthStore();
  const { data: teacherProfile } = useTeacherProfile();

  // Fetch enrolled courses for students
  const { data: enrollmentsData } = useQuery({
    queryKey: ["my-enrollments", user?.id],
    queryFn: () => courseService.getMyEnrollments(user!.id),
    enabled: isStudent() && !!user?.id,
  });

  // Fetch all courses (used for teachers and admins)
  const { data: coursesData } = useQuery({
    queryKey: ["courses", "list"],
    queryFn: () => courseService.list({ limit: 100 }),
    enabled: !isStudent(),
  });

  // Get course IDs based on role
  const courseIds = useMemo(() => {
    if (isStudent() && enrollmentsData) {
      const enrollments = Array.isArray(enrollmentsData)
        ? enrollmentsData
        : (enrollmentsData as any)?.enrollments || [];
      return enrollments
        .filter((e: any) => e.status === "active")
        .map((e: any) => e.course_id);
    }
    // For teachers - filter courses by ownership
    if ((isTeacher() || !isStudent()) && coursesData) {
      const allCourses = Array.isArray(coursesData)
        ? coursesData
        : (coursesData as any)?.courses || [];
      if (isTeacher() && teacherProfile?.id) {
        return allCourses
          .filter((c: any) => c.owner_teacher_id === teacherProfile.id)
          .map((c: any) => c.id);
      }
      return allCourses.map((c: any) => c.id);
    }
    return [];
  }, [
    isStudent,
    isTeacher,
    enrollmentsData,
    coursesData,
    user,
    teacherProfile,
  ]);

  // Fetch assignments for each enrolled course
  const { data: assignmentsData, isLoading: assignmentsLoading } = useQuery({
    queryKey: ["my-course-assignments", courseIds],
    queryFn: async () => {
      if (courseIds.length === 0) return { assignments: [] };

      const results = await Promise.all(
        courseIds.map((courseId) =>
          assignmentService.getByCourse(courseId, { limit: 100 }),
        ),
      );

      const allAssignments = results.flatMap((r) => r.assignments || []);
      return {
        assignments: allAssignments,
        total: allAssignments.length,
      };
    },
    enabled: courseIds.length > 0,
  });

  // Fetch my submissions
  const { data: mySubmissionsData } = useQuery({
    queryKey: ["my-submissions", user?.id],
    queryFn: () => submissionService.getMySubmissions(user!.id, { limit: 200 }),
    enabled: isStudent() && !!user?.id,
  });

  // Create a map of submissions by assignment ID
  const submissionMap = useMemo(() => {
    const map = new Map<string, Submission>();
    mySubmissionsData?.submissions?.forEach((s) => {
      map.set(s.assignment_id, s);
    });
    return map;
  }, [mySubmissionsData]);

  // Merge assignments with submissions
  const assignmentsWithSubmissions: AssignmentWithSubmission[] = useMemo(() => {
    return (
      assignmentsData?.assignments?.map((assignment) => ({
        ...assignment,
        submission: submissionMap.get(assignment.id),
      })) || []
    );
  }, [assignmentsData, submissionMap]);

  // Get unique courses for filter
  const courses = useMemo(() => {
    if (isStudent() && enrollmentsData) {
      const enrollments = Array.isArray(enrollmentsData)
        ? enrollmentsData
        : (enrollmentsData as any)?.enrollments || [];
      return enrollments
        .filter((e: any) => e.status === "active")
        .map((e: any) => ({
          id: e.course_id,
          title: e.course_title || "Unknown Course",
        }));
    }
    // For teachers - filter their own courses
    const allCourses = Array.isArray(coursesData)
      ? coursesData
      : (coursesData as any)?.courses || [];
    if (isTeacher() && teacherProfile?.id) {
      return allCourses
        .filter((c: any) => c.owner_teacher_id === teacherProfile.id)
        .map((c: any) => ({ id: c.id, title: c.title }));
    }
    return allCourses.map((c: any) => ({ id: c.id, title: c.title }));
  }, [isStudent, isTeacher, enrollmentsData, coursesData, teacherProfile]);

  // Filter assignments
  const filteredAssignments = assignmentsWithSubmissions.filter(
    (assignment) => {
      // Search filter
      if (
        search &&
        !assignment.title.toLowerCase().includes(search.toLowerCase()) &&
        !assignment.course_title?.toLowerCase().includes(search.toLowerCase())
      ) {
        return false;
      }

      // Course filter
      if (courseFilter !== "all" && assignment.course_id !== courseFilter) {
        return false;
      }

      // Status filter
      const status = getSubmissionStatus(assignment);
      if (statusFilter === "pending" && status !== "pending") return false;
      if (statusFilter === "submitted" && status !== "submitted") return false;
      if (statusFilter === "graded" && status !== "graded") return false;
      if (statusFilter === "overdue" && status !== "overdue") return false;

      return true;
    },
  );

  // Sort by due date (closest first, overdue at top, no due date at end)
  const sortedAssignments = [...filteredAssignments].sort((a, b) => {
    // Overdue and not submitted first
    const aOverdue = a.due_at && isOverdue(a.due_at) && !a.submission;
    const bOverdue = b.due_at && isOverdue(b.due_at) && !b.submission;
    if (aOverdue && !bOverdue) return -1;
    if (!aOverdue && bOverdue) return 1;

    // Then by due date
    if (!a.due_at && !b.due_at) return 0;
    if (!a.due_at) return 1;
    if (!b.due_at) return -1;
    return new Date(a.due_at).getTime() - new Date(b.due_at).getTime();
  });

  // Stats
  const stats = useMemo(() => {
    const total = assignmentsWithSubmissions.length;
    const submitted = assignmentsWithSubmissions.filter(
      (a) => a.submission,
    ).length;
    const graded = assignmentsWithSubmissions.filter(
      (a) => a.submission?.grade_score !== undefined,
    ).length;
    const overdue = assignmentsWithSubmissions.filter(
      (a) => a.due_at && isOverdue(a.due_at) && !a.submission,
    ).length;
    const pending = total - submitted;

    return { total, submitted, graded, overdue, pending };
  }, [assignmentsWithSubmissions]);

  const isLoading = assignmentsLoading;

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="space-y-6"
    >
      <PageHeader
        title={isTeacher() ? "Course Assignments" : "My Assignments"}
        description={
          isTeacher()
            ? "View and manage assignments in your courses"
            : "View and track all your assignments across enrolled courses"
        }
      />

      {/* Stats Cards */}
      <motion.div variants={itemVariants} className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Pending
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">
              {stats.pending}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Submitted
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {stats.submitted}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Graded
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {stats.graded}
            </div>
          </CardContent>
        </Card>
        <Card
          className={
            stats.overdue > 0 ? "border-red-200 dark:border-red-900" : ""
          }
        >
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Overdue
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {stats.overdue}
            </div>
          </CardContent>
        </Card>
      </motion.div>

      {/* Progress */}
      {stats.total > 0 && (
        <motion.div variants={itemVariants}>
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium">Completion Progress</span>
                <span className="text-sm text-muted-foreground">
                  {stats.submitted}/{stats.total} submitted (
                  {Math.round((stats.submitted / stats.total) * 100)}%)
                </span>
              </div>
              <div className="h-2 w-full bg-muted rounded-full overflow-hidden">
                <div
                  className="h-full bg-primary transition-all"
                  style={{ width: `${(stats.submitted / stats.total) * 100}%` }}
                />
              </div>
            </CardContent>
          </Card>
        </motion.div>
      )}

      {/* Filters */}
      <motion.div
        variants={itemVariants}
        className="flex flex-col sm:flex-row gap-4"
      >
        <SearchInput
          value={search}
          onChange={setSearch}
          placeholder="Search assignments..."
          className="sm:max-w-sm"
        />

        <div className="flex gap-2 flex-wrap">
          <Select value={courseFilter} onValueChange={setCourseFilter}>
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder="All courses" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All courses</SelectItem>
              {courses.map((course) => (
                <SelectItem key={course.id} value={course.id}>
                  {course.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Tabs
            value={statusFilter}
            onValueChange={(v) => setStatusFilter(v as FilterStatus)}
          >
            <TabsList>
              <TabsTrigger value="all">All</TabsTrigger>
              <TabsTrigger value="pending">Pending</TabsTrigger>
              <TabsTrigger value="submitted">Submitted</TabsTrigger>
              <TabsTrigger value="graded">Graded</TabsTrigger>
              <TabsTrigger value="overdue">Overdue</TabsTrigger>
            </TabsList>
          </Tabs>
        </div>
      </motion.div>

      {/* Assignments List */}
      <motion.div variants={itemVariants}>
        {isLoading ? (
          <div className="space-y-3">
            {[1, 2, 3, 4, 5].map((i) => (
              <AssignmentSkeleton key={i} />
            ))}
          </div>
        ) : sortedAssignments.length === 0 ? (
          <EmptyState
            icon={<FileText className="h-8 w-8 text-muted-foreground" />}
            title="No assignments found"
            description={
              search || statusFilter !== "all" || courseFilter !== "all"
                ? "Try adjusting your filters"
                : "You don't have any assignments in your enrolled courses yet"
            }
          />
        ) : (
          <div className="space-y-3">
            {sortedAssignments.map((assignment) => {
              const status = getSubmissionStatus(assignment);
              const isOverdueNotSubmitted = status === "overdue";

              return (
                <motion.div key={assignment.id} variants={itemVariants}>
                  <Link
                    href={`/app/courses/${assignment.course_id}/assignments/${assignment.id}`}
                    className="block"
                  >
                    <Card
                      className={cn(
                        "hover:bg-muted/50 transition-colors",
                        isOverdueNotSubmitted &&
                          "border-red-200 dark:border-red-900 bg-red-50/50 dark:bg-red-900/10",
                      )}
                    >
                      <CardContent className="flex items-start gap-4 p-4">
                        <div
                          className={cn(
                            "flex h-12 w-12 items-center justify-center rounded-lg shrink-0",
                            status === "graded"
                              ? "bg-green-100 dark:bg-green-900/30"
                              : status === "submitted"
                                ? "bg-blue-100 dark:bg-blue-900/30"
                                : isOverdueNotSubmitted
                                  ? "bg-red-100 dark:bg-red-900/30"
                                  : "bg-primary/10",
                          )}
                        >
                          {status === "graded" ? (
                            <Award className="h-6 w-6 text-green-600 dark:text-green-400" />
                          ) : status === "submitted" ? (
                            <CheckCircle2 className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                          ) : isOverdueNotSubmitted ? (
                            <XCircle className="h-6 w-6 text-red-600 dark:text-red-400" />
                          ) : (
                            <FileText className="h-6 w-6 text-primary" />
                          )}
                        </div>
                        <div className="flex-1 min-w-0 space-y-1">
                          <h3 className="font-medium">{assignment.title}</h3>
                          <p className="text-sm text-muted-foreground truncate">
                            {assignment.course_title}
                          </p>
                          {/* Time info */}
                          <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm">
                            {assignment.submission ? (
                              <span
                                className={cn(
                                  "flex items-center gap-1",
                                  assignment.due_at &&
                                    new Date(
                                      assignment.submission.submitted_at,
                                    ) > new Date(assignment.due_at)
                                    ? "text-red-600 dark:text-red-400"
                                    : "text-green-600 dark:text-green-400",
                                )}
                              >
                                <CheckCircle2 className="h-3.5 w-3.5" />
                                {getSubmissionTimeInfo(
                                  assignment.submission.submitted_at,
                                  assignment.due_at,
                                )}
                              </span>
                            ) : assignment.due_at ? (
                              <span
                                className={cn(
                                  "flex items-center gap-1",
                                  isOverdue(assignment.due_at)
                                    ? "text-red-600 dark:text-red-400"
                                    : getDaysUntil(assignment.due_at) <= 3
                                      ? "text-orange-600 dark:text-orange-400"
                                      : "text-muted-foreground",
                                )}
                              >
                                <Clock className="h-3.5 w-3.5" />
                                {getTimeRemaining(assignment.due_at)}
                              </span>
                            ) : (
                              <span className="text-muted-foreground flex items-center gap-1">
                                <Calendar className="h-3.5 w-3.5" />
                                No deadline
                              </span>
                            )}
                            {assignment.due_at && (
                              <span className="text-muted-foreground">
                                Due: {formatDateTime(assignment.due_at)}
                              </span>
                            )}
                          </div>
                          {/* Grade feedback */}
                          {assignment.submission?.grade_feedback && (
                            <p className="text-sm text-muted-foreground bg-muted/50 px-2 py-1 rounded mt-1">
                              Feedback: {assignment.submission.grade_feedback}
                            </p>
                          )}
                        </div>
                        <div className="flex flex-col items-end gap-2 shrink-0">
                          <StatusBadge assignment={assignment} />
                          <span className="text-xs text-muted-foreground">
                            {assignment.max_points}%
                          </span>
                        </div>
                      </CardContent>
                    </Card>
                  </Link>
                </motion.div>
              );
            })}
          </div>
        )}
      </motion.div>
    </motion.div>
  );
}
