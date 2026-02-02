"use client";

import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import Link from "next/link";
import {
  BookOpen,
  FileText,
  GraduationCap,
  Calendar,
  Users,
  ArrowLeft,
  MapPin,
  Video,
  Clock,
} from "lucide-react";
import { format } from "date-fns";

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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import {
  CreateAssignmentDialog,
  EditAssignmentDialog,
} from "@/components/assignments";
import { CreateLectureDialog, EditLectureDialog } from "@/components/lectures";
import { useAuthStore } from "@/store/auth-store";
import { getFullName, formatDate } from "@/lib/helpers";
import courseService from "@/services/courses";
import assignmentService from "@/services/assignments";
import { scheduleService } from "@/services/schedule";

const eventTypeLabels: Record<string, string> = {
  class: "Lecture",
  lab: "Lab",
  exam: "Exam",
  office_hours: "Office Hours",
  event: "Event",
};

export default function CourseDetailPage() {
  const params = useParams();
  const courseId = params.courseId as string;
  const { user } = useAuthStore();
  const isTeacher = user?.role === "teacher" || user?.role === "admin";

  const { data: course, isLoading: courseLoading } = useQuery({
    queryKey: ["course", courseId],
    queryFn: () => courseService.get(courseId),
  });

  const { data: assignmentsData, isLoading: assignmentsLoading } = useQuery({
    queryKey: ["course-assignments", courseId],
    queryFn: () => assignmentService.getByCourse(courseId),
  });

  const { data: scheduleData, isLoading: scheduleLoading } = useQuery({
    queryKey: ["course-schedule", courseId],
    queryFn: () => scheduleService.getCourseSchedule(courseId),
  });

  if (courseLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Skeleton className="h-10 w-10" />
          <div className="space-y-2">
            <Skeleton className="h-8 w-64" />
            <Skeleton className="h-4 w-32" />
          </div>
        </div>
        <Skeleton className="h-100 w-full" />
      </div>
    );
  }

  if (!course) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <h2 className="text-xl font-semibold">Course not found</h2>
        <Button asChild className="mt-4">
          <Link href="/app/courses">Back to courses</Link>
        </Button>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/app/courses">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <PageHeader
          title={course.title}
          description={course.description}
          action={
            <Badge variant={course.is_active ? "default" : "secondary"}>
              {course.is_active ? "Active" : "Inactive"}
            </Badge>
          }
        />
      </div>

      {/* Course Info Card */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-wrap gap-6">
            <div className="flex items-center gap-2">
              <Users className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm">
                Instructor:{" "}
                {getFullName(
                  course.teacher_first_name,
                  course.teacher_last_name,
                )}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm">
                Created: {formatDate(course.created_at)}
              </span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tabs */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList className="flex flex-wrap h-auto gap-2">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="lectures">Lectures</TabsTrigger>
          <TabsTrigger value="assignments">Assignments</TabsTrigger>
          <TabsTrigger value="quizzes">Quizzes</TabsTrigger>
          <TabsTrigger value="grades">Grades</TabsTrigger>
          <TabsTrigger value="attendance">Attendance</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-3">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Lectures</CardTitle>
                <BookOpen className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {scheduleLoading ? (
                    <Skeleton className="h-8 w-8" />
                  ) : (
                    scheduleData?.events?.length || 0
                  )}
                </div>
                <p className="text-xs text-muted-foreground">Total sessions</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">
                  Assignments
                </CardTitle>
                <FileText className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {assignmentsLoading ? (
                    <Skeleton className="h-8 w-8" />
                  ) : (
                    assignmentsData?.total || 0
                  )}
                </div>
                <p className="text-xs text-muted-foreground">
                  Total assignments
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Quizzes</CardTitle>
                <GraduationCap className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">4</div>
                <p className="text-xs text-muted-foreground">Total quizzes</p>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>About this course</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground">
                {course.description ||
                  "No description available for this course."}
              </p>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="lectures" className="space-y-4">
          <div className="flex justify-between items-center">
            <div>
              <h3 className="text-lg font-semibold">Lectures & Sessions</h3>
              <p className="text-sm text-muted-foreground">
                Course schedule, lectures, and lab sessions
              </p>
            </div>
            {isTeacher && <CreateLectureDialog courseId={courseId} />}
          </div>

          {scheduleLoading ? (
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-24 w-full" />
              ))}
            </div>
          ) : scheduleData?.events && scheduleData.events.length > 0 ? (
            <div className="space-y-4">
              {scheduleData.events.map((event) => (
                <Card
                  key={event.id}
                  className="hover:border-primary/50 transition-colors"
                >
                  <CardContent className="p-4">
                    <div className="flex items-start justify-between">
                      <div className="flex items-start gap-4">
                        <div className="flex h-10 w-10 items-center justify-center rounded bg-primary/10">
                          {event.is_online ? (
                            <Video className="h-5 w-5 text-primary" />
                          ) : (
                            <BookOpen className="h-5 w-5 text-primary" />
                          )}
                        </div>
                        <div className="space-y-1">
                          <p className="font-medium">{event.title}</p>
                          {event.description && (
                            <p className="text-sm text-muted-foreground line-clamp-2">
                              {event.description}
                            </p>
                          )}
                          <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
                            <span className="flex items-center gap-1">
                              <Calendar className="h-4 w-4" />
                              {format(new Date(event.start_time), "PPP")}
                            </span>
                            <span className="flex items-center gap-1">
                              <Clock className="h-4 w-4" />
                              {format(
                                new Date(event.start_time),
                                "HH:mm",
                              )} - {format(new Date(event.end_time), "HH:mm")}
                            </span>
                            {event.location && (
                              <span className="flex items-center gap-1">
                                <MapPin className="h-4 w-4" />
                                {event.location}
                              </span>
                            )}
                            {event.is_online && event.meeting_url && (
                              <a
                                href={event.meeting_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex items-center gap-1 text-primary hover:underline"
                              >
                                <Video className="h-4 w-4" />
                                Join online
                              </a>
                            )}
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant="outline">
                          {eventTypeLabels[event.event_type] ||
                            event.event_type}
                        </Badge>
                        {event.recurrence !== "none" && (
                          <Badge variant="secondary">{event.recurrence}</Badge>
                        )}
                        {isTeacher && (
                          <EditLectureDialog
                            lecture={event}
                            courseId={courseId}
                          />
                        )}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="py-8 text-center text-muted-foreground">
                No lectures or sessions scheduled yet
                {isTeacher && (
                  <p className="mt-2 text-sm">
                    Click "Add Lecture" to create your first session.
                  </p>
                )}
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="assignments" className="space-y-4">
          <div className="flex justify-between items-center">
            <div>
              <h3 className="text-lg font-semibold">Assignments</h3>
              <p className="text-sm text-muted-foreground">
                View and submit assignments
              </p>
            </div>
            {isTeacher && <CreateAssignmentDialog courseId={courseId} />}
          </div>

          {assignmentsLoading ? (
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-24 w-full" />
              ))}
            </div>
          ) : assignmentsData?.assignments.length ? (
            <div className="space-y-4">
              {assignmentsData.assignments.map((assignment) => (
                <Card
                  key={assignment.id}
                  className="hover:border-primary/50 transition-colors"
                >
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <Link
                        href={`/app/courses/${courseId}/assignments/${assignment.id}`}
                        className="flex-1"
                      >
                        <div>
                          <CardTitle className="text-base hover:text-primary transition-colors">
                            {assignment.title}
                          </CardTitle>
                          <CardDescription className="line-clamp-2 mt-1">
                            {assignment.description}
                          </CardDescription>
                        </div>
                      </Link>
                      <div className="flex items-center gap-2">
                        <Badge>{assignment.max_points} pts</Badge>
                        {isTeacher && (
                          <EditAssignmentDialog
                            assignment={assignment}
                            courseId={courseId}
                          />
                        )}
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <span>
                        Due:{" "}
                        {assignment.due_at
                          ? formatDate(assignment.due_at)
                          : "No deadline"}
                      </span>
                      <span>•</span>
                      <span>
                        {assignment.submission_count || 0} submissions
                      </span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="py-8 text-center text-muted-foreground">
                No assignments yet
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="quizzes" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Quizzes</CardTitle>
              <CardDescription>Test your knowledge</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <div
                    key={i}
                    className="flex items-center gap-4 p-4 rounded-lg border hover:bg-muted/50 transition-colors cursor-pointer"
                  >
                    <div className="flex h-10 w-10 items-center justify-center rounded bg-green-100 dark:bg-green-900/20">
                      <GraduationCap className="h-5 w-5 text-green-600 dark:text-green-400" />
                    </div>
                    <div className="flex-1">
                      <p className="font-medium">
                        Quiz {i}: Chapter {i} Review
                      </p>
                      <p className="text-sm text-muted-foreground">
                        20 questions • 30 minutes
                      </p>
                    </div>
                    <Badge variant={i === 1 ? "default" : "outline"}>
                      {i === 1 ? "Completed" : "Not started"}
                    </Badge>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="grades" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Grades</CardTitle>
              <CardDescription>Your grades for this course</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex justify-between items-center p-4 rounded-lg bg-muted">
                  <span className="font-medium">Overall Grade</span>
                  <span className="text-2xl font-bold text-primary">
                    A (92%)
                  </span>
                </div>
                <div className="space-y-2">
                  {[
                    { name: "Assignment 1", score: 95, max: 100 },
                    { name: "Quiz 1", score: 18, max: 20 },
                    { name: "Assignment 2", score: 88, max: 100 },
                  ].map((item, i) => (
                    <div
                      key={i}
                      className="flex justify-between items-center p-3 rounded border"
                    >
                      <span>{item.name}</span>
                      <span className="font-medium">
                        {item.score}/{item.max}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="attendance" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Attendance</CardTitle>
              <CardDescription>Your attendance record</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex justify-between items-center p-4 rounded-lg bg-muted">
                  <span className="font-medium">Attendance Rate</span>
                  <span className="text-2xl font-bold text-green-600">95%</span>
                </div>
                <div className="space-y-2">
                  {[
                    { date: "Jan 20, 2026", status: "Present" },
                    { date: "Jan 22, 2026", status: "Present" },
                    { date: "Jan 24, 2026", status: "Absent" },
                    { date: "Jan 27, 2026", status: "Present" },
                  ].map((item, i) => (
                    <div
                      key={i}
                      className="flex justify-between items-center p-3 rounded border"
                    >
                      <span>{item.date}</span>
                      <Badge
                        variant={
                          item.status === "Present" ? "default" : "destructive"
                        }
                      >
                        {item.status}
                      </Badge>
                    </div>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </motion.div>
  );
}
