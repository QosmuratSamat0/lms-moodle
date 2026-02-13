export interface UpcomingDeadline {
  id: string;
  title: string;
  course_title: string;
  due_at: string;
  days_remaining: number;
  status: "pending" | "completed" | "overdue";
}

export interface RecentSubmission {
  id: string;
  assignment_title: string;
  course_title: string;
  submitted_at: string;
  status: "pending" | "graded";
}

export interface RecentGrade {
  id: string;
  assignment_title: string;
  course_title: string;
  score: number;
  max_points: number;
  percentage: number;
  graded_at: string;
}

export interface GradeTrend {
  month: string;
  average_percentage: number;
}

export interface CourseStats {
  course_id: string;
  course_title: string;
  instructor_name: string;
  total_assignments: number;
  completed_assignments: number;
  current_grade: number;
  letter_grade: string;
  progress_percentage: number;
}

export interface CourseProgress {
  course_id: string;
  course_title: string;
  completed_modules: number;
  total_modules: number;
  progress_percentage: number;
}

export interface StudentDashboard {
  upcoming_deadlines: UpcomingDeadline[];
  recent_submissions: RecentSubmission[];
  recent_grades: RecentGrade[];
  grade_trends: GradeTrend[];
  course_stats: CourseStats[];
  course_progress: CourseProgress[];
  total_courses: number;
  overall_gpa: number;
}
