"use client";

import { useEffect, type ReactNode } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth-store";
import { tokenStorage } from "@/lib/api-client";
import authService from "@/services/auth";

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const { setUser, setLoading, logout } = useAuthStore();
  const router = useRouter();

  useEffect(() => {
    const initAuth = async () => {
      const token = tokenStorage.getAccessToken();

      if (!token) {
        setLoading(false);
        return;
      }

      try {
        const user = await authService.getMe();
        setUser(user);
      } catch {
        // Token is invalid or expired
        logout();
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, [setUser, setLoading, logout, router]);

  return <>{children}</>;
}
