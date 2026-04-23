import { ToastProvider } from '@shared/ui';
import { AppProviders } from './providers';
import { AppRouter } from './router';

export const App = (): JSX.Element => (
  <AppProviders>
    <ToastProvider>
      <AppRouter />
    </ToastProvider>
  </AppProviders>
);
