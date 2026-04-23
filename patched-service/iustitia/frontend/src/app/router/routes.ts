import type { Role } from '@shared/api';

export const ROUTES = {
  login: '/login',
  cases: '/cases',
  caseNew: '/cases/new',
  caseView: (id = ':id') => `/cases/${id}`,
  caseOpinion: (id = ':id') => `/cases/${id}/opinion`,
  caseVerdict: (id = ':id') => `/cases/${id}/verdict`,
  hearings: '/hearings',
  registry: '/registry',
  archive: '/archive',
  complaintNew: '/complaints/new',
} as const;

// HOME_BY_ROLE задаёт куда редиректить после логина и при role-mismatch в
// RoleGuard. Раньше был единый POST_LOGIN_ROUTE = /cases, из-за чего citizen
// после логина попадал на /cases, фронт слал GET /api/cases, а бек отвечал
// 403 (backend ACL: /api/cases требует Judge|Prosecutor). Теперь каждый
// юзер идёт на ту страницу, которая ему реально доступна по backend-ACL.
export const HOME_BY_ROLE: Record<Role, string> = {
  citizen: ROUTES.caseNew,
  prosecutor: ROUTES.cases,
  judge: ROUTES.hearings,
  registrar: ROUTES.registry,
};

// DEFAULT_HOME - fallback для catch-all маршрута и неизвестных ролей;
// выбран /cases как наиболее общая страница авторизованной зоны - guard
// дальше перекинет куда надо по HOME_BY_ROLE, если роль известна.
export const DEFAULT_HOME = ROUTES.cases;

// ROUTE_ROLES выровнен с backend RBAC в
// patched-services/iustitia/backend/internal/controller/restapi/v1/router.go:
//   GET  /api/cases              -> RequireRole(Judge, Prosecutor)
//   GET  /api/cases/{id}         -> RequireRole(Judge, Prosecutor)
//   POST /api/complaints         -> RequireRole(Citizen)
//   GET  /api/documents/{id}     -> RequireRole(Prosecutor, Judge)
//   POST /api/documents/generate -> RequireRole(Judge)
//   GET  /api/archive            -> авторизация без роли
//   PATCH /api/archive/{id}      -> RequireRole(Judge)
// docgen на фронте - только Judge (прокурор не генерирует приговоры вручную).
export const ROUTE_ROLES: Record<string, Role[]> = {
  [ROUTES.cases]: ['prosecutor', 'judge'],
  [ROUTES.caseNew]: ['citizen'],
  [ROUTES.complaintNew]: ['citizen'],
  [ROUTES.caseOpinion()]: ['prosecutor'],
  [ROUTES.caseVerdict()]: ['judge'],
  [ROUTES.hearings]: ['judge'],
  [ROUTES.registry]: ['registrar'],
  [ROUTES.archive]: ['judge'],
};
