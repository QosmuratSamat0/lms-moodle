"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";

import { loginSchema, type LoginFormData } from "@/schemas/auth";
import { useAuthStore } from "@/store/auth-store";
import { useUIStore } from "@/store/ui-store";
import authService from "@/services/auth";

export function LoginDialog() {
  const router = useRouter();
  const { login } = useAuthStore();
  const { loginModalOpen, setLoginModalOpen } = useUIStore();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
      remember: false,
    },
  });

  const loginMutation = useMutation({
    mutationFn: authService.login,
    onSuccess: (data) => {
      login(data.user, data.access_token, data.refresh_token);
      toast.success("Welcome back!");
      setLoginModalOpen(false);
      reset();
      router.push("/app/dashboard");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to login");
    },
  });

  const onSubmit = (data: LoginFormData) => {
    loginMutation.mutate({
      email: data.email,
      password: data.password,
    });
  };

  const handleOpenChange = (open: boolean) => {
    setLoginModalOpen(open);
    if (!open) {
      reset();
    }
  };

  return (
    <Dialog open={loginModalOpen} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="text-xl sm:text-2xl">Welcome back</DialogTitle>
          <DialogDescription>
            Enter your credentials to access your account
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 pt-2 sm:pt-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="name@example.com"
              className="h-11 sm:h-10 text-base sm:text-sm"
              autoComplete="email"
              autoCapitalize="none"
              autoCorrect="off"
              {...register("email")}
              disabled={loginMutation.isPending}
            />
            {errors.email && (
              <p className="text-sm text-destructive">{errors.email.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? "text" : "password"}
                placeholder="Enter your password"
                className="h-11 sm:h-10 text-base sm:text-sm pr-16"
                autoComplete="current-password"
                {...register("password")}
                disabled={loginMutation.isPending}
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent touch-manipulation"
                onClick={() => setShowPassword(!showPassword)}
              >
                {showPassword ? "Hide" : "Show"}
              </Button>
            </div>
            {errors.password && (
              <p className="text-sm text-destructive">
                {errors.password.message}
              </p>
            )}
          </div>
          <div className="flex items-center space-x-2 py-1">
            <Checkbox id="remember" {...register("remember")} className="h-5 w-5" />
            <Label htmlFor="remember" className="text-sm font-normal">
              Remember me
            </Label>
          </div>
          <Button
            type="submit"
            className="w-full h-11 sm:h-10 text-base sm:text-sm touch-manipulation"
            disabled={loginMutation.isPending}
          >
            {loginMutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            Sign in
          </Button>
        </form>
        <div className="mt-2 sm:mt-4 text-center text-sm text-muted-foreground pb-2 sm:pb-0">
          Contact your administrator if you need an account.
        </div>
      </DialogContent>
    </Dialog>
  );
}
