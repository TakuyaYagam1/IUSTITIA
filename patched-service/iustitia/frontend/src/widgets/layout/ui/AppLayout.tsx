import { Header } from '@widgets/header';
import { Outlet } from 'react-router-dom';
import styles from './AppLayout.module.css';

export const AppLayout = (): JSX.Element => (
  <div className={styles.layout}>
    <Header />
    <main className={styles.main}>
      <Outlet />
    </main>
  </div>
);
