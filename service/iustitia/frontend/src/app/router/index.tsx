import { ArchivePage } from '@pages/archive';
import { CaseNewPage } from '@pages/case-new';
import { CaseOpinionPage } from '@pages/case-opinion';
import { CaseVerdictPage } from '@pages/case-verdict';
import { CaseViewPage } from '@pages/case-view';
import { CasesPage } from '@pages/cases';
import { ComplaintPage } from '@pages/complaint';
import { HearingsPage } from '@pages/hearings';
import { LoginPage } from '@pages/login';
import { RegistryPage } from '@pages/registry';
import { AppLayout } from '@widgets/layout';
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';
import { PublicOnly } from './PublicOnly';
import { RoleGuard } from './RoleGuard';
import { DEFAULT_HOME, ROUTE_ROLES, ROUTES } from './routes';

const router = createBrowserRouter([
  {
    path: ROUTES.login,
    element: (
      <PublicOnly>
        <LoginPage />
      </PublicOnly>
    ),
  },
  {
    element: (
      <RoleGuard>
        <AppLayout />
      </RoleGuard>
    ),
    children: [
      {
        path: ROUTES.cases,
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.cases]}>
            <CasesPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.caseNew,
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.caseNew]}>
            <CaseNewPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.complaintNew,
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.complaintNew]}>
            <ComplaintPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.caseOpinion(),
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.caseOpinion()]}>
            <CaseOpinionPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.caseVerdict(),
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.caseVerdict()]}>
            <CaseVerdictPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.caseView(),
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.cases]}>
            <CaseViewPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.hearings,
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.hearings]}>
            <HearingsPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.registry,
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.registry]}>
            <RegistryPage />
          </RoleGuard>
        ),
      },
      {
        path: ROUTES.archive,
        element: (
          <RoleGuard roles={ROUTE_ROLES[ROUTES.archive]}>
            <ArchivePage />
          </RoleGuard>
        ),
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to={DEFAULT_HOME} replace />,
  },
]);

export const AppRouter = (): JSX.Element => <RouterProvider router={router} />;
