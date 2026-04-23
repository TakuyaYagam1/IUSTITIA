import { caseApi } from '@entities/case';
import type { Case, CaseAcceptRequest, CaseDismissRequest, VerdictResult } from '@shared/api';
import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';

export const useCaseAccept = (
  caseId: string,
): UseMutationResult<Case, Error, CaseAcceptRequest> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body) => caseApi.accept(caseId, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cases'] });
      queryClient.invalidateQueries({ queryKey: ['cases', 'byId', caseId] });
    },
  });
};

export const useCaseDismiss = (
  caseId: string,
): UseMutationResult<VerdictResult, Error, CaseDismissRequest> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body) => caseApi.dismiss(caseId, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cases'] });
      queryClient.invalidateQueries({ queryKey: ['archive'] });
    },
  });
};

export const useProsecutors = () => {
  const queryClient = useQueryClient();
  return { queryClient };
};
