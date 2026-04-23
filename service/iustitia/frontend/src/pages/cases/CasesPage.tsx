import { ROUTES } from '@/app/router/routes';
import { useAuthStore } from '@entities/user';
import { useCaseFilterStore, useCaseList, type StatusFilter } from '@features/case-view';
import type { Case } from '@shared/api';
import { Seal } from '@shared/ui';
import { useMemo } from 'react';
import { Link } from 'react-router-dom';
import styles from './CasesPage.module.css';

const STATUS_FILTERS: { value: StatusFilter; label: string }[] = [
  { value: 'all', label: 'Все' },
  { value: 'open', label: 'Открытые' },
  { value: 'hearing', label: 'Слушание' },
  { value: 'closed', label: 'Закрытые' },
];

const statusText = (c: Case): { text: string; cls: string } => {
  if (c.verdict === 'guilty') return { text: 'Виновен', cls: styles.statusGuilty ?? '' };
  if (c.verdict === 'acquitted') return { text: 'Оправдан', cls: styles.statusAcquit ?? '' };
  if (c.status === 'hearing') return { text: 'Слушание', cls: styles.statusPending ?? '' };
  if (c.status === 'closed') return { text: 'Закрыто', cls: '' };
  return { text: 'Открыто', cls: styles.statusPending ?? '' };
};

export const CasesPage = (): JSX.Element => {
  const { data, isLoading, isError, error } = useCaseList();
  const status = useCaseFilterStore((s) => s.status);
  const setStatus = useCaseFilterStore((s) => s.setStatus);
  const role = useAuthStore((s) => s.session?.role);
  const isCitizen = role === 'citizen';

  const filtered = useMemo(() => {
    if (!data) return [];
    if (status === 'all') return data;
    return data.filter((c) => c.status === status);
  }, [data, status]);

  return (
    <div className={styles.wrap}>
      <div className={styles.header}>
        <div className={styles.titleBlock}>
          <h1 className={styles.title}>Картотека дел</h1>
          <span className={styles.subtitle}>Трибунал Одиннадцатого Государства</span>
        </div>
        <div className={styles.filters}>
          <div className={styles.filterGroup}>
            {STATUS_FILTERS.map((f) => (
              <button
                key={f.value}
                type="button"
                className={`${styles.filterBtn} ${status === f.value ? styles.filterBtnActive : ''}`}
                onClick={() => setStatus(f.value)}
              >
                {f.label}
              </button>
            ))}
          </div>
          {isCitizen && (
            <Link to={ROUTES.caseNew} className={styles.ctaLink}>
              Новое дело
            </Link>
          )}
        </div>
      </div>

      {isLoading && <div className={styles.message}>Получение картотеки…</div>}
      {isError && (
        <div className={`${styles.message} ${styles.error}`}>
          Ошибка связи с картотекой: {error instanceof Error ? error.message : 'неизвестно'}
        </div>
      )}
      {!isLoading && !isError && filtered.length === 0 && (
        <div className={styles.message}>Нет дел по выбранному фильтру.</div>
      )}

      <div className={styles.grid}>
        {filtered.map((c) => {
          const st = statusText(c);
          return (
            <Link key={c.id} to={ROUTES.caseView(c.id)} className={styles.card}>
              <div className={styles.cardHead}>
                <span className={styles.cardNumber}>Дело №{c.seq_num}/Δ</span>
                {c.classified_note && <Seal variant="topsecret">СЕКРЕТНО</Seal>}
              </div>
              <h3 className={styles.defendant}>{c.defendant}</h3>
              <p className={styles.crime}>{c.crime}</p>
              <div className={styles.cardFoot}>
                <span className={`${styles.status} ${st.cls}`}>{st.text}</span>
                <span>{new Date(c.created_at).toLocaleDateString('ru-RU')}</span>
              </div>
            </Link>
          );
        })}
      </div>
    </div>
  );
};
