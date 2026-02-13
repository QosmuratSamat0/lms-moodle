import { BarChart, Students } from "lucide-react";

interface StudentStatsProps {
  totalCourses: number;
  overallGPA: number;
  isLoading?: boolean;
}

export function StudentStats({
  totalCourses,
  overallGPA,
  isLoading = false,
}: StudentStatsProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-2 gap-4">
        <div className="h-32 animate-pulse rounded-lg bg-muted" />
        <div className="h-32 animate-pulse rounded-lg bg-muted" />
      </div>
    );
  }

  const gpaColor =
    overallGPA >= 3.5
      ? "text-green-600 dark:text-green-400"
      : overallGPA >= 3.0
        ? "text-blue-600 dark:text-blue-400"
        : overallGPA >= 2.5
          ? "text-yellow-600 dark:text-yellow-400"
          : "text-red-600 dark:text-red-400";

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-muted-foreground">Total Courses</p>
            <p className="text-3xl font-bold mt-1">{totalCourses}</p>
          </div>
          <Students className="h-10 w-10 text-muted-foreground" />
        </div>
      </div>

      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-muted-foreground">Overall Grade</p>
            <p className={`text-3xl font-bold mt-1 ${gpaColor}`}>
              {overallGPA.toFixed(2)}
            </p>
          </div>
          <BarChart className="h-10 w-10 text-muted-foreground" />
        </div>
      </div>
    </div>
  );
}
