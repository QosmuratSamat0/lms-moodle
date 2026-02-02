"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import {
  Calendar,
  Clock,
  AlertTriangle,
  CheckCircle2,
  FileText,
  BookOpen,
  Filter,
} from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
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
  formatDate,
  formatDateTime,
  isOverdue,
  getDaysUntil,
  cn,
} from "@/lib/helpers";
import assignmentService from "@/services/assignments";
import type { Assignment } from "@/types";

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

function DeadlineSkeleton() {
  return (
    <div className="flex items-center gap-4 p-4 rounded-lg border">
      <Skeleton className="h-12 w-12 rounded-lg" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-48" />
        <Skeleton className="h-3 w-32" />
      </div>
      <Skeleton className="h-6 w-20" />
    </div>
  );
}

function DeadlineCard({ assignment }: { assignment: Assignment }) {
  const overdue = assignment.due_at ? isOverdue(assignment.due_at) : false;
  const daysUntil = assignment.due_at ? getDaysUntil(assignment.due_at) : null;

  const getUrgencyColor = () => {
    if (overdue) return "border-destructive bg-destructive/5";
    if (daysUntil !== null && daysUntil <= 1)
      return "border-orange-500 bg-orange-500/5";
    if (daysUntil !== null && daysUntil <= 3)
      return "border-yellow-500 bg-yellow-500/5";
    return "";
  };

  const getUrgencyBadge = () => {
    if (overdue) return { label: "Overdue", variant: "destructive" as const };
    if (daysUntil !== null && daysUntil <= 1)
      return { label: "Due Tomorrow", variant: "default" as const };
    if (daysUntil !== null && daysUntil <= 3)
      return { label: `${daysUntil} days left`, variant: "secondary" as const };
    return { label: `${daysUntil} days left`, variant: "outline" as const };
  };

  const urgencyBadge = assignment.due_at ? getUrgencyBadge() : null;

  return (
    <motion.div variants={itemVariants}>
      <Link
        href={`/app/courses/${assignment.course_id}/assignments/${assignment.id}`}
      >
        <Card
          className={cn("transition-all hover:shadow-lg", getUrgencyColor())}
        >
          <CardContent className="flex items-center gap-4 p-4">
            <div
              className={cn(
                "flex h-12 w-12 items-center justify-center rounded-lg",
                overdue
                  ? "bg-destructive/10 text-destructive"
                  : "bg-primary/10 text-primary",
              )}
            >
              <FileText className="h-6 w-6" />
            </div>
            <div className="flex-1 min-w-0">
              <h4 className="font-medium truncate">{assignment.title}</h4>
              <div className="flex items-center gap-2 mt-1">
                <BookOpen className="h-3 w-3 text-muted-foreground" />
                <span className="text-sm text-muted-foreground truncate">
                  {assignment.course_title}
                </span>
              </div>
            </div>
            <div className="flex flex-col items-end gap-2">
              {urgencyBadge && (
                <Badge variant={urgencyBadge.variant}>
                  {urgencyBadge.label}
                </Badge>
              )}
              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                <Clock className="h-3 w-3" />
                <span>
                  {assignment.due_at
                    ? formatDateTime(assignment.due_at)
                    : "No deadline"}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </Link>
    </motion.div>
  );
}

export default function DeadlinesPage() {
  const [filter, setFilter] = useState<"all" | "upcoming" | "overdue">("all");
  const [courseFilter, setCourseFilter] = useState<string>("all");

  const { data, isLoading, error } = useQuery({
    queryKey: ["assignments"],
    queryFn: () => assignmentService.list(),
  });

  const assignments = data?.assignments || [];

  // Get unique courses for filter
  const courses = Array.from(
    new Set(
      assignments.map((a) =>
        JSON.stringify({ id: a.course_id, title: a.course_title }),
      ),
    ),
  ).map((s) => JSON.parse(s));

  // Filter assignments with due dates
  let filteredAssignments = assignments.filter((a) => a.due_at);

  // Apply course filter
  if (courseFilter !== "all") {
    filteredAssignments = filteredAssignments.filter(
      (a) => a.course_id === courseFilter,
    );
  }

  // Apply status filter
  if (filter === "upcoming") {
    filteredAssignments = filteredAssignments.filter(
      (a) => !isOverdue(a.due_at!),
    );
  } else if (filter === "overdue") {
    filteredAssignments = filteredAssignments.filter((a) =>
      isOverdue(a.due_at!),
    );
  }

  // Sort by due date
  filteredAssignments.sort((a, b) => {
    const dateA = new Date(a.due_at!).getTime();
    const dateB = new Date(b.due_at!).getTime();
    return dateA - dateB;
  });

  // Group by date
  const groupedByDate = filteredAssignments.reduce(
    (acc, assignment) => {
      const date = formatDate(assignment.due_at!);
      if (!acc[date]) acc[date] = [];
      acc[date].push(assignment);
      return acc;
    },
    {} as Record<string, Assignment[]>,
  );

  const overdueCount = assignments.filter(
    (a) => a.due_at && isOverdue(a.due_at),
  ).length;
  const upcomingCount = assignments.filter(
    (a) => a.due_at && !isOverdue(a.due_at),
  ).length;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Deadlines"
        description="Track all your upcoming assignment deadlines"
      />

      {/* Summary Cards */}
      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Total</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {assignments.filter((a) => a.due_at).length}
            </div>
            <p className="text-xs text-muted-foreground">Active deadlines</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Upcoming</CardTitle>
            <CheckCircle2 className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {upcomingCount}
            </div>
            <p className="text-xs text-muted-foreground">On track</p>
          </CardContent>
        </Card>
        <Card className={overdueCount > 0 ? "border-destructive" : ""}>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Overdue</CardTitle>
            <AlertTriangle className="h-4 w-4 text-destructive" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-destructive">
              {overdueCount}
            </div>
            <p className="text-xs text-muted-foreground">Need attention</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          value={filter}
          onValueChange={(v) => setFilter(v as "all" | "upcoming" | "overdue")}
        >
          <TabsList>
            <TabsTrigger value="all">All</TabsTrigger>
            <TabsTrigger value="upcoming">Upcoming</TabsTrigger>
            <TabsTrigger value="overdue">
              Overdue
              {overdueCount > 0 && (
                <Badge variant="destructive" className="ml-2">
                  {overdueCount}
                </Badge>
              )}
            </TabsTrigger>
          </TabsList>
        </Tabs>

        <Select value={courseFilter} onValueChange={setCourseFilter}>
          <SelectTrigger className="w-[200px]">
            <Filter className="h-4 w-4 mr-2" />
            <SelectValue placeholder="Filter by course" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Courses</SelectItem>
            {courses.map((course) => (
              <SelectItem key={course.id} value={course.id}>
                {course.title}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Deadlines List */}
      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3, 4, 5].map((i) => (
            <DeadlineSkeleton key={i} />
          ))}
        </div>
      ) : error ? (
        <EmptyState
          title="Error loading deadlines"
          description="There was a problem loading your deadlines."
          action={
            <Button onClick={() => window.location.reload()}>Retry</Button>
          }
        />
      ) : filteredAssignments.length === 0 ? (
        <EmptyState
          icon={<Calendar className="h-8 w-8 text-muted-foreground" />}
          title="No deadlines"
          description={
            filter === "overdue"
              ? "Great! You don't have any overdue assignments."
              : "You don't have any upcoming deadlines."
          }
        />
      ) : (
        <div className="space-y-6">
          {Object.entries(groupedByDate).map(([date, dateAssignments]) => (
            <div key={date}>
              <h3 className="font-medium text-sm text-muted-foreground mb-3 sticky top-0 bg-background py-2">
                {date}
              </h3>
              <motion.div
                variants={containerVariants}
                initial="hidden"
                animate="visible"
                className="space-y-3"
              >
                {dateAssignments.map((assignment) => (
                  <DeadlineCard key={assignment.id} assignment={assignment} />
                ))}
              </motion.div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
