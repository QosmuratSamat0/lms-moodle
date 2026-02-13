"use client";

import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { dashboardService } from "@/services";
import {
  UpcomingDeadlines,
  RecentGrades,
  CourseStatsCard,
  StudentStats,
} from "@/components/dashboard";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";

export default function DashboardPage() {
  const {
    data: dashboard,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["dashboard", "student"],
    queryFn: dashboardService.getStudentDashboard,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="h-16 animate-pulse rounded-lg bg-muted" />
        <div className="grid gap-6 md:grid-cols-2">
          <div className="h-32 animate-pulse rounded-lg bg-muted" />
          <div className="h-32 animate-pulse rounded-lg bg-muted" />
        </div>
        <div className="grid gap-6 md:grid-cols-2">
          <div className="h-64 animate-pulse rounded-lg bg-muted" />
          <div className="h-64 animate-pulse rounded-lg bg-muted" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          Failed to load dashboard. Please try again later.
        </AlertDescription>
      </Alert>
    );
  }

  if (!dashboard) {
    return (
      <Alert>
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>No dashboard data available.</AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome back! Here's an overview of your courses and grades.
        </p>
      </div>

      {/* Stats Overview */}
      <StudentStats
        totalCourses={dashboard.total_courses}
        overallGPA={dashboard.overall_gpa}
      />

      {/* Main Content Grid */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Left Column */}
        <div className="space-y-6">
          <UpcomingDeadlines deadlines={dashboard.upcoming_deadlines} />
          <RecentGrades grades={dashboard.recent_grades} />
        </div>

        {/* Right Column */}
        <div className="space-y-6">
          <CourseStatsCard courses={dashboard.course_stats} />
        </div>
      </div>
    </div>
  );
}
