import { QueryClientProvider as TanstackQueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { queryClient } from './queryClient';

interface QueryClientProviderProps {
  children: ReactNode;
}

export const QueryClientProvider = ({ children }: QueryClientProviderProps): JSX.Element => (
  <TanstackQueryClientProvider client={queryClient}>{children}</TanstackQueryClientProvider>
);
