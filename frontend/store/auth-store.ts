import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User, UserRole } from "@/types/user";
import { tokenStorage } from "@/lib/api-client";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  setUser: (user: User | null) => void;
  setLoading: (loading: boolean) => void;
  login: (user: User, accessToken: string, refreshToken: string, rememberMe?: boolean) => void;
  logout: () => void;
  hasRole: (roles: UserRole | UserRole[]) => boolean;
  isTeacher: () => boolean;
  isStudent: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,

      setUser: (user) =>
        set({
          user,
          isAuthenticated: !!user,
        }),

      setLoading: (isLoading) => set({ isLoading }),

      login: (user, accessToken, refreshToken, rememberMe = false) => {
        tokenStorage.setTokens(accessToken, refreshToken, rememberMe);
        set({
          user,
          isAuthenticated: true,
          isLoading: false,
        });
      },

      logout: () => {
        tokenStorage.clearTokens();
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },

      hasRole: (roles) => {
        const { user } = get();
        if (!user) return false;
        const roleArray = Array.isArray(roles) ? roles : [roles];
        return roleArray.includes(user.role);
      },

      isTeacher: () => {
        const { user } = get();
        return user?.role === "teacher";
      },

      isStudent: () => {
        const { user } = get();
        return user?.role === "student";
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
);
