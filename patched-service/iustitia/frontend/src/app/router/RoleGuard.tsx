import { useAuthStore } from '@entities/user';
import type { Role } from '@shared/api';
import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { DEFAULT_HOME, HOME_BY_ROLE, ROUTES } from './routes';

interface RoleGuardProps {
  roles?: Role[];
  children: ReactNode;
}

export const RoleGuard = ({ roles, children }: RoleGuardProps) => {
  const session = useAuthStore((s) => s.session);
  const location = useLocation();

  if (!session) {
    const next = `${location.pathname}${location.search}`;
    const target = `${ROUTES.login}?next=${encodeURIComponent(next)}`;
    return <Navigate to={target} state={{ from: location.pathname }} replace />;
  }

  if (roles && roles.length > 0 && !roles.includes(session.role)) {
    // При несоответствии роли редиректим на ЛИЧНЫЙ home (HOME_BY_ROLE),
    // а не на общий /cases. Иначе citizen, ткнувший руками /docgen или
    // /archive в адресной строке, снова попал бы на /cases и получил 403
    // от GET /api/cases (backend его туда не пускает).
    const landing = HOME_BY_ROLE[session.role] ?? DEFAULT_HOME;
    return <Navigate to={landing} replace />;
  }

  return <>{children}</>;
};
