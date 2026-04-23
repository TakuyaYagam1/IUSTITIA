import { ROUTES } from '@/app/router/routes';
import { useCase, useCaseDocuments, useCaseOpinion } from '@features/case-view';
import type { Document as DocgenDocument } from '@features/docgen';
import { Preview, useGenerateDocument } from '@features/docgen';
import { VerdictForm } from '@features/verdict-file';
import type { CaseDocument, PreliminaryVerdict } from '@shared/api';
import { Paper, useToast } from '@shared/ui';
import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import styles from './CaseVerdictPage.module.css';

const VERDICT_RU: Record<PreliminaryVerdict, string> = {
  guilty: 'виновным',
  acquitted: 'невиновным',
  dismissed: 'прекратить дело',
};

export const CaseVerdictPage = (): JSX.Element => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const toast = useToast();
  const { data: caseData, isLoading: caseLoading } = useCase(id);
  const { data: opinion, isLoading: opinionLoading } = useCaseOpinion(id, caseData?.status);
  const { data: docs, refetch: refetchDocs } = useCaseDocuments(id);
  const generate = useGenerateDocument();

  const [currentDoc, setCurrentDoc] = useState<DocgenDocument | null>(null);
  const [approvedDocId, setApprovedDocId] = useState<string | null>(null);

  const canVerdict = approvedDocId !== null;

  if (!id) {
    return <div className={styles.status}>Идентификатор дела не указан.</div>;
  }

  const handleGenerate = (template: 'indictment' | 'summons') => {
    generate.mutate(
      { case_id: id, template },
      {
        onSuccess: (d) => {
          setCurrentDoc(d);
          setApprovedDocId(null);
          void refetchDocs();
          toast.push({
            kind: 'success',
            title: 'Документ сгенерирован',
            message: 'Проверьте превью и утвердите.',
          });
        },
        onError: () => {
          toast.push({
            kind: 'error',
            title: 'Ошибка',
            message: 'Не удалось сгенерировать документ.',
          });
        },
      },
    );
  };

  const handleSelectDoc = (d: CaseDocument) => {
    setCurrentDoc(d as unknown as DocgenDocument);
    setApprovedDocId(null);
  };

  const handleApprove = (d: DocgenDocument) => {
    setApprovedDocId(d.id);
    toast.push({
      kind: 'success',
      title: 'Документ утверждён',
      message: 'Теперь можно вынести приговор.',
    });
  };

  return (
    <div className={styles.wrap}>
      <div className={styles.intro}>
        <h1 className={styles.title}>Вынесение приговора</h1>
        <span className={styles.subtitle}>Решение трибунала · Окончательно</span>
      </div>

      {(caseLoading || opinionLoading) && <div className={styles.status}>Загрузка…</div>}

      {caseData && (
        <div className={styles.block}>
          <span className={styles.blockTitle}>Материалы дела</span>
          <span className={styles.meta}>
            Дело №{caseData.seq_num}/Δ · {caseData.defendant}
          </span>
          <span className={styles.reasoning}>{caseData.crime}</span>
        </div>
      )}

      {opinion && (
        <div className={styles.block}>
          <span className={styles.blockTitle}>
            Заключение прокурора: {VERDICT_RU[opinion.preliminary_verdict]}
          </span>
          <span className={styles.reasoning}>{opinion.reasoning}</span>
        </div>
      )}

      <Paper variant="paper">
        <div className={styles.docsSection}>
          <span className={styles.blockTitle}>Документы по делу</span>
          {docs && docs.length > 0 && (
            <ul className={styles.docList}>
              {docs.map((d: CaseDocument) => (
                <li
                  key={d.id}
                  className={`${styles.docItem} ${currentDoc?.id === d.id ? styles.docItemActive : ''}`}
                  onClick={() => handleSelectDoc(d)}
                >
                  №{d.id.slice(0, 8)} · {d.template}
                </li>
              ))}
            </ul>
          )}
          <div className={styles.docActions}>
            <button
              className={styles.genBtn}
              onClick={() => handleGenerate('indictment')}
              disabled={generate.isPending}
            >
              {generate.isPending ? 'Генерация…' : 'Обвинительное заключение'}
            </button>
            <button
              className={styles.genBtn}
              onClick={() => handleGenerate('summons')}
              disabled={generate.isPending}
            >
              {generate.isPending ? 'Генерация…' : 'Повестка в трибунал'}
            </button>
          </div>
        </div>
      </Paper>

      <Paper variant="paper">
        <Preview doc={currentDoc} onApprove={handleApprove} />
      </Paper>

      <Paper variant="paper">
        <VerdictForm
          caseId={id}
          disabled={!canVerdict}
          gateHint={
            canVerdict
              ? undefined
              : 'Сначала сгенерируйте документ по делу и утвердите его в превью.'
          }
          onSuccess={(res) => {
            toast.push({
              kind: 'success',
              title: 'Приговор вынесен',
              message: 'Дело закрыто и помещено в архив.',
            });
            navigate(ROUTES.archive + '#' + res.archive_entry.id, { replace: true });
          }}
        />
      </Paper>
    </div>
  );
};
