"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { SidebarNav } from "./sidebar-nav";
import { TopBar } from "./top-bar";
import { MobileAppBanner } from "@/components/common/mobile-app-banner";
import { useAuthStore } from "@/store/auth-store";
import { useUIStore } from "@/store/ui-store";
import { cn } from "@/lib/utils";

interface AppShellProps {
  children: React.ReactNode;
}

export function AppShell({ children }: AppShellProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuthStore();
  const { sidebarCollapsed } = useUIStore();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center safe-area-inset">
        <div className="flex flex-col items-center gap-2">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
          <p className="text-sm text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Sidebar - hidden on mobile, fixed on desktop */}
      <SidebarNav />

      {/* Main content area */}
      <div
        className={cn(
          "min-h-screen transition-[margin] duration-200 ease-in-out",
          "lg:ml-64", // Default sidebar width on desktop
          sidebarCollapsed && "lg:ml-[72px]", // Collapsed sidebar width
        )}
      >
        <TopBar />
        <main className="flex-1 px-3 py-4 sm:px-4 sm:py-6 md:px-6 md:py-8 lg:px-8 lg:py-10 pb-20 sm:pb-6 safe-area-bottom">
          <div className="mx-auto w-full max-w-7xl">{children}</div>
        </main>
      </div>
      
      {/* Mobile App Install Banner - only shows on mobile */}
      <MobileAppBanner />
    </div>
  );
}
