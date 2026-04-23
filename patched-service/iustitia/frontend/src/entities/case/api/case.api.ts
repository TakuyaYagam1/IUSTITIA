import {
  api,
  unwrap,
  type Case,
  type CaseAcceptRequest,
  type CaseCreateRequest,
  type CaseDismissRequest,
  type CaseDocument,
  type CaseOpinion,
  type Complaint,
  type HearingItem,
  type OpinionCreateRequest,
  type VerdictRequest,
  type VerdictResult,
} from '@shared/api';

type SearchBody = {
  q?: string;
  order_by?: string;
  direction?: 'ASC' | 'DESC';
  limit?: number;
  offset?: number;
};

// /api/cases           -> RequireRole(Judge, Prosecutor)
// /api/cases/{id}      -> RequireRole(Judge, Prosecutor)
// /api/cases/search    -> auth only, все роли (единственный роут картотеки,
//                         доступный registrar - см. router.go).
export const caseApi = {
  list: (params?: { limit?: number; offset?: number }): Promise<Case[]> =>
    unwrap(
      api.GET('/api/cases', {
        params: { query: params ?? {} },
      }),
    ),

  get: (id: string): Promise<Case> =>
    unwrap(
      api.GET('/api/cases/{id}', {
        params: { path: { id } },
      }),
    ),

  search: (body?: SearchBody): Promise<Case[]> =>
    unwrap(
      api.POST('/api/cases/search', {
        body: {
          q: body?.q ?? '',
          order_by: body?.order_by ?? 'created_at',
          direction: body?.direction ?? 'DESC',
          limit: body?.limit ?? 50,
          offset: body?.offset ?? 0,
        },
      }),
    ),

  complaints: (caseId: string): Promise<Complaint[]> =>
    unwrap(
      api.GET('/api/complaints/{case_id}', {
        params: { path: { case_id: caseId } },
      }),
    ),

  // Trial workflow endpoints
  create: (body: CaseCreateRequest): Promise<Case> => unwrap(api.POST('/api/cases', { body })),

  accept: (id: string, body: CaseAcceptRequest): Promise<Case> =>
    unwrap(
      api.POST('/api/cases/{id}/accept', {
        params: { path: { id } },
        body,
      }),
    ),

  dismiss: (id: string, body: CaseDismissRequest) =>
    unwrap(
      api.POST('/api/cases/{id}/dismiss', {
        params: { path: { id } },
        body,
      }),
    ),

  fileOpinion: (id: string, body: OpinionCreateRequest): Promise<CaseOpinion> =>
    unwrap(
      api.POST('/api/cases/{id}/opinion', {
        params: { path: { id } },
        body,
      }),
    ),

  getOpinion: (id: string): Promise<CaseOpinion> =>
    unwrap(
      api.GET('/api/cases/{id}/opinion', {
        params: { path: { id } },
      }),
    ),

  finalize: (id: string, body: VerdictRequest): Promise<VerdictResult> =>
    unwrap(
      api.POST('/api/cases/{id}/verdict', {
        params: { path: { id } },
        body,
      }),
    ),

  listHearings: (): Promise<HearingItem[]> => unwrap(api.GET('/api/hearings')),

  listDocuments: (id: string): Promise<CaseDocument[]> =>
    unwrap(
      api.GET('/api/cases/{id}/documents', {
        params: { path: { id } },
      }),
    ),
};
