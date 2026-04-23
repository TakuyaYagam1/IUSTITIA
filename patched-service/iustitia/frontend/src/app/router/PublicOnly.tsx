import { useAuthStore } from '@entities/user';
import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { DEFAULT_HOME, HOME_BY_ROLE } from './routes';

// Если юзер уже залогинен и лезет на /login, PublicOnly его редиректит на
// его личный landing (по роли). Раньше редиректило на единый
// POST_LOGIN_ROUTE = /cases, из-за чего citizen попадал на чужую страницу
// и ловил 403 от GET /api/cases. Теперь через HOME_BY_ROLE.
export const PublicOnly = ({ children }: { children: ReactNode }) => {
  const session = useAuthStore((s) => s.session);
  if (session) {
    const landing = HOME_BY_ROLE[session.role] ?? DEFAULT_HOME;
    return <Navigate to={landing} replace />;
  }
  return <>{children}</>;
};
