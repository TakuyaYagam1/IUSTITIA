import type { Case, Complaint } from '@shared/api';
import { Paper, Seal } from '@shared/ui';
import styles from './CaseView.module.css';

interface CaseViewProps {
  caseItem: Case;
  complaints?: Complaint[];
}

const verdictLabel = (c: Case): { text: string; cls: string } => {
  if (c.verdict === 'guilty') return { text: 'Виновен', cls: styles.statusGuilty ?? '' };
  if (c.verdict === 'acquitted') return { text: 'Оправдан', cls: styles.statusAcquit ?? '' };
  if (c.status === 'hearing') return { text: 'Слушание', cls: styles.statusPending ?? '' };
  if (c.status === 'closed') return { text: 'Закрыто', cls: '' };
  return { text: 'Открыто', cls: styles.statusPending ?? '' };
};

const formatDate = (iso: string) =>
  new Date(iso).toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

export const CaseView = ({ caseItem, complaints }: CaseViewProps): JSX.Element => {
  const verdict = verdictLabel(caseItem);
  const hasClassified = Boolean(caseItem.classified_note);

  return (
    <Paper variant="paper">
      <div className={styles.root}>
        <div className={styles.header}>
          <div className={styles.heading}>
            <span className={styles.caseNumber}>Дело №{caseItem.seq_num}/Δ</span>
            <h1 className={styles.defendant}>{caseItem.defendant}</h1>
          </div>
          <div className={styles.seals}>
            {hasClassified && <Seal variant="topsecret">СЕКРЕТНОЕ ПРИЛОЖЕНИЕ №0</Seal>}
            {caseItem.verdict === 'guilty' && <Seal variant="rejected" />}
            {caseItem.verdict === 'acquitted' && <Seal variant="approved" />}
          </div>
        </div>

        <div className={styles.meta}>
          <div className={styles.metaItem}>
            <span className={styles.metaLabel}>Статус</span>
            <span className={`${styles.metaValue} ${verdict.cls}`}>{verdict.text}</span>
          </div>
          <div className={styles.metaItem}>
            <span className={styles.metaLabel}>Зарегистрировано</span>
            <span className={styles.metaValue}>{formatDate(caseItem.created_at)}</span>
          </div>
          <div className={styles.metaItem}>
            <span className={styles.metaLabel}>ID картотеки</span>
            <span
              className={styles.metaValue}
              style={{ fontFamily: 'var(--font-mono)', fontSize: 'var(--fs-sm)' }}
            >
              {caseItem.id}
            </span>
          </div>
        </div>

        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>Суть обвинения</h2>
          <p className={styles.crime}>{caseItem.crime}</p>
        </div>

        {hasClassified && caseItem.classified_note && (
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>Секретное Приложение №0</h2>
            <div
              className={styles.classified}
              dangerouslySetInnerHTML={{ __html: caseItem.classified_note }}
            />
          </div>
        )}

        {complaints && complaints.length > 0 && (
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>Доносы в деле ({complaints.length})</h2>
            <div className={styles.complaints}>
              {complaints.map((c) => (
                <div key={c.id} className={styles.complaint}>
                  <div className={styles.complaintMeta}>
                    {formatDate(c.created_at)} · автор {c.author_id.slice(0, 8)}
                  </div>
                  <div
                    className={styles.complaintText}
                    dangerouslySetInnerHTML={{ __html: c.text }}
                  />
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </Paper>
  );
};
