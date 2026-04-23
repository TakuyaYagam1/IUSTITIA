import type { ArchiveEntry, ArchivePatchRequest } from '@shared/api';
import { api, unwrap } from '@shared/api';

export const archiveApi = {
  list: (): Promise<ArchiveEntry[]> => unwrap(api.GET('/api/archive', { params: { query: {} } })),

  get: (id: string): Promise<ArchiveEntry> =>
    unwrap(api.GET('/api/archive/{id}', { params: { path: { id } } })),

  patch: (id: string, patch: ArchivePatchRequest): Promise<ArchiveEntry> =>
    unwrap(
      api.PATCH('/api/archive/{id}', {
        params: { path: { id } },
        body: patch,
      }),
    ),
};
