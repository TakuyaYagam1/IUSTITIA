import { useAuthStore } from '@entities/user';
import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { DEFAULT_HOME, HOME_BY_ROLE } from './routes';

export const PublicOnly = ({ children }: { children: ReactNode }) => {
  const session = useAuthStore((s) => s.session);
  if (session) {
    const landing = HOME_BY_ROLE[session.role] ?? DEFAULT_HOME;
    return <Navigate to={landing} replace />;
  }
  return <>{children}</>;
};
