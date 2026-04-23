import { caseApi } from '@entities/case';
import type { VerdictRequest, VerdictResult } from '@shared/api';
import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';

export const useVerdictFile = (
  caseId: string,
): UseMutationResult<VerdictResult, Error, VerdictRequest> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body) => caseApi.finalize(caseId, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cases'] });
      queryClient.invalidateQueries({ queryKey: ['cases', 'byId', caseId] });
      queryClient.invalidateQueries({ queryKey: ['hearings'] });
      queryClient.invalidateQueries({ queryKey: ['archive'] });
    },
  });
};
