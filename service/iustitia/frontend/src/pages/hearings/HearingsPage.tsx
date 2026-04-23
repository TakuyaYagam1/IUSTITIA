import { ROUTES } from '@/app/router/routes';
import { useHearings } from '@features/case-view';
import type { PreliminaryVerdict } from '@shared/api';
import { Link } from 'react-router-dom';
import styles from './HearingsPage.module.css';

const VERDICT_LABEL: Record<PreliminaryVerdict, string> = {
  guilty: 'Прокурор: виновен',
  acquitted: 'Прокурор: невиновен',
  dismissed: 'Прокурор: прекратить',
};

export const HearingsPage = (): JSX.Element => {
  const { data, isLoading, isError, error } = useHearings();

  return (
    <div className={styles.wrap}>
      <div className={styles.header}>
        <h1 className={styles.title}>Слушания</h1>
        <span className={styles.subtitle}>Очередь дел, готовых к приговору</span>
      </div>

      {isLoading && <div className={styles.message}>Получение очереди…</div>}
      {isError && (
        <div className={`${styles.message} ${styles.error}`}>
          Ошибка: {error instanceof Error ? error.message : 'неизвестно'}
        </div>
      )}
      {!isLoading && !isError && (data ?? []).length === 0 && (
        <div className={styles.message}>Дел в очереди нет.</div>
      )}

      <div className={styles.grid}>
        {(data ?? []).map((h) => (
          <Link key={h.case.id} to={ROUTES.caseVerdict(h.case.id)} className={styles.card}>
            <span className={styles.cardNumber}>Дело №{h.case.seq_num}/Δ</span>
            <h3 className={styles.defendant}>{h.case.defendant}</h3>
            <span className={styles.opinion}>{VERDICT_LABEL[h.opinion.preliminary_verdict]}</span>
            <p className={styles.reasoning}>{h.opinion.reasoning}</p>
          </Link>
        ))}
      </div>
    </div>
  );
};
