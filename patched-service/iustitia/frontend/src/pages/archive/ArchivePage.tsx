import {
  archiveApi,
  loadClassifiedConfig,
  saveClassifiedConfig,
  type ClassifiedConfig,
} from '@entities/archive';
import type { ArchiveEntry, Verdict } from '@shared/api';
import { setClassifiedConfig } from '@shared/api';
import { Seal } from '@shared/ui';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useMemo, useState } from 'react';
import { useLocation } from 'react-router-dom';
import styles from './ArchivePage.module.css';

const VERDICT_FILTERS: { value: ClassifiedConfig['verdict']; label: string }[] = [
  { value: 'all', label: 'Все' },
  { value: 'guilty', label: 'Виновны' },
  { value: 'acquitted', label: 'Оправданы' },
];

const verdictText = (v: Verdict): { text: string; cls: string } => {
  if (v === 'guilty') return { text: 'Виновен', cls: styles.verdictGuilty ?? '' };
  if (v === 'acquitted') return { text: 'Оправдан', cls: styles.verdictAcquitted ?? '' };
  return { text: 'На рассмотрении', cls: '' };
};

export const ArchivePage = (): JSX.Element => {
  const location = useLocation();
  const [config, setConfig] = useState<ClassifiedConfig>(() =>
    loadClassifiedConfig(location.search),
  );

  useEffect(() => {
    setClassifiedConfig(config);
  }, [config]);

  useEffect(() => {
    const next = loadClassifiedConfig(location.search);
    setConfig(next);
    saveClassifiedConfig(next);
  }, [location.search]);

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['archive', 'list', config.verdict, config.sortBy, config.order],
    queryFn: () => archiveApi.list(),
  });

  const filtered = useMemo<ArchiveEntry[]>(() => {
    if (!data) return [];
    const base =
      config.verdict === 'all' ? data : data.filter((e) => e.final_verdict === config.verdict);
    const sorted = [...base].sort((a, b) => {
      if (config.sortBy === 'defendant') return a.defendant.localeCompare(b.defendant);
      return a.archived_at.localeCompare(b.archived_at);
    });
    return config.order === 'asc' ? sorted : sorted.reverse();
  }, [data, config.verdict, config.sortBy, config.order]);

  const onFilter = (v: ClassifiedConfig['verdict']) => {
    setConfig((prev) => {
      const next = { ...prev, verdict: v };
      saveClassifiedConfig(next);
      return next;
    });
  };

  return (
    <div className={styles.wrap}>
      <div className={styles.head}>
        <div className={styles.titleBlock}>
          <h1 className={styles.title}>Архив</h1>
          <span className={styles.subtitle}>Бесконечная картотека трибунала</span>
        </div>
        <div className={styles.filters}>
          {VERDICT_FILTERS.map((f) => (
            <button
              key={f.value}
              type="button"
              className={`${styles.filterBtn} ${config.verdict === f.value ? styles.filterBtnActive : ''}`}
              onClick={() => onFilter(f.value)}
            >
              {f.label}
            </button>
          ))}
        </div>
      </div>

      {isLoading && <div className={styles.message}>Получение архивных дел…</div>}
      {isError && (
        <div className={`${styles.message} ${styles.error}`}>
          Архив недоступен: {error instanceof Error ? error.message : 'неизвестно'}
        </div>
      )}
      {!isLoading && !isError && filtered.length === 0 && (
        <div className={styles.message}>Нет записей по выбранному фильтру.</div>
      )}

      <div className={styles.grid}>
        {filtered.map((e) => {
          const v = verdictText(e.final_verdict);
          return (
            <div key={e.id} className={styles.card}>
              <div className={styles.cardHead}>
                <span className={styles.cardNumber}>{e.id.slice(0, 8).toUpperCase()}</span>
                {e.classified_note && <Seal variant="topsecret">СЕКРЕТНО</Seal>}
              </div>
              <h3 className={styles.defendant}>{e.defendant}</h3>
              <span className={`${styles.verdict} ${v.cls}`}>{v.text}</span>
              {e.sentence && <p className={styles.sentence}>{e.sentence}</p>}
              <div className={styles.cardFoot}>
                <span>Архивировано</span>
                <span>{new Date(e.archived_at).toLocaleDateString('ru-RU')}</span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
