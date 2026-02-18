import api from "@/lib/api-client";
import type { PaginationParams } from "@/types/common";

// Types
export type QuestionType = "multiple_choice" | "true_false" | "short_answer";

export interface Quiz {
  id: string;
  course_id: string;
  created_by: string;
  title: string;
  description?: string;
  time_limit_minutes?: number;
  max_attempts: number;
  start_date?: string;
  end_date?: string;
  published: boolean;
  created_at: string;
  updated_at: string;
}

export interface QuizQuestion {
  id: string;
  quiz_id: string;
  type: QuestionType;
  text: string;
  options?: string[];
  correct_answer: string;
  points: number;
  order_index: number;
}

export interface QuizAttempt {
  id: string;
  quiz_id: string;
  student_id: string;
  started_at: string;
  submitted_at?: string;
  score?: number;
  max_score: number;
  percentage?: number;
}

export interface AnswerResult {
  question_id: string;
  question_text: string;
  your_answer: string;
  correct_answer: string;
  is_correct: boolean;
  points_earned: number;
  max_points: number;
}

export interface QuizResult {
  attempt_id: string;
  quiz_id: string;
  quiz_title: string;
  student_id: string;
  score: number;
  max_score: number;
  percentage: number;
  passed: boolean;
  submitted_at: string;
  answers?: AnswerResult[];
}

// Input types
export interface QuestionInput {
  type: QuestionType;
  text: string;
  options?: string[];
  correct_answer: string;
  points: number;
}

export interface CreateQuizInput {
  course_id: string;
  title: string;
  description?: string;
  time_limit_minutes?: number;
  max_attempts?: number;
  start_date?: string;
  end_date?: string;
  questions: QuestionInput[];
}

export interface UpdateQuizInput {
  published?: boolean;
}

export interface AnswerInput {
  question_id: string;
  answer: string;
}

export interface SubmitQuizInput {
  answers: AnswerInput[];
}

// Response types
export interface QuizListResponse {
  quizzes: Quiz[];
  total: number;
}

export interface AttemptsListResponse {
  attempts: QuizAttempt[];
  total: number;
}

type QueryParams = Record<string, string | number | boolean | undefined>;

export const quizService = {
  // Get quiz by ID
  getById: async (id: string): Promise<Quiz> => {
    return api.get<Quiz>(`/quizzes/${id}`);
  },

  // Get quiz questions
  getQuestions: async (quizId: string): Promise<{ questions: QuizQuestion[] }> => {
    return api.get<{ questions: QuizQuestion[] }>(`/quizzes/${quizId}/questions`);
  },

  // List quizzes by course
  listByCourse: async (
    courseId: string,
    params?: PaginationParams
  ): Promise<QuizListResponse> => {
    return api.get<QuizListResponse>(`/courses/${courseId}/quizzes`, params as QueryParams);
  },

  // Create quiz (teacher/admin)
  create: async (data: CreateQuizInput): Promise<Quiz> => {
    return api.post<Quiz>("/quizzes", data);
  },

  // Update quiz (teacher/admin)
  update: async (id: string, data: UpdateQuizInput): Promise<Quiz> => {
    return api.put<Quiz>(`/quizzes/${id}`, data);
  },

  // Delete quiz (admin)
  delete: async (id: string): Promise<void> => {
    return api.delete(`/quizzes/${id}`);
  },

  // Start quiz attempt
  startAttempt: async (quizId: string): Promise<QuizAttempt> => {
    return api.post<QuizAttempt>(`/quizzes/${quizId}/attempts`);
  },

  // Get student's attempts for a quiz
  getAttempts: async (quizId: string): Promise<AttemptsListResponse> => {
    return api.get<AttemptsListResponse>(`/quizzes/${quizId}/attempts`);
  },

  // Submit quiz attempt
  submitAttempt: async (attemptId: string, data: SubmitQuizInput): Promise<QuizResult> => {
    return api.post<QuizResult>(`/quiz-attempts/${attemptId}/submit`, data);
  },

  // Get attempt result
  getAttemptResult: async (attemptId: string): Promise<QuizResult> => {
    return api.get<QuizResult>(`/quiz-attempts/${attemptId}/result`);
  },
};

export default quizService;
