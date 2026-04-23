import { createContext } from 'react';

export type ToastKind = 'info' | 'success' | 'error';

export interface ToastItem {
  id: string;
  kind: ToastKind;
  title: string;
  message: string;
}

export interface ToastContextValue {
  push: (toast: Omit<ToastItem, 'id'>) => void;
}

export const ToastContext = createContext<ToastContextValue | null>(null);
