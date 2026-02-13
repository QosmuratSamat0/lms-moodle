"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import {
  LayoutDashboard,
  BookOpen,
  Calendar,
  FileText,
  MessageSquare,
  Bell,
  User,
  Settings,
  ChevronLeft,
  GraduationCap,
} from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useUIStore } from "@/store/ui-store";

const sidebarLinks = [
  {
    title: "Main",
    links: [
      { href: "/app/dashboard", label: "Dashboard", icon: LayoutDashboard },
      { href: "/app/courses", label: "Courses", icon: BookOpen },
      { href: "/app/deadlines", label: "Deadlines", icon: Calendar },
      { href: "/app/assignments", label: "Assignments", icon: FileText },
    ],
  },
  {
    title: "Communication",
    links: [
      { href: "/app/chat", label: "Chat", icon: MessageSquare },
      { href: "/app/notifications", label: "Notifications", icon: Bell },
    ],
  },
  {
    title: "Account",
    links: [
      { href: "/app/profile", label: "Profile", icon: User },
      { href: "/app/settings", label: "Settings", icon: Settings },
    ],
  },
];

export function SidebarNav() {
  const pathname = usePathname();
  const { sidebarCollapsed, toggleSidebarCollapsed } = useUIStore();

  return (
    <motion.aside
      initial={false}
      animate={{ width: sidebarCollapsed ? 72 : 256 }}
      transition={{ duration: 0.2, ease: "easeInOut" }}
      className="fixed left-0 top-0 z-40 hidden h-screen border-r bg-sidebar lg:block"
    >
      <div className="flex h-full flex-col">
        {/* Logo */}
        <div className="flex h-16 items-center border-b px-4">
          <Link href="/app/dashboard" className="flex items-center gap-2">
            <GraduationCap className="h-6 w-6 text-sidebar-primary" />
            <AnimatePresence>
              {!sidebarCollapsed && (
                <motion.span
                  initial={{ opacity: 0, width: 0 }}
                  animate={{ opacity: 1, width: "auto" }}
                  exit={{ opacity: 0, width: 0 }}
                  className="font-bold text-xl overflow-hidden whitespace-nowrap"
                >
                  LMS
                </motion.span>
              )}
            </AnimatePresence>
          </Link>
        </div>

        {/* Navigation */}
        <ScrollArea className="flex-1 px-3 py-4">
          <TooltipProvider delayDuration={0}>
            <nav className="space-y-6">
              {sidebarLinks.map((section) => (
                <div key={section.title} className="space-y-1">
                  <AnimatePresence>
                    {!sidebarCollapsed && (
                      <motion.h4
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="px-2 text-xs font-semibold uppercase tracking-wider text-sidebar-foreground/60"
                      >
                        {section.title}
                      </motion.h4>
                    )}
                  </AnimatePresence>
                  {section.links.map((link) => {
                    const isActive =
                      pathname === link.href ||
                      pathname.startsWith(link.href + "/");
                    const LinkContent = (
                      <Link
                        href={link.href}
                        className={cn(
                          "flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors",
                          isActive
                            ? "bg-sidebar-accent text-sidebar-accent-foreground"
                            : "text-sidebar-foreground/80 hover:bg-sidebar-accent/50 hover:text-sidebar-accent-foreground",
                        )}
                      >
                        <link.icon className="h-4 w-4 shrink-0" />
                        <AnimatePresence>
                          {!sidebarCollapsed && (
                            <motion.span
                              initial={{ opacity: 0, width: 0 }}
                              animate={{ opacity: 1, width: "auto" }}
                              exit={{ opacity: 0, width: 0 }}
                              className="overflow-hidden whitespace-nowrap"
                            >
                              {link.label}
                            </motion.span>
                          )}
                        </AnimatePresence>
                      </Link>
                    );

                    if (sidebarCollapsed) {
                      return (
                        <Tooltip key={link.href}>
                          <TooltipTrigger asChild>{LinkContent}</TooltipTrigger>
                          <TooltipContent side="right">
                            {link.label}
                          </TooltipContent>
                        </Tooltip>
                      );
                    }

                    return <div key={link.href}>{LinkContent}</div>;
                  })}
                </div>
              ))}
            </nav>
          </TooltipProvider>
        </ScrollArea>

        {/* Collapse Button */}
        <div className="border-t p-3">
          <Button
            variant="ghost"
            size="sm"
            onClick={toggleSidebarCollapsed}
            className="w-full justify-start"
          >
            <motion.div
              animate={{ rotate: sidebarCollapsed ? 180 : 0 }}
              transition={{ duration: 0.2 }}
            >
              <ChevronLeft className="h-4 w-4" />
            </motion.div>
            <AnimatePresence>
              {!sidebarCollapsed && (
                <motion.span
                  initial={{ opacity: 0, width: 0 }}
                  animate={{ opacity: 1, width: "auto" }}
                  exit={{ opacity: 0, width: 0 }}
                  className="ml-2 overflow-hidden whitespace-nowrap"
                >
                  Collapse
                </motion.span>
              )}
            </AnimatePresence>
          </Button>
        </div>
      </div>
    </motion.aside>
  );
}
