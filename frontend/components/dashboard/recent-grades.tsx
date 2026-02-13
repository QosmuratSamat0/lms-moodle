import { RecentGrade } from "@/types/dashboard";
import { formatDistanceToNow, parseISO } from "date-fns";
import { Award } from "lucide-react";

interface RecentGradesProps {
  grades: RecentGrade[];
  isLoading?: boolean;
}

export function RecentGrades({ grades, isLoading = false }: RecentGradesProps) {
  if (isLoading) {
    return (
      <div className="rounded-lg border bg-card p-6">
        <h2 className="mb-4 text-lg font-semibold">Recent Grades</h2>
        <div className="space-y-3">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-12 animate-pulse rounded bg-muted" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-lg border bg-card p-6">
      <h2 className="mb-4 text-lg font-semibold flex items-center gap-2">
        <Award className="h-5 w-5" />
        Recent Grades
      </h2>

      {grades.length === 0 ? (
        <p className="text-sm text-muted-foreground">No grades yet</p>
      ) : (
        <div className="space-y-3">
          {grades.slice(0, 5).map((grade) => {
            const percentage = Math.round(grade.percentage);
            const gradeColor =
              percentage >= 90
                ? "text-green-600 dark:text-green-400"
                : percentage >= 80
                  ? "text-blue-600 dark:text-blue-400"
                  : percentage >= 70
                    ? "text-yellow-600 dark:text-yellow-400"
                    : "text-red-600 dark:text-red-400";

            return (
              <div
                key={grade.id}
                className="flex items-center justify-between rounded-md border p-3"
              >
                <div className="flex-1">
                  <p className="font-medium text-sm">
                    {grade.assignment_title}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {grade.course_title} •{" "}
                    {formatDistanceToNow(parseISO(grade.graded_at), {
                      addSuffix: true,
                    })}
                  </p>
                </div>
                <div className="text-right">
                  <p className={`font-bold ${gradeColor}`}>
                    {grade.score}/{grade.max_points}
                  </p>
                  <p className="text-xs text-muted-foreground">{percentage}%</p>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
