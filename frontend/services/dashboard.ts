import { api } from "@/lib/api-client";
import { StudentDashboard } from "@/types";

const dashboardService = {
  /**
   * Get student dashboard with grades, deadlines, and course stats
   */
  getStudentDashboard: async (): Promise<StudentDashboard> => {
    return api.get("/dashboard/student");
  },

  /**
   * Get course-specific stats for a student
   */
  getCourseStat: async (courseId: string) => {
    return api.get(`/dashboard/student/courses/${courseId}`);
  },
};

export default dashboardService;
