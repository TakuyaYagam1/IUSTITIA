import { useContext } from 'react';
import { ToastContext, type ToastContextValue } from './context';

export const useToast = (): ToastContextValue => {
  const ctx = useContext(ToastContext);
  if (!ctx) {
    throw new Error('useToast must be used inside <ToastProvider>');
  }
  return ctx;
};
