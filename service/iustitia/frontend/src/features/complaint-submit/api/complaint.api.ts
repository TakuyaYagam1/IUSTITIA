import { api, unwrap, type Complaint } from '@shared/api';

export interface ComplaintCreatePayload {
  case_id: string;
  text: string;
}

export const complaintApi = {
  create: (payload: ComplaintCreatePayload): Promise<Complaint> =>
    unwrap(api.POST('/api/complaints', { body: payload })),
};
