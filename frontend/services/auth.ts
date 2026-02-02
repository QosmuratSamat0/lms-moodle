import api, { tokenStorage } from "@/lib/api-client";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RefreshTokenRequest,
  RefreshTokenResponse,
} from "@/types/auth";
import type { User } from "@/types/user";
import type { UpdateProfileData, ChangePasswordData } from "@/schemas/profile";

export const authService = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>("/auth/login", data);
    tokenStorage.setTokens(response.access_token, response.refresh_token);
    return response;
  },

  register: async (data: RegisterRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>("/auth/register", data);
    tokenStorage.setTokens(response.access_token, response.refresh_token);
    return response;
  },

  logout: async (): Promise<void> => {
    try {
      await api.post("/auth/logout");
    } finally {
      tokenStorage.clearTokens();
    }
  },

  refreshToken: async (): Promise<RefreshTokenResponse> => {
    const refreshToken = tokenStorage.getRefreshToken();
    if (!refreshToken) {
      throw new Error("No refresh token available");
    }

    const data: RefreshTokenRequest = { refresh_token: refreshToken };
    const response = await api.post<RefreshTokenResponse>(
      "/auth/refresh",
      data,
    );
    tokenStorage.setTokens(response.access_token, response.refresh_token);
    return response;
  },

  getMe: async (): Promise<User> => {
    return api.get<User>("/users/me");
  },

  updateProfile: async (data: UpdateProfileData): Promise<User> => {
    return api.put<User>("/users/me", data);
  },

  changePassword: async (data: ChangePasswordData): Promise<void> => {
    await api.put("/users/me/password", {
      old_password: data.current_password,
      new_password: data.new_password,
    });
  },
};

// Mock implementations for development
export const authServiceMock = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 500));

    const mockUser: User = {
      id: "1",
      email: data.email,
      role: data.email.includes("teacher") ? "teacher" : "student",
      is_active: true,
      first_name: "John",
      last_name: "Doe",
      created_at: new Date().toISOString(),
    };

    const response: LoginResponse = {
      access_token: "mock_access_token_" + Date.now(),
      refresh_token: "mock_refresh_token_" + Date.now(),
      user: mockUser,
    };

    tokenStorage.setTokens(response.access_token, response.refresh_token);
    return response;
  },

  register: async (data: RegisterRequest): Promise<LoginResponse> => {
    await new Promise((resolve) => setTimeout(resolve, 500));

    const mockUser: User = {
      id: "1",
      email: data.email,
      role: data.role,
      is_active: true,
      first_name: data.first_name,
      last_name: data.last_name,
      group_name: data.group_name,
      department: data.department,
      created_at: new Date().toISOString(),
    };

    const response: LoginResponse = {
      access_token: "mock_access_token_" + Date.now(),
      refresh_token: "mock_refresh_token_" + Date.now(),
      user: mockUser,
    };

    tokenStorage.setTokens(response.access_token, response.refresh_token);
    return response;
  },

  logout: async (): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 200));
    tokenStorage.clearTokens();
  },

  getMe: async (): Promise<User> => {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return {
      id: "1",
      email: "user@example.com",
      role: "student",
      is_active: true,
      first_name: "John",
      last_name: "Doe",
      created_at: new Date().toISOString(),
    };
  },

  updateProfile: async (data: UpdateProfileData): Promise<User> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    return {
      id: "1",
      email: data.email || "user@example.com",
      role: "student",
      is_active: true,
      first_name: data.first_name,
      last_name: data.last_name,
      group_name: data.group_name,
      department: data.department,
      created_at: new Date().toISOString(),
    };
  },

  changePassword: async (data: ChangePasswordData): Promise<void> => {
    await new Promise((resolve) => setTimeout(resolve, 500));
    // Mock password change - in real implementation would validate current password
    console.log(
      "Mock password change for:",
      data.current_password ? "user" : "unknown",
    );
  },
};

// Export the service to use (switch between real and mock)
const useMock = process.env.NEXT_PUBLIC_USE_MOCK === "true";
export default useMock ? authServiceMock : authService;
