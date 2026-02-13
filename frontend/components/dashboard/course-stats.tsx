import { CourseStats } from "@/types";
import { Progress } from "@/components/ui/progress";
import { BookOpen } from "lucide-react";

interface CourseStatsProps {
  courses: CourseStats[];
  isLoading?: boolean;
}

export function CourseStatsCard({
  courses,
  isLoading = false,
}: CourseStatsProps) {
  if (isLoading) {
    return (
      <div className="rounded-lg border bg-card p-6">
        <h2 className="mb-4 text-lg font-semibold">Courses</h2>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-24 animate-pulse rounded bg-muted" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-lg border bg-card p-6">
      <h2 className="mb-4 text-lg font-semibold flex items-center gap-2">
        <BookOpen className="h-5 w-5" />
        My Courses
      </h2>

      {courses.length === 0 ? (
        <p className="text-sm text-muted-foreground">No enrolled courses</p>
      ) : (
        <div className="space-y-4">
          {courses.map((course) => {
            const gradeColor =
              course.current_grade >= 90
                ? "text-green-600 dark:text-green-400"
                : course.current_grade >= 80
                  ? "text-blue-600 dark:text-blue-400"
                  : course.current_grade >= 70
                    ? "text-yellow-600 dark:text-yellow-400"
                    : "text-red-600 dark:text-red-400";

            return (
              <div
                key={course.course_id}
                className="space-y-2 rounded-md border p-3"
              >
                <div className="flex items-start justify-between">
                  <div>
                    <p className="font-medium text-sm">{course.course_title}</p>
                    <p className="text-xs text-muted-foreground">
                      {course.instructor_name}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className={`font-bold text-lg ${gradeColor}`}>
                      {course.letter_grade}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {Math.round(course.current_grade)}%
                    </p>
                  </div>
                </div>
                <div className="space-y-1">
                  <div className="flex justify-between text-xs">
                    <span className="text-muted-foreground">Progress</span>
                    <span className="text-muted-foreground">
                      {course.completed_assignments}/{course.total_assignments}
                    </span>
                  </div>
                  <Progress
                    value={course.progress_percentage}
                    className="h-2"
                  />
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
