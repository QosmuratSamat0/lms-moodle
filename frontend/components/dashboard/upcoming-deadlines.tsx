import { UpcomingDeadline } from "@/types";
import { formatDistanceToNow, parseISO } from "date-fns";
import { AlertCircle, Calendar } from "lucide-react";

interface UpcomingDeadlinesProps {
  deadlines: UpcomingDeadline[];
  isLoading?: boolean;
}

export function UpcomingDeadlines({
  deadlines,
  isLoading = false,
}: UpcomingDeadlinesProps) {
  if (isLoading) {
    return (
      <div className="rounded-lg border bg-card p-6">
        <h2 className="mb-4 text-lg font-semibold">Upcoming Deadlines</h2>
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
        <Calendar className="h-5 w-5" />
        Upcoming Deadlines
      </h2>

      {deadlines.length === 0 ? (
        <p className="text-sm text-muted-foreground">No upcoming deadlines</p>
      ) : (
        <div className="space-y-3">
          {deadlines.map((deadline) => (
            <div
              key={deadline.id}
              className={`flex items-start justify-between rounded-md border p-3 ${
                deadline.status === "overdue"
                  ? "border-red-200 bg-red-50 dark:border-red-900 dark:bg-red-950"
                  : deadline.days_remaining <= 3
                    ? "border-yellow-200 bg-yellow-50 dark:border-yellow-900 dark:bg-yellow-950"
                    : "border-green-200 bg-green-50 dark:border-green-900 dark:bg-green-950"
              }`}
            >
              <div className="flex-1">
                <p className="font-medium text-sm">{deadline.title}</p>
                <p className="text-xs text-muted-foreground">
                  {deadline.course_title}
                </p>
              </div>
              <div className="text-right">
                <p className="text-xs font-semibold">
                  {deadline.status === "overdue" ? (
                    <span className="text-red-600 dark:text-red-400 flex items-center gap-1">
                      <AlertCircle className="h-3 w-3" />
                      Overdue
                    </span>
                  ) : (
                    `${deadline.days_remaining}d left`
                  )}
                </p>
                <p className="text-xs text-muted-foreground">
                  {formatDistanceToNow(parseISO(deadline.due_at), {
                    addSuffix: true,
                  })}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
