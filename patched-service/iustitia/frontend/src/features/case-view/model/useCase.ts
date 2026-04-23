import { caseApi } from '@entities/case';
import { userApi } from '@entities/user';
import type {
    Case,
    CaseDocument,
    CaseOpinion,
    Complaint,
    HearingItem,
    Role,
    UserListItem,
} from '@shared/api';
import { useQuery, type UseQueryResult } from '@tanstack/react-query';

export const useCaseList = (): UseQueryResult<Case[]> =>
  useQuery({
    queryKey: ['cases', 'list'],
    queryFn: () => caseApi.list({ limit: 100, offset: 0 }),
  });

// useCaseSearch вместо useCaseList на страницах, которые должны быть
// доступны registrar: /api/cases заблокирован для него бекендом (см.
// patched router.go: RequireRole(Judge, Prosecutor)), а /api/cases/search
// открыт для всех авторизованных. Фильтрацию делаем клиентом через q="".
export const useCaseSearch = (): UseQueryResult<Case[]> =>
  useQuery({
    queryKey: ['cases', 'search', 'all'],
    queryFn: () => caseApi.search({ limit: 100 }),
    staleTime: 0,
    refetchOnWindowFocus: true,
  });

export const useCase = (id: string | undefined): UseQueryResult<Case> =>
  useQuery({
    queryKey: ['cases', 'detail', id],
    queryFn: () => caseApi.get(id as string),
    enabled: Boolean(id),
  });

export const useCaseComplaints = (id: string | undefined): UseQueryResult<Complaint[]> =>
  useQuery({
    queryKey: ['cases', 'complaints', id],
    queryFn: () => caseApi.complaints(id as string),
    enabled: Boolean(id),
  });

export const useCaseOpinion = (
  id: string | undefined,
  caseStatus?: string,
): UseQueryResult<CaseOpinion> =>
  useQuery({
    queryKey: ['cases', 'opinion', id],
    queryFn: () => caseApi.getOpinion(id as string),
    enabled: Boolean(id) && (caseStatus === 'hearing' || caseStatus === 'closed'),
    retry: false,
  });

export const useHearings = (): UseQueryResult<HearingItem[]> =>
  useQuery({
    queryKey: ['hearings'],
    queryFn: () => caseApi.listHearings(),
  });

export const useUsersByRole = (role: Role, enabled = true): UseQueryResult<UserListItem[]> =>
  useQuery({
    queryKey: ['users', 'byRole', role],
    queryFn: () => userApi.listByRole(role),
    enabled,
  });

export const useCaseDocuments = (id: string | undefined): UseQueryResult<CaseDocument[]> =>
  useQuery({
    queryKey: ['cases', 'documents', id],
    queryFn: () => caseApi.listDocuments(id as string),
    enabled: Boolean(id),
  });
