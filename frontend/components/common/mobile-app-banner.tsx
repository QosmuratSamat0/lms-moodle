"use client";

import { useState, useEffect } from "react";
import { X, Download, Smartphone } from "lucide-react";
import { Button } from "@/components/ui/button";

export function MobileAppBanner() {
  const [isVisible, setIsVisible] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    // Check if mobile device
    const checkMobile = () => {
      const userAgent = navigator.userAgent || navigator.vendor;
      const isMobileDevice = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(
        userAgent.toLowerCase()
      );
      // Also check screen width
      const isSmallScreen = window.innerWidth < 768;
      return isMobileDevice || isSmallScreen;
    };

    // Check if banner was dismissed
    const dismissed = localStorage.getItem("app-banner-dismissed");
    const dismissedTime = dismissed ? parseInt(dismissed) : 0;
    const oneDayAgo = Date.now() - 24 * 60 * 60 * 1000;

    setIsMobile(checkMobile());
    // Show banner if mobile and not dismissed in last 24 hours
    setIsVisible(checkMobile() && dismissedTime < oneDayAgo);
  }, []);

  const handleDismiss = () => {
    setIsVisible(false);
    localStorage.setItem("app-banner-dismissed", Date.now().toString());
  };

  const handleInstall = () => {
    // Link to APK download - update this URL to your actual APK location
    window.open("/downloads/ednova.apk", "_blank");
  };

  if (!isMobile || !isVisible) {
    return null;
  }

  return (
    <div className="fixed bottom-0 left-0 right-0 z-50 safe-area-bottom">
      <div className="bg-primary text-primary-foreground p-4 shadow-lg">
        <div className="flex items-center gap-3">
          <div className="flex-shrink-0 h-12 w-12 bg-primary-foreground/20 rounded-xl flex items-center justify-center">
            <Smartphone className="h-6 w-6" />
          </div>
          <div className="flex-1 min-w-0">
            <h4 className="font-semibold text-sm">Get EDnova App</h4>
            <p className="text-xs text-primary-foreground/80 truncate">
              Better experience with our mobile app
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Button
              size="sm"
              variant="secondary"
              className="h-9 px-3 text-xs font-medium"
              onClick={handleInstall}
            >
              <Download className="h-4 w-4 mr-1" />
              Install
            </Button>
            <Button
              size="icon"
              variant="ghost"
              className="h-8 w-8 text-primary-foreground/70 hover:text-primary-foreground hover:bg-primary-foreground/10"
              onClick={handleDismiss}
            >
              <X className="h-4 w-4" />
              <span className="sr-only">Dismiss</span>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
