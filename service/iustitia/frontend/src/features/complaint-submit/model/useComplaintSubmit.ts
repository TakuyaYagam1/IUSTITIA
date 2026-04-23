import type { Complaint } from '@shared/api';
import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';
import { complaintApi, type ComplaintCreatePayload } from '../api/complaint.api';

export const useComplaintSubmit = (): UseMutationResult<
  Complaint,
  Error,
  ComplaintCreatePayload
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: complaintApi.create,
    onSuccess: (complaint) => {
      queryClient.invalidateQueries({ queryKey: ['cases', 'complaints', complaint.case_id] });
      queryClient.invalidateQueries({ queryKey: ['cases', 'list'] });
    },
  });
};
