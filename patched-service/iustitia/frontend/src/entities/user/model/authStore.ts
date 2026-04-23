import type { Role } from '@shared/api';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthSession {
  token: string;
  userId: string;
  role: Role;
  dome: string;
  username: string;
}

interface AuthState {
  session: AuthSession | null;
  setSession: (session: AuthSession) => void;
  clearSession: () => void;
  isAuthenticated: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      session: null,
      setSession: (session) => set({ session }),
      clearSession: () => set({ session: null }),
      isAuthenticated: () => get().session !== null,
    }),
    {
      name: 'iustitia.auth',
      partialize: (state) => ({ session: state.session }),
    },
  ),
);

export const getAuthToken = (): string | null => useAuthStore.getState().session?.token ?? null;
