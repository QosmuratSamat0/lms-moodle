"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion } from "framer-motion";
import { Menu, GraduationCap, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetTitle,
  SheetTrigger,
  SheetClose,
} from "@/components/ui/sheet";
import { ThemeToggle } from "@/components/common/theme-toggle";
import { useUIStore } from "@/store/ui-store";
import { cn } from "@/lib/utils";

const navLinks = [
  { href: "/", label: "Home" },
  { href: "/about", label: "About" },
  { href: "/contact", label: "Contact" },
];

export function Navbar() {
  const pathname = usePathname();
  const { openLoginModal } = useUIStore();

  return (
    <motion.header
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 safe-area-top"
    >
      <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8 flex h-14 sm:h-16 items-center justify-between">
        <div className="flex items-center gap-4 md:gap-10">
          <Link href="/" className="flex items-center gap-2">
            <GraduationCap className="h-6 w-6 text-primary" />
            <span className="font-bold text-lg sm:text-xl">LMS</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex gap-6">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={cn(
                  "text-sm font-medium transition-colors hover:text-primary",
                  pathname === link.href
                    ? "text-foreground"
                    : "text-muted-foreground",
                )}
              >
                {link.label}
              </Link>
            ))}
          </nav>
        </div>

        {/* Desktop Actions */}
        <div className="hidden md:flex items-center gap-4">
          <ThemeToggle />
          <Button onClick={openLoginModal}>Login</Button>
        </div>

        {/* Mobile Navigation */}
        <div className="flex md:hidden items-center gap-1">
          <ThemeToggle />
          <Sheet>
            <SheetTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-10 w-10 touch-manipulation"
                aria-label="Open menu"
              >
                <Menu className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent
              side="right"
              showCloseButton={false}
            >
              <SheetTitle className="sr-only">Navigation Menu</SheetTitle>
              {/* Mobile Menu Header */}
              <div className="flex h-14 items-center justify-between border-b px-4">
                <SheetClose asChild>
                  <Link
                    href="/"
                    className="flex items-center gap-2"
                  >
                    <GraduationCap className="h-6 w-6 text-primary" />
                    <span className="font-bold text-lg">LMS</span>
                  </Link>
                </SheetClose>
                <SheetClose asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-10 w-10 touch-manipulation"
                  >
                    <X className="h-5 w-5" />
                    <span className="sr-only">Close menu</span>
                  </Button>
                </SheetClose>
              </div>

              {/* Mobile Menu Content */}
              <nav className="flex flex-col p-4 flex-1">
                {navLinks.map((link) => (
                  <SheetClose asChild key={link.href}>
                    <Link
                      href={link.href}
                      className={cn(
                        "flex items-center py-4 px-3 text-lg font-medium transition-colors hover:text-primary hover:bg-accent rounded-lg touch-manipulation",
                        pathname === link.href
                          ? "text-foreground bg-accent/50"
                          : "text-muted-foreground",
                      )}
                    >
                      {link.label}
                    </Link>
                  </SheetClose>
                ))}

                {/* Login Button */}
                <div className="mt-auto pt-6 border-t">
                  <SheetClose asChild>
                    <Button
                      className="w-full h-12 text-base touch-manipulation"
                      onClick={openLoginModal}
                    >
                      Login
                    </Button>
                  </SheetClose>
                </div>
              </nav>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </motion.header>
  );
}
