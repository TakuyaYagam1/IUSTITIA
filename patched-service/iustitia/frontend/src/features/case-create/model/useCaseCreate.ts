import { caseApi } from '@entities/case';
import type { Case, CaseCreateRequest } from '@shared/api';
import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';

export const useCaseCreate = (): UseMutationResult<Case, Error, CaseCreateRequest> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: caseApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cases', 'list'] });
      queryClient.invalidateQueries({ queryKey: ['cases', 'search'] });
    },
  });
};
