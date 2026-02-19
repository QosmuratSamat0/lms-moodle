import { User } from "./user";

export interface LoginRequest {
  email: string;
  password: string;
  remember?: boolean;
}

export interface RegisterRequest {
  email: string;
  password: string;
  role: "student" | "teacher";
  first_name?: string;
  last_name?: string;
  group_name?: string; // for students
  department?: string; // for teachers
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}
