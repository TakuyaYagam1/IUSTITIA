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

export const HOME_BY_ROLE: Record<Role, string> = {
  citizen: ROUTES.caseNew,
  prosecutor: ROUTES.cases,
  judge: ROUTES.hearings,
  registrar: ROUTES.registry,
};

export const DEFAULT_HOME = ROUTES.cases;

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
