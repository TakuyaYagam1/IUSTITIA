import { ROUTES } from '@/app/router/routes';
import { useAuthStore, userApi } from '@entities/user';
import type { Role } from '@shared/api';
import { Button, Nameplate } from '@shared/ui';
import { LogOut } from 'lucide-react';
import { useCallback } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import styles from './Header.module.css';

interface NavItem {
  to: string;
  label: string;
  roles: Role[];
}

const NAV: NavItem[] = [
  { to: ROUTES.caseNew, label: 'Новое дело', roles: ['citizen'] },
  { to: ROUTES.complaintNew, label: 'Подать заявление', roles: ['citizen'] },
  { to: ROUTES.cases, label: 'Дела', roles: ['prosecutor'] },
  { to: ROUTES.hearings, label: 'Слушания', roles: ['judge'] },
  { to: ROUTES.registry, label: 'Канцелярия', roles: ['registrar'] },
  { to: ROUTES.archive, label: 'Архив', roles: ['judge'] },
];

const ROLE_LABELS: Record<Role, string> = {
  citizen: 'Гражданин',
  registrar: 'Регистратор',
  prosecutor: 'Прокурор',
  judge: 'Судья',
};

export const Header = () => {
  const session = useAuthStore((s) => s.session);
  const clearSession = useAuthStore((s) => s.clearSession);
  const navigate = useNavigate();

  const handleLogout = useCallback(async () => {
    try {
      await userApi.logout();
    } catch {
      /* noop */
    } finally {
      clearSession();
      navigate(ROUTES.login, { replace: true });
    }
  }, [clearSession, navigate]);

  if (!session) {
    return null;
  }

  const visibleNav = NAV.filter((item) => item.roles.includes(session.role));

  return (
    <header className={styles.header}>
      <div className={styles.left}>
        <div className={styles.title}>
          <span className={styles.titleMain}>IUSTITIA</span>
          <span className={styles.titleSub}>Одиннадцатое Государство · Свободный Марс</span>
        </div>
      </div>

      <nav className={styles.nav}>
        {visibleNav.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.to === ROUTES.cases}
            className={({ isActive }) =>
              [styles.navLink, isActive && styles.navLinkActive].filter(Boolean).join(' ')
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className={styles.right}>
        <Nameplate name={session.username} role={ROLE_LABELS[session.role]} />
        <Button variant="icon" onClick={handleLogout} aria-label="Выход">
          <LogOut size={16} />
        </Button>
      </div>
    </header>
  );
};
