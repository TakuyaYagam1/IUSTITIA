import { caseApi } from '@entities/case';
import type { CaseOpinion, OpinionCreateRequest } from '@shared/api';
import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';

export const useOpinionFile = (
  caseId: string,
): UseMutationResult<CaseOpinion, Error, OpinionCreateRequest> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body) => caseApi.fileOpinion(caseId, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cases'] });
      queryClient.invalidateQueries({ queryKey: ['cases', 'byId', caseId] });
      queryClient.invalidateQueries({ queryKey: ['cases', 'opinion', caseId] });
      queryClient.invalidateQueries({ queryKey: ['hearings'] });
    },
  });
};
