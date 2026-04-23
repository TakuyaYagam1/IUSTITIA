import type { ReactNode } from 'react';
import { ApiBootstrap } from './ApiBootstrap';
import { QueryClientProvider } from './QueryClientProvider';

interface AppProvidersProps {
  children: ReactNode;
}

export const AppProviders = ({ children }: AppProvidersProps): JSX.Element => (
  <QueryClientProvider>
    <ApiBootstrap>{children}</ApiBootstrap>
  </QueryClientProvider>
);
