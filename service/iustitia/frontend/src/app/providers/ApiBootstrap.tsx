import { useAuthStore } from '@entities/user';
import { setOnUnauthorized, setTokenGetter } from '@shared/api';
import { type ReactNode } from 'react';

interface ApiBootstrapProps {
  children: ReactNode;
}

setTokenGetter(() => useAuthStore.getState().session?.token ?? null);
setOnUnauthorized(() => {
  useAuthStore.getState().clearSession();
});

export const ApiBootstrap = ({ children }: ApiBootstrapProps): JSX.Element => {
  return <>{children}</>;
};
