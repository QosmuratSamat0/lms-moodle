import { create } from "zustand";
import { persist } from "zustand/middleware";

interface UIState {
  // Sidebar
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  toggleSidebarCollapsed: () => void;

  // Modals
  loginModalOpen: boolean;
  signupModalOpen: boolean;
  setLoginModalOpen: (open: boolean) => void;
  setSignupModalOpen: (open: boolean) => void;
  openLoginModal: () => void;
  openSignupModal: () => void;
  closeAuthModals: () => void;

  // Command palette
  commandPaletteOpen: boolean;
  setCommandPaletteOpen: (open: boolean) => void;

  // Mobile
  isMobile: boolean;
  setIsMobile: (isMobile: boolean) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      // Sidebar
      sidebarOpen: true,
      sidebarCollapsed: false,
      toggleSidebar: () =>
        set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
      toggleSidebarCollapsed: () =>
        set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),

      // Modals
      loginModalOpen: false,
      signupModalOpen: false,
      setLoginModalOpen: (open) =>
        set({ loginModalOpen: open, signupModalOpen: false }),
      setSignupModalOpen: (open) =>
        set({ signupModalOpen: open, loginModalOpen: false }),
      openLoginModal: () =>
        set({ loginModalOpen: true, signupModalOpen: false }),
      openSignupModal: () =>
        set({ signupModalOpen: true, loginModalOpen: false }),
      closeAuthModals: () =>
        set({ loginModalOpen: false, signupModalOpen: false }),

      // Command palette
      commandPaletteOpen: false,
      setCommandPaletteOpen: (open) => set({ commandPaletteOpen: open }),

      // Mobile
      isMobile: false,
      setIsMobile: (isMobile) => set({ isMobile }),
    }),
    {
      name: "ui-storage",
      partialize: (state) => ({
        sidebarCollapsed: state.sidebarCollapsed,
      }),
    },
  ),
);
