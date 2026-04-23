import { api, unwrap } from '@shared/api';
import type { components } from '@shared/api/schema';

export type DocumentGenerateRequest = components['schemas']['DocumentGenerateRequest'];
export type Document = components['schemas']['Document'];

export const docgenApi = {
  generate: (payload: DocumentGenerateRequest): Promise<Document> =>
    unwrap(api.POST('/api/documents/generate', { body: payload })),

  getById: (id: string): Promise<Document> =>
    unwrap(api.GET('/api/documents/{id}', { params: { path: { id } } })),
};
