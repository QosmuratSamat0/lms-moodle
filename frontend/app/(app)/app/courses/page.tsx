"use client";

import { useState, useMemo } from "react";
import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { BookOpen, Users } from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import { SearchInput } from "@/components/common/search-input";
import { EmptyState } from "@/components/common/empty-state";
import { GridSkeleton } from "@/components/common/skeletons";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { getFullName, formatDate } from "@/lib/helpers";
import { useAuthStore } from "@/store/auth-store";
import courseService from "@/services/courses";
import type { Course } from "@/types/course";
import { useTeacherProfile } from "@/hooks/use-profile";

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

export default function CoursesPage() {
  const [search, setSearch] = useState("");
  const { user } = useAuthStore();
  const isTeacher = user?.role === "teacher";
  const isStudent = user?.role === "student";

  // Fetch all courses
  const { data, isLoading, error } = useQuery({
    queryKey: ["courses", search],
    queryFn: () => courseService.list({ search: search || undefined }),
  });

  // Fetch teacher profile for filtering by owner_teacher_id
  const { data: teacherProfile } = useTeacherProfile();

  // Fetch student enrollments to filter courses
  const { data: enrollmentsData, isLoading: enrollmentsLoading } = useQuery({
    queryKey: ["enrollments", "student", user?.id],
    queryFn: () => courseService.getMyEnrollments(user!.id),
    enabled: isStudent && !!user?.id,
  });

  // Get set of enrolled course IDs for students
  const enrolledCourseIds = useMemo(() => {
    const enrollments = Array.isArray(enrollmentsData) ? enrollmentsData : [];
    return new Set(enrollments.map((e: { course_id: string }) => e.course_id));
  }, [enrollmentsData]);

  // Filter courses based on user role
  const filteredCourses = useMemo(() => {
    const courses = Array.isArray(data) ? data : [];

    // For teachers, show courses they own
    if (isTeacher) {
      if (teacherProfile?.id) {
        return courses.filter(
          (course: Course) => course.owner_teacher_id === teacherProfile.id,
        );
      }
      return [];
    }

    // For students, only show enrolled courses
    if (isStudent && enrolledCourseIds.size > 0) {
      return courses.filter((course: Course) => enrolledCourseIds.has(course.id));
    }
    if (isStudent) return [];

    return courses;
  }, [data, isTeacher, isStudent, teacherProfile, enrolledCourseIds]);

  const loading = isLoading || (isStudent && enrollmentsLoading);

  return (
    <div className="space-y-6">
      <PageHeader
        title="Courses"
        description={
          isTeacher ? "Your assigned courses" : "Browse and manage your courses"
        }
      />

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <SearchInput
          placeholder="Search courses..."
          value={search}
          onChange={setSearch}
          className="sm:max-w-sm"
        />
      </div>

      {loading ? (
        <GridSkeleton count={6} />
      ) : error ? (
        <EmptyState
          title="Error loading courses"
          description="There was a problem loading your courses. Please try again."
          action={
            <Button onClick={() => window.location.reload()}>Retry</Button>
          }
        />
      ) : filteredCourses.length === 0 ? (
        <EmptyState
          icon={<BookOpen className="h-8 w-8 text-muted-foreground" />}
          title="No courses found"
          description={
            search
              ? "No courses match your search. Try different keywords."
              : "You haven't enrolled in any courses yet."
          }
        />
      ) : (
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
        >
          {filteredCourses.map((course) => (
            <motion.div key={course.id} variants={itemVariants}>
              <Link href={`/app/courses/${course.id}`}>
                <Card className="h-full transition-all hover:shadow-lg hover:border-primary/50">
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                        <BookOpen className="h-6 w-6 text-primary" />
                      </div>
                      <Badge
                        variant={course.is_active ? "default" : "secondary"}
                      >
                        {course.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </div>
                    <CardTitle className="line-clamp-2">
                      {course.title}
                    </CardTitle>
                    <CardDescription className="line-clamp-2">
                      {course.description || "No description available"}
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Users className="h-4 w-4" />
                        <span>
                          {getFullName(
                            course.teacher_first_name,
                            course.teacher_last_name,
                          )}
                        </span>
                      </div>
                    </div>
                    <p className="mt-2 text-xs text-muted-foreground">
                      Updated {formatDate(course.updated_at)}
                    </p>
                  </CardContent>
                </Card>
              </Link>
            </motion.div>
          ))}
        </motion.div>
      )}
    </div>
  );
}
