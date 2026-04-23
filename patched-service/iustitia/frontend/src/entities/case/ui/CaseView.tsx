// Patch F1 (Stored XSS render).
// Raw-версия рендерила classified_note и complaint.text через
// dangerouslySetInnerHTML напрямую, поэтому любая <script>/<iframe>-нагрузка
// из донесений (см. POST /api/complaints) выполнялась в контексте фронта
// при первом просмотре дела. Серверная санитизация на create (bluemonday)
// закрывает только storage-путь; на случай легаси-данных и прочих источников
// добавляем client-side sanitize через DOMPurify с явным whitelist-ом тегов
// и атрибутов. Теги оставлены минимальные - бумажно-документальная стилистика
// (b/i/em/u/code/br/p/span) + class для CSS-модулей.
import type { Case, CaseOpinion, Complaint, PreliminaryVerdict } from '@shared/api';
import { Paper, Seal } from '@shared/ui';
import DOMPurify from 'dompurify';
import { type ReactNode } from 'react';
import styles from './CaseView.module.css';

const RICH_TEXT_CONFIG = {
  ALLOWED_TAGS: ['b', 'i', 'em', 'strong', 'u', 'code', 'br', 'p', 'span'],
  ALLOWED_ATTR: ['class'],
};

const sanitize = (html: string): string => DOMPurify.sanitize(html, RICH_TEXT_CONFIG);

interface CaseViewProps {
  caseItem: Case;
  complaints?: Complaint[];
  opinion?: CaseOpinion | null;
  actions?: ReactNode;
}

const PRELIM_LABEL: Record<PreliminaryVerdict, string> = {
  guilty: 'Виновен',
  acquitted: 'Невиновен',
  dismissed: 'Прекратить',
};

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

export const CaseView = ({
  caseItem,
  complaints,
  opinion,
  actions,
}: CaseViewProps): JSX.Element => {
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
              dangerouslySetInnerHTML={{ __html: sanitize(caseItem.classified_note) }}
            />
          </div>
        )}

        {opinion && (
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>
              Заключение прокурора: {PRELIM_LABEL[opinion.preliminary_verdict]}
            </h2>
            <p className={styles.crime}>{opinion.reasoning}</p>
          </div>
        )}

        {actions && <div className={styles.section}>{actions}</div>}

        {complaints && complaints.length > 0 && (
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>Заявления в деле ({complaints.length})</h2>
            <div className={styles.complaints}>
              {complaints.map((c) => (
                <div key={c.id} className={styles.complaint}>
                  <div className={styles.complaintMeta}>
                    {formatDate(c.created_at)} · автор {c.author_id.slice(0, 8)}
                  </div>
                  <div
                    className={styles.complaintText}
                    dangerouslySetInnerHTML={{ __html: sanitize(c.text) }}
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
