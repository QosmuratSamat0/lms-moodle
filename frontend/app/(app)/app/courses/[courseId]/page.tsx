"use client";

import { useState, useEffect, useMemo } from "react";
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
  ChevronRight,
  ChevronDown,
  Plus,
  Link as LinkIcon,
  File,
  Clock,
  CheckCircle2,
  AlertCircle,
  Award,
  Video,
  ExternalLink,
  Settings,
  Percent,
  Download,
  UsersRound,
  Pencil,
  Trash2,
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
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { Progress } from "@/components/ui/progress";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CreateContentDialog, SyllabusConfigDialog } from "@/components/course";
import { EditAssignmentDialog } from "@/components/assignments";
import {
  StudentAttendanceView,
  TeacherAttendanceView,
} from "@/components/attendance/attendance-views";
import { useAuthStore } from "@/store/auth-store";
import { getFullName, formatDate } from "@/lib/helpers";
import courseService from "@/services/courses";
import assignmentService from "@/services/assignments";
import submissionService from "@/services/submissions";
import { scheduleService } from "@/services/schedule";
import uploadService from "@/services/uploads";
import groupService from "@/services/groups";

// Total weeks for the course
const TOTAL_WEEKS = 10;

// Content type detection from description markers
type ContentType = "assignment" | "quiz" | "lecture" | "resource" | "link";

const getContentType = (description?: string): ContentType => {
  if (!description) return "assignment";
  if (description.includes("[QUIZ]")) return "quiz";
  if (description.includes("[LECTURE]")) return "lecture";
  if (description.includes("[RESOURCE]")) return "resource";
  if (description.includes("[LINK]")) return "link";
  return "assignment";
};

const getWeekNumber = (title: string): number => {
  const match = title.match(/Week\s+(\d+)/i);
  return match ? parseInt(match[1], 10) : 0;
};

// Extract file name from description if present
const getFileName = (description?: string): string | null => {
  if (!description) return null;
  const match = description.match(/\[FILE:(.+?)\]/);
  return match ? match[1] : null;
};

// Content type icons and colors
const contentTypeConfig: Record<
  ContentType,
  { icon: React.ElementType; color: string; bgColor: string }
> = {
  assignment: {
    icon: FileText,
    color: "text-blue-600 dark:text-blue-400",
    bgColor: "bg-blue-100 dark:bg-blue-900/30",
  },
  quiz: {
    icon: GraduationCap,
    color: "text-purple-600 dark:text-purple-400",
    bgColor: "bg-purple-100 dark:bg-purple-900/30",
  },
  lecture: {
    icon: Video,
    color: "text-green-600 dark:text-green-400",
    bgColor: "bg-green-100 dark:bg-green-900/30",
  },
  resource: {
    icon: File,
    color: "text-orange-600 dark:text-orange-400",
    bgColor: "bg-orange-100 dark:bg-orange-900/30",
  },
  link: {
    icon: LinkIcon,
    color: "text-cyan-600 dark:text-cyan-400",
    bgColor: "bg-cyan-100 dark:bg-cyan-900/30",
  },
};

// Syllabus grading configuration type
interface SyllabusConfig {
  registerMidterm: {
    weight: number; // 30%
    components: {
      name: string;
      count: number;
      weightEach: number; // percentage of register midterm
    }[];
  };
  registerEndterm: {
    weight: number; // 30%
    components: {
      name: string;
      count: number;
      weightEach: number;
    }[];
  };
  final: {
    weight: number; // 40%
  };
}

// Default syllabus config (Advanced Programming style)
const DEFAULT_SYLLABUS: SyllabusConfig = {
  registerMidterm: {
    weight: 30,
    components: [
      { name: "Assignment", count: 2, weightEach: 50 }, // 2 assignments, 50% each of register midterm
    ],
  },
  registerEndterm: {
    weight: 30,
    components: [
      { name: "Assignment", count: 2, weightEach: 50 }, // 2 more assignments
    ],
  },
  final: {
    weight: 40,
  },
};

export default function CourseDetailPage() {
  const params = useParams();
  const courseId = params.courseId as string;
  const { user } = useAuthStore();
  const isTeacher = user?.role === "teacher" || user?.role === "admin";
  const isStudent = user?.role === "student";

  // Selected group for teachers (to view specific group's data)
  const [selectedGroupId, setSelectedGroupId] = useState<string | null>(null);

  // Syllabus config state (would be fetched from backend in production)
  const [syllabusConfig, setSyllabusConfig] =
    useState<SyllabusConfig>(DEFAULT_SYLLABUS);
  const [syllabusDialogOpen, setSyllabusDialogOpen] = useState(false);

  // Expanded weeks state - will be computed based on content
  const [expandedWeeks, setExpandedWeeks] = useState<Set<number>>(new Set());
  const [addContentWeek, setAddContentWeek] = useState<number | null>(null);

  // Fetch teacher's group assignments for this course
  const { data: teacherAssignments, isLoading: assignmentsLoading2 } = useQuery(
    {
      queryKey: ["teacher-group-assignments"],
      queryFn: () => groupService.getMyAssignments({ limit: 100 }),
      enabled: isTeacher,
    },
  );

  // Filter assignments for this course only and get unique groups
  const courseGroups = useMemo(() => {
    if (!teacherAssignments?.assignments) return [];
    return teacherAssignments.assignments
      .filter((a) => a.course_id === courseId)
      .sort((a, b) => a.group_code.localeCompare(b.group_code));
  }, [teacherAssignments, courseId]);

  // Set default selected group when data loads
  useEffect(() => {
    if (courseGroups.length > 0 && !selectedGroupId) {
      setSelectedGroupId(courseGroups[0].group_id);
    }
  }, [courseGroups, selectedGroupId]);

  const toggleWeek = (week: number) => {
    setExpandedWeeks((prev) => {
      const next = new Set(prev);
      if (next.has(week)) {
        next.delete(week);
      } else {
        next.add(week);
      }
      return next;
    });
  };

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

  // Fetch student's submissions to show grades

  const { data: mySubmissions } = useQuery({
    queryKey: ["my-submissions", user?.id],
    queryFn: () => submissionService.getMySubmissions(user!.id),
    enabled: isStudent && !!user?.id,
  });

  // Helper to get submission for an assignment
  const getSubmissionForAssignment = (assignmentId: string) => {
    return mySubmissions?.submissions?.find(
      (s) => s.assignment_id === assignmentId,
    );
  };

  // Group content by week
  const getWeekContent = (weekNumber: number) => {
    const weekAssignments =
      assignmentsData?.assignments?.filter((a) => {
        const week = getWeekNumber(a.title);
        return week === weekNumber;
      }) || [];

    const weekLectures =
      scheduleData?.events?.filter((e) => {
        const week = getWeekNumber(e.title);
        return week === weekNumber;
      }) || [];

    return { assignments: weekAssignments, lectures: weekLectures };
  };

  // Get general content (week 0 or no week specified)
  const getGeneralContent = () => {
    const generalAssignments =
      assignmentsData?.assignments?.filter((a) => {
        const week = getWeekNumber(a.title);
        return week === 0;
      }) || [];

    const generalLectures =
      scheduleData?.events?.filter((e) => {
        const week = getWeekNumber(e.title);
        return week === 0;
      }) || [];

    return { assignments: generalAssignments, lectures: generalLectures };
  };

  // Compute which weeks have content
  const weeksWithContent = useMemo(() => {
    const weeks = new Set<number>();
    for (let i = 1; i <= TOTAL_WEEKS; i++) {
      const content = getWeekContent(i);
      if (content.assignments.length > 0 || content.lectures.length > 0) {
        weeks.add(i);
      }
    }
    return weeks;
  }, [assignmentsData, scheduleData]);

  // Auto-expand weeks with content when data loads
  useEffect(() => {
    if (!assignmentsLoading && !scheduleLoading) {
      setExpandedWeeks(new Set(weeksWithContent));
    }
  }, [weeksWithContent, assignmentsLoading, scheduleLoading]);

  // Calculate grades based on syllabus
  const calculateGrades = useMemo(() => {
    if (!mySubmissions?.submissions || !assignmentsData?.assignments) {
      return null;
    }

    const gradedSubmissions = mySubmissions.submissions
      .filter((s) => s.status === "graded" && s.grade_score !== undefined)
      .map((s) => {
        const assignment = assignmentsData.assignments.find(
          (a) => a.id === s.assignment_id,
        );
        return {
          ...s,
          assignment,
          percentage: assignment
            ? (s.grade_score! / assignment.max_points) * 100
            : 0,
        };
      });

    // Sort by assignment order (assuming title has order info)
    gradedSubmissions.sort((a, b) => {
      const orderA = a.assignment?.title?.match(/(\d+)/)?.[1] || "0";
      const orderB = b.assignment?.title?.match(/(\d+)/)?.[1] || "0";
      return parseInt(orderA) - parseInt(orderB);
    });

    // Calculate register midterm (first half of assignments)
    const midtermAssignmentCount = syllabusConfig.registerMidterm.components
      .filter((c) => c.name === "Assignment")
      .reduce((sum, c) => sum + c.count, 0);

    const midtermAssignments = gradedSubmissions.slice(
      0,
      midtermAssignmentCount,
    );
    const midtermAssignmentAvg =
      midtermAssignments.length > 0
        ? midtermAssignments.reduce((sum, s) => sum + s.percentage, 0) /
          midtermAssignmentCount
        : 0;
    const registerMidtermScore =
      (midtermAssignmentAvg / 100) * syllabusConfig.registerMidterm.weight;

    // Calculate register endterm (second half of assignments)
    const endtermAssignmentCount = syllabusConfig.registerEndterm.components
      .filter((c) => c.name === "Assignment")
      .reduce((sum, c) => sum + c.count, 0);

    const endtermAssignments = gradedSubmissions.slice(
      midtermAssignmentCount,
      midtermAssignmentCount + endtermAssignmentCount,
    );
    const endtermAssignmentAvg =
      endtermAssignments.length > 0
        ? endtermAssignments.reduce((sum, s) => sum + s.percentage, 0) /
          endtermAssignmentCount
        : 0;
    const registerEndtermScore =
      (endtermAssignmentAvg / 100) * syllabusConfig.registerEndterm.weight;

    // Final exam (not implemented yet, placeholder)
    const finalScore = 0;

    const totalScore = registerMidtermScore + registerEndtermScore + finalScore;
    const maxPossibleNow =
      (midtermAssignments.length > 0
        ? syllabusConfig.registerMidterm.weight
        : 0) +
      (endtermAssignments.length > 0
        ? syllabusConfig.registerEndterm.weight
        : 0);

    return {
      gradedSubmissions,
      registerMidterm: {
        score: registerMidtermScore,
        maxWeight: syllabusConfig.registerMidterm.weight,
        assignments: midtermAssignments,
        assignmentAvg: midtermAssignmentAvg,
      },
      registerEndterm: {
        score: registerEndtermScore,
        maxWeight: syllabusConfig.registerEndterm.weight,
        assignments: endtermAssignments,
        assignmentAvg: endtermAssignmentAvg,
      },
      final: {
        score: finalScore,
        maxWeight: syllabusConfig.final.weight,
      },
      totalScore,
      maxPossibleNow,
      currentPercentage:
        maxPossibleNow > 0 ? (totalScore / maxPossibleNow) * 100 : 0,
    };
  }, [mySubmissions, assignmentsData, syllabusConfig]);

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
          <div className="flex flex-wrap items-center justify-between gap-4">
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

            {/* Group Selector for Teachers */}
            {isTeacher && courseGroups.length > 0 && (
              <div className="flex items-center gap-2">
                <UsersRound className="h-4 w-4 text-muted-foreground" />
                <Select
                  value={selectedGroupId || undefined}
                  onValueChange={setSelectedGroupId}
                >
                  <SelectTrigger className="w-[140px] h-9">
                    <SelectValue placeholder="Select Group" />
                  </SelectTrigger>
                  <SelectContent>
                    {courseGroups.map((assignment) => (
                      <SelectItem
                        key={assignment.group_id}
                        value={assignment.group_id}
                      >
                        {assignment.group_code}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <Badge variant="outline" className="ml-1">
                  {courseGroups.length}{" "}
                  {courseGroups.length === 1 ? "group" : "groups"}
                </Badge>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Tabs - Simplified to Course, Grades, Attendance */}
      <Tabs defaultValue="course" className="space-y-4">
        <TabsList className="flex flex-wrap h-auto gap-2">
          <TabsTrigger value="course">
            <BookOpen className="h-4 w-4 mr-2" />
            Course
          </TabsTrigger>
          <TabsTrigger value="grades">
            <Award className="h-4 w-4 mr-2" />
            Grades
          </TabsTrigger>
          <TabsTrigger value="attendance">
            <CheckCircle2 className="h-4 w-4 mr-2" />
            Attendance
          </TabsTrigger>
        </TabsList>

        {/* Course Content Tab - Moodle Style */}
        <TabsContent value="course" className="space-y-4">
          {/* General Course Materials (Syllabus, Resources, etc.) */}
          {(() => {
            const generalContent = getGeneralContent();
            const hasGeneralContent =
              generalContent.assignments.length > 0 ||
              generalContent.lectures.length > 0;

            return hasGeneralContent ? (
              <Card className="border-primary/20">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <BookOpen className="h-5 w-5 text-primary" />
                      <CardTitle className="text-lg">
                        Course Materials
                      </CardTitle>
                    </div>
                    {isTeacher && (
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => setAddContentWeek(0)}
                      >
                        <Plus className="h-4 w-4 mr-1" />
                        Add
                      </Button>
                    )}
                  </div>
                </CardHeader>
                <CardContent className="space-y-2">
                  {generalContent.assignments.map((item) => (
                    <ContentItem
                      key={item.id}
                      item={item}
                      type="assignment"
                      courseId={courseId}
                      isStudent={isStudent}
                      isTeacher={isTeacher}
                      getSubmissionForAssignment={getSubmissionForAssignment}
                    />
                  ))}
                  {generalContent.lectures.map((item) => (
                    <ContentItem
                      key={item.id}
                      item={item}
                      type="lecture"
                      courseId={courseId}
                      isStudent={isStudent}
                      isTeacher={isTeacher}
                    />
                  ))}
                </CardContent>
              </Card>
            ) : isTeacher ? (
              <Card className="border-dashed border-primary/30">
                <CardContent className="py-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <BookOpen className="h-5 w-5" />
                      <span>Course Materials (Syllabus, Resources)</span>
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setAddContentWeek(0)}
                    >
                      <Plus className="h-4 w-4 mr-1" />
                      Add
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ) : null;
          })()}

          {/* Weekly Sections */}
          {assignmentsLoading || scheduleLoading ? (
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-16 w-full" />
              ))}
            </div>
          ) : (
            <div className="space-y-3">
              {Array.from({ length: TOTAL_WEEKS }, (_, i) => i + 1).map(
                (weekNum) => {
                  const weekContent = getWeekContent(weekNum);
                  const hasContent =
                    weekContent.assignments.length > 0 ||
                    weekContent.lectures.length > 0;
                  const isExpanded = expandedWeeks.has(weekNum);

                  return (
                    <Collapsible
                      key={weekNum}
                      open={isExpanded}
                      onOpenChange={() => toggleWeek(weekNum)}
                    >
                      <Card className={hasContent ? "" : "border-dashed"}>
                        <CollapsibleTrigger asChild>
                          <CardHeader className="py-3 cursor-pointer hover:bg-muted/50 transition-colors">
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-3">
                                {isExpanded ? (
                                  <ChevronDown className="h-5 w-5 text-muted-foreground" />
                                ) : (
                                  <ChevronRight className="h-5 w-5 text-muted-foreground" />
                                )}
                                <CardTitle className="text-base font-medium">
                                  Week {weekNum}
                                </CardTitle>
                                {hasContent && (
                                  <Badge
                                    variant="secondary"
                                    className="text-xs"
                                  >
                                    {weekContent.assignments.length +
                                      weekContent.lectures.length}{" "}
                                    items
                                  </Badge>
                                )}
                              </div>
                              {isTeacher && isExpanded && (
                                <Button
                                  size="sm"
                                  variant="outline"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    setAddContentWeek(weekNum);
                                  }}
                                >
                                  <Plus className="h-4 w-4 mr-1" />
                                  Add content
                                </Button>
                              )}
                            </div>
                          </CardHeader>
                        </CollapsibleTrigger>
                        <CollapsibleContent>
                          <CardContent className="pt-0 pb-4">
                            {hasContent ? (
                              <div className="space-y-2 border-l-2 border-muted ml-2 pl-4">
                                {weekContent.assignments.map((item) => (
                                  <ContentItem
                                    key={item.id}
                                    item={item}
                                    type="assignment"
                                    courseId={courseId}
                                    isStudent={isStudent}
                                    isTeacher={isTeacher}
                                    getSubmissionForAssignment={
                                      getSubmissionForAssignment
                                    }
                                  />
                                ))}
                                {weekContent.lectures.map((item) => (
                                  <ContentItem
                                    key={item.id}
                                    item={item}
                                    type="lecture"
                                    courseId={courseId}
                                    isStudent={isStudent}
                                    isTeacher={isTeacher}
                                  />
                                ))}
                              </div>
                            ) : (
                              <div className="text-center py-4 text-muted-foreground text-sm">
                                No content for this week
                                {isTeacher && (
                                  <Button
                                    size="sm"
                                    variant="link"
                                    className="ml-2"
                                    onClick={() => setAddContentWeek(weekNum)}
                                  >
                                    Add content
                                  </Button>
                                )}
                              </div>
                            )}
                          </CardContent>
                        </CollapsibleContent>
                      </Card>
                    </Collapsible>
                  );
                },
              )}
            </div>
          )}

          {/* Create Content Dialog */}
          <CreateContentDialog
            courseId={courseId}
            weekNumber={addContentWeek}
            open={addContentWeek !== null}
            onOpenChange={(open) => {
              if (!open) setAddContentWeek(null);
            }}
          />
        </TabsContent>

        <TabsContent value="grades" className="space-y-4">
          {/* Syllabus Configuration for Teachers */}
          {isTeacher && (
            <div className="flex justify-end">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setSyllabusDialogOpen(true)}
              >
                <Settings className="h-4 w-4 mr-2" />
                Configure Syllabus
              </Button>
            </div>
          )}

          {/* Grading Structure Overview */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Percent className="h-5 w-5" />
                Grading Structure
              </CardTitle>
              <CardDescription>
                Based on course syllabus: Register Midterm (
                {syllabusConfig.registerMidterm.weight}%) + Register Endterm (
                {syllabusConfig.registerEndterm.weight}%) + Final (
                {syllabusConfig.final.weight}%) = 100%
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-3">
                <div className="p-4 rounded-lg border bg-blue-50 dark:bg-blue-900/20">
                  <div className="text-sm font-medium text-blue-700 dark:text-blue-400">
                    Register Midterm
                  </div>
                  <div className="text-2xl font-bold text-blue-800 dark:text-blue-300">
                    {syllabusConfig.registerMidterm.weight}%
                  </div>
                  <div className="text-xs text-blue-600 dark:text-blue-400 mt-1">
                    {syllabusConfig.registerMidterm.components
                      .map((c) => `${c.count} ${c.name}(s)`)
                      .join(", ")}
                  </div>
                </div>
                <div className="p-4 rounded-lg border bg-purple-50 dark:bg-purple-900/20">
                  <div className="text-sm font-medium text-purple-700 dark:text-purple-400">
                    Register Endterm
                  </div>
                  <div className="text-2xl font-bold text-purple-800 dark:text-purple-300">
                    {syllabusConfig.registerEndterm.weight}%
                  </div>
                  <div className="text-xs text-purple-600 dark:text-purple-400 mt-1">
                    {syllabusConfig.registerEndterm.components
                      .map((c) => `${c.count} ${c.name}(s)`)
                      .join(", ")}
                  </div>
                </div>
                <div className="p-4 rounded-lg border bg-green-50 dark:bg-green-900/20">
                  <div className="text-sm font-medium text-green-700 dark:text-green-400">
                    Final Exam
                  </div>
                  <div className="text-2xl font-bold text-green-800 dark:text-green-300">
                    {syllabusConfig.final.weight}%
                  </div>
                  <div className="text-xs text-green-600 dark:text-green-400 mt-1">
                    Written/Oral Exam
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Current Grades */}
          {isStudent && calculateGrades && (
            <>
              {/* Overall Progress */}
              <Card>
                <CardHeader>
                  <CardTitle>Your Progress</CardTitle>
                  <CardDescription>Current grade calculation</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-6">
                    {/* Overall Score */}
                    <div className="flex justify-between items-center p-4 rounded-lg bg-muted">
                      <div>
                        <span className="font-medium">Current Grade</span>
                        <p className="text-sm text-muted-foreground">
                          Based on {calculateGrades.gradedSubmissions.length}{" "}
                          graded items
                        </p>
                      </div>
                      <div className="text-right">
                        <span
                          className={`text-3xl font-bold ${
                            calculateGrades.currentPercentage >= 90
                              ? "text-green-600"
                              : calculateGrades.currentPercentage >= 70
                                ? "text-blue-600"
                                : calculateGrades.currentPercentage >= 50
                                  ? "text-yellow-600"
                                  : "text-red-600"
                          }`}
                        >
                          {calculateGrades.currentPercentage.toFixed(1)}%
                        </span>
                        <p className="text-sm text-muted-foreground">
                          {calculateGrades.totalScore.toFixed(1)} /{" "}
                          {calculateGrades.maxPossibleNow.toFixed(1)} points
                          earned
                        </p>
                      </div>
                    </div>

                    {/* Progress Bar */}
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>Overall Progress</span>
                        <span>{calculateGrades.totalScore.toFixed(1)}%</span>
                      </div>
                      <Progress
                        value={calculateGrades.totalScore}
                        className="h-3"
                      />
                    </div>

                    {/* Register Midterm */}
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <h4 className="font-medium text-blue-700 dark:text-blue-400">
                          Register Midterm (
                          {syllabusConfig.registerMidterm.weight}%)
                        </h4>
                        <Badge
                          variant="outline"
                          className="bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
                        >
                          {calculateGrades.registerMidterm.score.toFixed(1)} /{" "}
                          {calculateGrades.registerMidterm.maxWeight}%
                        </Badge>
                      </div>
                      {calculateGrades.registerMidterm.assignments.length >
                      0 ? (
                        <div className="space-y-2 pl-4 border-l-2 border-blue-200 dark:border-blue-800">
                          {calculateGrades.registerMidterm.assignments.map(
                            (s, idx) => (
                              <div
                                key={s.id}
                                className="flex justify-between items-center p-2 rounded bg-blue-50/50 dark:bg-blue-900/10"
                              >
                                <span className="text-sm">
                                  {s.assignment?.title ||
                                    `Assignment ${idx + 1}`}
                                </span>
                                <span
                                  className={`font-medium ${
                                    s.percentage >= 90
                                      ? "text-green-600"
                                      : s.percentage >= 70
                                        ? "text-blue-600"
                                        : s.percentage >= 50
                                          ? "text-yellow-600"
                                          : "text-red-600"
                                  }`}
                                >
                                  {s.percentage.toFixed(0)}%
                                </span>
                              </div>
                            ),
                          )}
                          <div className="flex justify-between items-center p-2 font-medium">
                            <span className="text-sm">Average</span>
                            <span>
                              {calculateGrades.registerMidterm.assignmentAvg.toFixed(
                                1,
                              )}
                              %
                            </span>
                          </div>
                        </div>
                      ) : (
                        <p className="text-sm text-muted-foreground pl-4">
                          No graded assignments yet
                        </p>
                      )}
                    </div>

                    {/* Register Endterm */}
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <h4 className="font-medium text-purple-700 dark:text-purple-400">
                          Register Endterm (
                          {syllabusConfig.registerEndterm.weight}%)
                        </h4>
                        <Badge
                          variant="outline"
                          className="bg-purple-50 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400"
                        >
                          {calculateGrades.registerEndterm.score.toFixed(1)} /{" "}
                          {calculateGrades.registerEndterm.maxWeight}%
                        </Badge>
                      </div>
                      {calculateGrades.registerEndterm.assignments.length >
                      0 ? (
                        <div className="space-y-2 pl-4 border-l-2 border-purple-200 dark:border-purple-800">
                          {calculateGrades.registerEndterm.assignments.map(
                            (s, idx) => (
                              <div
                                key={s.id}
                                className="flex justify-between items-center p-2 rounded bg-purple-50/50 dark:bg-purple-900/10"
                              >
                                <span className="text-sm">
                                  {s.assignment?.title ||
                                    `Assignment ${idx + 1}`}
                                </span>
                                <span
                                  className={`font-medium ${
                                    s.percentage >= 90
                                      ? "text-green-600"
                                      : s.percentage >= 70
                                        ? "text-blue-600"
                                        : s.percentage >= 50
                                          ? "text-yellow-600"
                                          : "text-red-600"
                                  }`}
                                >
                                  {s.percentage.toFixed(0)}%
                                </span>
                              </div>
                            ),
                          )}
                          <div className="flex justify-between items-center p-2 font-medium">
                            <span className="text-sm">Average</span>
                            <span>
                              {calculateGrades.registerEndterm.assignmentAvg.toFixed(
                                1,
                              )}
                              %
                            </span>
                          </div>
                        </div>
                      ) : (
                        <p className="text-sm text-muted-foreground pl-4">
                          No graded assignments yet
                        </p>
                      )}
                    </div>

                    {/* Final */}
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <h4 className="font-medium text-green-700 dark:text-green-400">
                          Final Exam ({syllabusConfig.final.weight}%)
                        </h4>
                        <Badge
                          variant="outline"
                          className="bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-400"
                        >
                          {calculateGrades.final.score.toFixed(1)} /{" "}
                          {calculateGrades.final.maxWeight}%
                        </Badge>
                      </div>
                      <p className="text-sm text-muted-foreground pl-4">
                        Not yet scheduled
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </>
          )}

          {/* Teacher View - All Students' Grades Summary */}
          {isTeacher && (
            <Card>
              <CardHeader>
                <CardTitle>Students&apos; Grades Overview</CardTitle>
                <CardDescription>
                  Grade distribution based on syllabus
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  Configure the syllabus grading structure to automatically
                  calculate student grades.
                </p>
              </CardContent>
            </Card>
          )}

          {/* Syllabus Config Dialog */}
          <SyllabusConfigDialog
            open={syllabusDialogOpen}
            onOpenChange={setSyllabusDialogOpen}
            config={syllabusConfig}
            onSave={setSyllabusConfig}
          />
        </TabsContent>

        <TabsContent value="attendance" className="space-y-4">
          {isStudent && <StudentAttendanceView courseId={courseId} />}
          {isTeacher && (
            <TeacherAttendanceView
              courseId={courseId}
              courseGroups={courseGroups}
            />
          )}
        </TabsContent>
      </Tabs>
    </motion.div>
  );
}

// Content Item Component
interface ContentItemProps {
  item: {
    id: string;
    title: string;
    description?: string;
    due_at?: string;
    max_points?: number;
    start_time?: string;
    meeting_url?: string;
    is_online?: boolean;
    file_url?: string;
    course_id?: string;
    allow_late?: boolean;
    weight_percentage?: number;
    grading_category?: string;
    created_at?: string;
  };
  type: "assignment" | "lecture";
  courseId: string;
  isStudent?: boolean;
  isTeacher?: boolean;
  getSubmissionForAssignment?: (id: string) =>
    | {
        status?: string;
        grade_score?: number;
      }
    | undefined;
}

function ContentItem({
  item,
  type,
  courseId,
  isStudent,
  isTeacher,
  getSubmissionForAssignment,
}: ContentItemProps) {
  const contentType =
    type === "assignment" ? getContentType(item.description) : "lecture";
  const config = contentTypeConfig[contentType];
  const Icon = config.icon;

  // Get submission status for assignments
  const submission =
    type === "assignment" && isStudent && getSubmissionForAssignment
      ? getSubmissionForAssignment(item.id)
      : null;
  const isGraded = submission?.status === "graded";
  const gradePercentage =
    isGraded && submission?.grade_score !== undefined && item.max_points
      ? Math.round((submission.grade_score / item.max_points) * 100)
      : null;

  // Clean title (remove week prefix if present for display)
  const cleanTitle = item.title.replace(/^Week\s+\d+:\s*/i, "");

  // Determine if it's a link type
  const isLink = contentType === "link";
  const linkUrl =
    isLink && item.description
      ? item.description.replace("[LINK]", "").split("\n")[0].trim()
      : null;

  // Check if there's an attached file (from file_url field or legacy [FILE:] tag)
  const fileName = item.file_url ? null : getFileName(item.description);
  const directFileUrl = item.file_url || null;

  // Only try uploadService for legacy items that have [FILE:] in description
  const { data: attachedFiles } = useQuery({
    queryKey: ["assignment-files", item.id],
    queryFn: () => uploadService.getByReference("assignment", item.id),
    enabled: !!fileName && !directFileUrl,
    staleTime: 60000,
  });

  const attachedFile = attachedFiles?.[0];

  return (
    <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-muted/50 transition-colors group">
      <div
        className={`flex h-8 w-8 items-center justify-center rounded ${config.bgColor}`}
      >
        <Icon className={`h-4 w-4 ${config.color}`} />
      </div>
      <div className="flex-1 min-w-0">
        {isLink && linkUrl ? (
          <a
            href={linkUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="font-medium text-sm hover:text-primary flex items-center gap-1"
          >
            {cleanTitle}
            <ExternalLink className="h-3 w-3" />
          </a>
        ) : type === "assignment" &&
          contentType !== "resource" &&
          contentType !== "lecture" ? (
          <Link
            href={`/app/courses/${courseId}/assignments/${item.id}`}
            className="font-medium text-sm hover:text-primary"
          >
            {cleanTitle}
          </Link>
        ) : (
          <span className="font-medium text-sm">{cleanTitle}</span>
        )}
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          {item.due_at && (
            <span className="flex items-center gap-1">
              <Clock className="h-3 w-3" />
              Due: {format(new Date(item.due_at), "MMM d")}
            </span>
          )}
          {item.start_time && (
            <span className="flex items-center gap-1">
              <Calendar className="h-3 w-3" />
              {format(new Date(item.start_time), "MMM d, HH:mm")}
            </span>
          )}
          {/* Show file indicator */}
          {(directFileUrl || fileName || attachedFile) && (
            <span className="flex items-center gap-1 text-primary">
              <File className="h-3 w-3" />
              {attachedFile?.original_name || fileName || "Attachment"}
            </span>
          )}
        </div>
      </div>
      <div className="flex items-center gap-2">
        {/* Download button for direct file_url on the assignment */}
        {directFileUrl && (
          <Button variant="ghost" size="sm" className="h-7 text-xs" asChild>
            <a href={directFileUrl} target="_blank" rel="noopener noreferrer">
              <Download className="h-3 w-3 mr-1" />
              Download
            </a>
          </Button>
        )}
        {/* Download button for legacy attached files from upload service */}
        {!directFileUrl && attachedFile && (
          <Button variant="ghost" size="sm" className="h-7 text-xs" asChild>
            <a
              href={attachedFile.secure_url || attachedFile.url}
              target="_blank"
              rel="noopener noreferrer"
              download={attachedFile.original_name}
            >
              <Download className="h-3 w-3 mr-1" />
              Download
            </a>
          </Button>
        )}
        {/* Status badges for students */}
        {isStudent &&
          type === "assignment" &&
          contentType !== "resource" &&
          contentType !== "lecture" && (
            <>
              {isGraded && gradePercentage !== null && (
                <Badge
                  variant="outline"
                  className={`text-xs ${
                    gradePercentage >= 90
                      ? "bg-green-100 text-green-700 border-green-200 dark:bg-green-900/30 dark:text-green-400"
                      : gradePercentage >= 70
                        ? "bg-blue-100 text-blue-700 border-blue-200 dark:bg-blue-900/30 dark:text-blue-400"
                        : gradePercentage >= 50
                          ? "bg-yellow-100 text-yellow-700 border-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400"
                          : "bg-red-100 text-red-700 border-red-200 dark:bg-red-900/30 dark:text-red-400"
                  }`}
                >
                  <Award className="h-3 w-3 mr-1" />
                  {gradePercentage}%
                </Badge>
              )}
              {submission && !isGraded && (
                <Badge
                  variant="outline"
                  className="text-xs bg-blue-100 text-blue-700 border-blue-200 dark:bg-blue-900/30 dark:text-blue-400"
                >
                  <CheckCircle2 className="h-3 w-3 mr-1" />
                  Submitted
                </Badge>
              )}
              {!submission && (
                <Badge
                  variant="outline"
                  className="text-xs bg-orange-100 text-orange-700 border-orange-200 dark:bg-orange-900/30 dark:text-orange-400"
                >
                  <AlertCircle className="h-3 w-3 mr-1" />
                  Pending
                </Badge>
              )}
            </>
          )}
        {/* Join link for online lectures */}
        {type === "lecture" && item.is_online && item.meeting_url && (
          <a
            href={item.meeting_url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-primary hover:underline flex items-center gap-1"
          >
            <Video className="h-3 w-3" />
            Join
          </a>
        )}
        {/* Edit/Delete for teachers */}
        {isTeacher && type === "assignment" && (
          <EditAssignmentDialog
            assignment={
              {
                id: item.id,
                title: item.title,
                description: item.description || "",
                due_at: item.due_at,
                max_points: item.max_points || 0,
                course_id: courseId,
                allow_late: item.allow_late || false,
                weight_percentage: item.weight_percentage || 0,
                grading_category: item.grading_category || "",
                created_at: item.created_at || "",
                file_url: item.file_url,
              } as any
            }
            courseId={courseId}
            trigger={
              <Button
                variant="ghost"
                size="sm"
                className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <Pencil className="h-3.5 w-3.5" />
              </Button>
            }
          />
        )}
      </div>
    </div>
  );
}
