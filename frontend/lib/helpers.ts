// Re-export cn from utils
export { cn } from "./utils";

// Date formatting utilities
export function formatDate(date: string | Date): string {
  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  }).format(new Date(date));
}

export function formatDateTime(date: string | Date): string {
  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(date));
}

export function formatTime(date: string | Date): string {
  return new Intl.DateTimeFormat("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(date));
}

export function formatRelativeTime(date: string | Date): string {
  const now = new Date();
  const then = new Date(date);
  const diffInSeconds = Math.floor((now.getTime() - then.getTime()) / 1000);

  if (diffInSeconds < 60) return "just now";
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`;
  if (diffInSeconds < 604800)
    return `${Math.floor(diffInSeconds / 86400)}d ago`;

  return formatDate(date);
}

export function isOverdue(dueDate: string | Date): boolean {
  return new Date(dueDate) < new Date();
}

export function getDaysUntil(date: string | Date): number {
  const now = new Date();
  const target = new Date(date);
  const diffInMs = target.getTime() - now.getTime();
  return Math.ceil(diffInMs / (1000 * 60 * 60 * 24));
}

export function getTimeRemaining(dueDate: string | Date): string {
  const now = new Date();
  const due = new Date(dueDate);
  const diffInMs = due.getTime() - now.getTime();

  if (diffInMs < 0) {
    // Overdue
    const overMs = Math.abs(diffInMs);
    const days = Math.floor(overMs / (1000 * 60 * 60 * 24));
    const hours = Math.floor(
      (overMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60),
    );
    const mins = Math.floor((overMs % (1000 * 60 * 60)) / (1000 * 60));

    if (days > 0)
      return `${days} day${days > 1 ? "s" : ""} ${hours} hour${hours !== 1 ? "s" : ""} overdue`;
    if (hours > 0)
      return `${hours} hour${hours > 1 ? "s" : ""} ${mins} min${mins !== 1 ? "s" : ""} overdue`;
    return `${mins} minute${mins !== 1 ? "s" : ""} overdue`;
  }

  const days = Math.floor(diffInMs / (1000 * 60 * 60 * 24));
  const hours = Math.floor(
    (diffInMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60),
  );
  const mins = Math.floor((diffInMs % (1000 * 60 * 60)) / (1000 * 60));

  if (days > 0)
    return `${days} day${days > 1 ? "s" : ""} ${hours} hour${hours !== 1 ? "s" : ""} remaining`;
  if (hours > 0)
    return `${hours} hour${hours > 1 ? "s" : ""} ${mins} min${mins !== 1 ? "s" : ""} remaining`;
  return `${mins} minute${mins !== 1 ? "s" : ""} remaining`;
}

export function getSubmissionTimeInfo(
  submittedAt: string,
  dueDate?: string,
): string {
  if (!dueDate) return `Submitted ${formatDateTime(submittedAt)}`;

  const submitted = new Date(submittedAt);
  const due = new Date(dueDate);
  const diffInMs = due.getTime() - submitted.getTime();

  if (diffInMs > 0) {
    // Early
    const days = Math.floor(diffInMs / (1000 * 60 * 60 * 24));
    const hours = Math.floor(
      (diffInMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60),
    );
    const mins = Math.floor((diffInMs % (1000 * 60 * 60)) / (1000 * 60));

    if (days > 0)
      return `Submitted ${days} day${days > 1 ? "s" : ""} ${hours} hour${hours !== 1 ? "s" : ""} early`;
    if (hours > 0)
      return `Submitted ${hours} hour${hours > 1 ? "s" : ""} ${mins} min${mins !== 1 ? "s" : ""} early`;
    return `Submitted ${mins} minute${mins !== 1 ? "s" : ""} early`;
  } else {
    // Late
    const lateMs = Math.abs(diffInMs);
    const days = Math.floor(lateMs / (1000 * 60 * 60 * 24));
    const hours = Math.floor(
      (lateMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60),
    );
    const mins = Math.floor((lateMs % (1000 * 60 * 60)) / (1000 * 60));

    if (days > 0)
      return `Submitted ${days} day${days > 1 ? "s" : ""} ${hours} hour${hours !== 1 ? "s" : ""} late`;
    if (hours > 0)
      return `Submitted ${hours} hour${hours > 1 ? "s" : ""} ${mins} min${mins !== 1 ? "s" : ""} late`;
    return `Submitted ${mins} minute${mins !== 1 ? "s" : ""} late`;
  }
}

// String utilities
export function truncate(str: string, length: number): string {
  if (str.length <= length) return str;
  return str.slice(0, length) + "...";
}

export function getInitials(firstName?: string, lastName?: string): string {
  const first = firstName?.charAt(0).toUpperCase() || "";
  const last = lastName?.charAt(0).toUpperCase() || "";
  return first + last || "?";
}

export function getFullName(firstName?: string, lastName?: string): string {
  return [firstName, lastName].filter(Boolean).join(" ") || "Unknown";
}

// Number utilities
export function formatPercentage(value: number): string {
  return `${Math.round(value * 100) / 100}%`;
}

export function formatScore(score: number, maxPoints: number): string {
  return `${score}/${maxPoints}`;
}

// Role utilities
export function getRoleLabel(role: string): string {
  const labels: Record<string, string> = {
    student: "Student",
    teacher: "Teacher",
    manager: "Manager",
    admin: "Administrator",
  };
  return labels[role] || role;
}

export function getRoleColor(role: string): string {
  const colors: Record<string, string> = {
    student: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
    teacher:
      "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
    manager:
      "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200",
    admin: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
  };
  return (
    colors[role] ||
    "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
  );
}

// Status utilities
export function getSubmissionStatusColor(status: string): string {
  const colors: Record<string, string> = {
    pending:
      "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200",
    submitted: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
    graded: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
    late: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
  };
  return colors[status] || "bg-gray-100 text-gray-800";
}

export function getEnrollmentStatusColor(status: string): string {
  const colors: Record<string, string> = {
    active: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
    pending:
      "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200",
    rejected: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
    dropped: "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200",
  };
  return colors[status] || "bg-gray-100 text-gray-800";
}
