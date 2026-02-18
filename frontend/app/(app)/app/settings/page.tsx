"use client";

import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { toast } from "sonner";
import {
  Settings,
  Globe,
  Clock,
  Download,
  Trash2,
  Shield,
  Key,
  Smartphone,
  Monitor,
} from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { tokenStorage } from "@/lib/api-client";

export default function SettingsPage() {
  const [rememberMe, setRememberMe] = useState(false);

  useEffect(() => {
    setRememberMe(tokenStorage.getRememberMe());
  }, []);

  const handleRememberMeChange = (checked: boolean) => {
    setRememberMe(checked);
    tokenStorage.setRememberMe(checked);
    toast.success(checked ? "Session will persist after browser close" : "Session will end when browser closes");
  };

  const handleExportData = () => {
    toast.info("Data export feature coming soon");
  };

  const handleClearCache = () => {
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth-storage");
      toast.success("Cache cleared successfully");
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <PageHeader
        title="Settings"
        description="Configure application preferences and manage your data"
      />

      {/* General Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Settings className="h-5 w-5" />
            General
          </CardTitle>
          <CardDescription>
            Basic application settings
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 sm:gap-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Globe className="h-4 w-4 text-muted-foreground" />
                <Label className="font-medium">Language</Label>
              </div>
              <p className="text-sm text-muted-foreground">
                Select your preferred language
              </p>
            </div>
            <Select defaultValue="en">
              <SelectTrigger className="w-full sm:w-[180px]">
                <SelectValue placeholder="Select language" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="en">English</SelectItem>
                <SelectItem value="ru">Русский</SelectItem>
                <SelectItem value="kk">Қазақша</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <Separator />

          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 sm:gap-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Clock className="h-4 w-4 text-muted-foreground" />
                <Label className="font-medium">Timezone</Label>
              </div>
              <p className="text-sm text-muted-foreground">
                Set your local timezone for dates
              </p>
            </div>
            <Select defaultValue="asia-almaty">
              <SelectTrigger className="w-full sm:w-[180px]">
                <SelectValue placeholder="Select timezone" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="asia-almaty">Asia/Almaty (GMT+5)</SelectItem>
                <SelectItem value="asia-astana">Asia/Astana (GMT+6)</SelectItem>
                <SelectItem value="europe-moscow">Europe/Moscow (GMT+3)</SelectItem>
                <SelectItem value="utc">UTC (GMT+0)</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Session Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Key className="h-5 w-5" />
            Session
          </CardTitle>
          <CardDescription>
            Manage your login session preferences
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 sm:gap-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Monitor className="h-4 w-4 text-muted-foreground" />
                <Label className="font-medium">Remember Me</Label>
              </div>
              <p className="text-sm text-muted-foreground">
                Stay logged in after closing the browser
              </p>
            </div>
            <Switch
              checked={rememberMe}
              onCheckedChange={handleRememberMeChange}
            />
          </div>

          <Separator />

          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 sm:gap-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Smartphone className="h-4 w-4 text-muted-foreground" />
                <Label className="font-medium">Active Sessions</Label>
              </div>
              <p className="text-sm text-muted-foreground">
                View and manage your active sessions
              </p>
            </div>
            <Button variant="outline" size="sm">
              View Sessions
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Privacy & Data */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Privacy & Data
          </CardTitle>
          <CardDescription>
            Manage your data and privacy settings
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Download className="h-4 w-4 text-muted-foreground" />
                <Label className="font-medium">Export Data</Label>
              </div>
              <p className="text-sm text-muted-foreground">
                Download a copy of your data
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleExportData}
              className="w-full sm:w-auto"
            >
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
          </div>

          <Separator />

          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Trash2 className="h-4 w-4 text-muted-foreground" />
                <Label className="font-medium">Clear Cache</Label>
              </div>
              <p className="text-sm text-muted-foreground">
                Clear locally stored data and preferences
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleClearCache}
              className="w-full sm:w-auto"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Clear
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* About Section */}
      <Card>
        <CardHeader>
          <CardTitle>About</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted-foreground">
          <p><strong>Mini Moodle LMS</strong></p>
          <p>Version 1.0.0</p>
          <p>© 2025 Mini Moodle. All rights reserved.</p>
        </CardContent>
      </Card>
    </motion.div>
  );
}
